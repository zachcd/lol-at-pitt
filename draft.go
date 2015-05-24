package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lab-d8/lol-at-pitt/draft"
)

type DraftHandler func(msg Message, room *DraftRoom)

/////////////////////////
const (
	startingCountdownTime = 20
	countdownEventTime    = 5
)

var (
	currentCountdown                         = startingCountdownTime
	allowTicks       bool                    = true // If this is false, dont continue to count down
	mainHandler      map[string]DraftHandler = map[string]DraftHandler{}
)

/////////////////////////

func RegisterDraftHandler(msg_type string, handle DraftHandler) {
	mainHandler[msg_type] = handle
}

func Handle(msg Message) {
	if mainHandler[msg.Type] != nil {
		mainHandler[msg.Type](msg, room)
	}
}

func Init() {
	timer_handler()

	RegisterDraftHandler("login", handle_login)
	RegisterDraftHandler("bid", handle_bid)
	RegisterDraftHandler("event", handle_event)
	RegisterDraftHandler("captains", handle_captains)
	RegisterDraftHandler("timer_reset", handle_timer_reset)
	RegisterDraftHandler("captains", handle_captains)
	RegisterDraftHandler("upcoming", handle_upcoming)
	RegisterDraftHandler("current-player", handle_current_player)
	RegisterDraftHandler("current-header", handle_header)
}

func handle_bid(msg Message, room *DraftRoom) {
	// TODO: Update with maria code
	amt, err := strconv.Atoi(msg.Text)

	if err == nil {
		formattedStr := fmt.Sprintf("<h5>Amount: <span  class='text-success'>%d</span></h5>", amt)
		go Handle(Message{Type: "event", Text: formattedStr})
	}
}

func handle_event(msg Message, room *DraftRoom) {
	room.broadcast(&msg)
}

func handle_captains(msg Message, room *DraftRoom) {
	// TODO: do formatting of text here. Make it a json blob
	text := ""
	format := `<li class='list-group-item'>%s (%s)<span class='text-info'> %d </span></li>`
	captains := draft.GetCaptains()
	for _, captain := range captains {
		res := fmt.Sprintf(format, captain.TeamName, captain.Name, captain.Points)
		text += res
	}
	room.broadcast(&Message{Type: "captains", Text: text})
}

func handle_upcoming(msg Message, room *DraftRoom) {
	text := ""
	format := `<li class='list-group-item'> %s <span class='text-muted'> %s </span></li>`
	players := draft.GetPlayers()
	for _, player := range players[1:len(players)] {
		res := fmt.Sprintf(format, player.Ign, player.Tier)
		text += res
	}
	room.broadcast(&Message{Type: "upcoming", Text: text})
}

func handle_current_player(msg Message, room *DraftRoom) {
	var format string = `
		<div class="row">
			<div class="col-md-3">%s</div>
			<div class="col-md-8">%s</div>
	</div>
	<div class="row">
			<div id="current_tier" class="col-md-3 text-muted">%s</div>
	</div>
	</div>
	`
	player := draft.GetPlayers()[0]

	res := fmt.Sprintf(format, player.Ign, player.Roles, player.Tier)
	room.broadcast(&Message{Type: "current-player", Text: res})
}

func handle_header(msg Message, room *DraftRoom) {
	player := draft.GetPlayers()[0]

	room.broadcast(&Message{Type: "current-header", Text: player.Ign})
}

// handle_login will give the player their stats, captains, current player, and upcoming players.
func handle_login(msg Message, room *DraftRoom) {
	Handle(Message{Type: "captains"})
	Handle(Message{Type: "upcoming"})
	Handle(Message{Type: "current-player"})
	Handle(Message{Type: "current-header"})
}

func handle_timer_reset(msg Message, room *DraftRoom) {
	currentCountdown = startingCountdownTime
}

func timer_handler() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			_ = now

			if !allowTicks {
				continue
			}

			currentCountdown--

			if currentCountdown < countdownEventTime {
				res := fmt.Sprintf("%d seconds remaining...", currentCountdown)
				Handle(Message{Type: "event", Text: res})
			}

			if currentCountdown == 0 {
				allowTicks = false
				Handle(Message{Type: "timer-end"})
			}
		}
	}()
}

package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lab-d8/lol-at-pitt/draft"
)

type DraftHandler func(msg Message, room *DraftRoom)

/////////////////////////
const (
	startingCountdownTime = 10
	countUpEventTime      = 10
	countdownEventTime    = 5
)

var (
	currentCountdown                         = startingCountdownTime
	allowTicks       bool                    = false // If this is false, dont continue to count down
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
	draft.Init()
	timer_handler()

	RegisterDraftHandler("login", handle_update)
	RegisterDraftHandler("update", handle_update)
	RegisterDraftHandler("bid", handle_bid)
	RegisterDraftHandler("bid-more", handle_more_bid)
	RegisterDraftHandler("bidder", handle_bidder)
	RegisterDraftHandler("event", handle_event)
	RegisterDraftHandler("captains", handle_captains)
	RegisterDraftHandler("upcoming", handle_upcoming)
	RegisterDraftHandler("refresh", handle_refresh)
	RegisterDraftHandler("timer-end", handle_timer_end)
	RegisterDraftHandler("current-player", handle_current_player)
	// winner
	// final-ticks
}

func handle_refresh(msg Message, room *DraftRoom) {
	draft.Init()
	Handle(Message{Type: "update"})
}

func handle_more_bid(msg Message, room *DraftRoom) {
	amt, err := strconv.Atoi(msg.Text)
	log.Println(msg, err)
	if err == nil {
		amount := draft.GetCurrentPlayer().HighestBid + amt
		Handle(Message{Type: "bid", From: msg.From, Text: strconv.Itoa(amount)})
	}
}

func handle_bid(msg Message, room *DraftRoom) {
	amt, err := strconv.Atoi(msg.Text)
	log.Println(msg)
	if err == nil {
		bidSuccess := draft.Bid(msg.From, amt)
		captain := draft.GetAuctioner(msg.From)
		if bidSuccess {
			formattedStr := fmt.Sprintf("<h5>%s bid <span  class='text-success'>%d</span> on <span class='text-success'>%s</span></h5>",
				captain.TeamName, amt, draft.GetCurrentPlayer().Ign)
			go Handle(Message{Type: "event", Text: formattedStr})
			currentCountdown = startingCountdownTime
			allowTicks = true
		}
	}
}

func handle_event(msg Message, room *DraftRoom) {
	room.broadcast(&msg)
}

func handle_captains(msg Message, room *DraftRoom) {
	text := ""
	format := `<li class='list-group-item'>%s (%s)<span class='text-info'> %d </span></li>`
	captains := draft.GetSortedCaptains()
	for _, captain := range captains {
		res := fmt.Sprintf(format, captain.TeamName, captain.Name, captain.Points)
		text += res
	}
	room.broadcast(&Message{Type: "captains", Text: text})
}

func handle_upcoming(msg Message, room *DraftRoom) {
	text := ""
	format := `<li class='list-group-item'> %s <span class='text-muted'> %d </span></li>`
	players := draft.GetPlayers()
	for _, player := range players {
		res := fmt.Sprintf(format, player.Ign, player.Score)
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
	player := draft.GetCurrentPlayer()
	res := fmt.Sprintf(format, player.Ign, player.Roles, player.Tier)
	room.broadcast(&Message{Type: "current-header", Text: player.Ign})
	room.broadcast(&Message{Type: "current-player", Text: res})
}

func handle_bidder(msg Message, room *DraftRoom) {
	captain := draft.GetAuctioner(msg.From)
	if captain != nil {
		str := fmt.Sprintf("%d", captain.Points)
		room.messageWithID(msg.From, &Message{Type: "points", Text: str})
		room.messageWithID(msg.From, &Message{Type: "team", Text: captain.TeamName})
	}

}

// handle_login will give the player their stats, captains, current player, and upcoming players.
func handle_update(msg Message, room *DraftRoom) {
	Handle(Message{Type: "captains"})
	Handle(Message{Type: "upcoming"})
	Handle(Message{Type: "current-player"})
	Handle(Message{Type: "current-header"})
	Handle(Message{Type: "event", Text: "Currently waiting to bid on.." + draft.GetCurrentPlayer().Ign})
	for _, client := range room.clients {
		Handle(Message{Type: "bidder", From: client.ID})
	}


}

func handle_winner(msg Message, room *DraftRoom) {
	Handle(Message{Type: "event", Text: draft.GetCurrentPlayer().Team + " bought" + draft.GetCurrentPlayer().Ign + " for " + strconv.Itoa(draft.GetCurrentPlayer().HighestBid)})
	draft.Win()
	Handle(Message{Type: "update"})
	draft.Paused = true
}

func handle_timer_reset(msg Message, room *DraftRoom) {
	currentCountdown = startingCountdownTime
}

func handle_timer_end(msg Message, room *DraftRoom) {
	current := draft.GetCurrentPlayer()
	if current.HighestBid > 0 {
		draft.Paused = true
		handle_winner(msg, room)
	}
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

package draft

import (
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"log"
	"sort"
	"strconv"
	"time"
)

type DraftPlayer struct {
	Id     int64
	Ign    string
	Bid    int
	Team   string
	Player ols.Player
}

type Bid struct {
	Team   string
	Amount int
}

type Auctioner struct {
	Id     int64
	Points int
	Team   string
}

type Draft struct {
	Current DraftPlayer
	History
	Unassigned []DraftPlayer
	Auctioners map[string]*Auctioner
	paused     bool
	started    bool
	dao        DraftDAO
	queuedBids chan Bid
}

type History struct {
	max    int
	Values []string
}

func InitHistory(size int) History {
	return History{max: size, Values: []string{}}
}

func (h *History) Add(val string) {
	h.Values = append([]string{val}, h.Values...)
	if len(h.Values) > h.max {
		h.Values = h.Values[:h.max]
	}

}

func InitNewDraft(db *mgo.Database) *Draft {
	draftees := []DraftPlayer{}
	auctioners := map[string]*Auctioner{}
	allPlayers := dao.GetPlayersDAO().All()

	captains := allPlayers.Filter(func(player ols.Player) bool {
		return player.Captain
	})
	players := allPlayers.Filter(func(player ols.Player) bool {
		return !player.Captain && player.Team == ""
	})

	sort.Sort(players)

	for _, captain := range captains {
		auctioners[captain.Team] = &Auctioner{Id: captain.Id, Team: captain.Team, Points: captain.Score}
	}

	for _, player := range players {
		draftPlay := DraftPlayer{Id: player.Id, Ign: player.Ign, Player: *player, Bid: -1}
		draftees = append(draftees, draftPlay)
	}

	var current DraftPlayer
	current, draftees = draftees[0], draftees[1:]
	draft := Draft{
		Current:    current,
		Unassigned: draftees,
		Auctioners: auctioners,
		History:    InitHistory(20),
		paused:     true,
		started:    false,
		queuedBids: make(chan Bid, 30),
	}
	go DraftRunner(&draft)
	draft.History.Add("Starting Draft..")
	return &draft
}

func Load(db *mgo.Database) *Draft {
	dao := InitDraftDAO(db)
	draft := dao.Load()
	draft.dao = dao
	draft.History = InitHistory(20)
	return draft
}

func (d *Draft) Pause() {
	d.History.Add("PAUSE: Draft is paused by the admins")
	d.paused = true
}

func (d *Draft) Resume() {
	d.History.Add("RESUME: Draft is now resuming")
	d.paused = false
}

func (d *Draft) Bid(amount int, team string) {
	log.Println("Bid...", amount)
	d.queuedBids <- Bid{team, amount}
}

// Returns: true if the bid went through, false otherwise
func (d *Draft) bid(amount int, team string) bool {
	auctioner, ok := d.Auctioners[team]

	if d.Current.Team == team || d.Current.Bid >= amount || !ok || auctioner.Points < amount {
		return false
	}

	d.Current.Bid = amount
	d.Current.Team = team
	d.History.Add(team + " bid " + strconv.Itoa(amount) + " points for " + d.Current.Ign)
	return true

}

func (d *Draft) ArePlayersLeft() bool {
	return len(d.Unassigned) > 0
}

func (d *Draft) Finalize() {
	d.Pause()
	d.started = false
	auctioner, _ := d.Auctioners[d.Current.Team]
	auctioner.Points -= d.Current.Bid
	log.Println(auctioner.Points)
	// Save player team
	player := dao.GetPlayersDAO().Load(d.Current.Id)
	player.Team = d.Current.Team
	player.Score = d.Current.Bid
	dao.GetPlayersDAO().Save(player)

	// Save captain point value
	captain := dao.GetPlayersDAO().Load(auctioner.Id)
	captain.Score = auctioner.Points
	dao.GetPlayersDAO().Save(captain)

	d.History.Add("WINNER: " + d.Current.Team + " won " + d.Current.Ign + " for " + strconv.Itoa(d.Current.Bid))
}

func (d *Draft) Start() {
	if d.started {
		return
	}
	d.started = true
	d.Resume()
	d.History.Add("STARTING: Now bidding on " + d.Current.Ign)
	go DraftTimer(d)
}

func (d *Draft) Next() {
	d.started = false
	d.Current, d.Unassigned = d.Unassigned[0], d.Unassigned[1:]
	d.History.Add("NEXT: Going to next player, " + d.Current.Ign)
}

func DraftRunner(draft *Draft) {
	log.Println("draft runner started..")
	for {
		log.Println("draft runnerrr", draft.paused)
		bid := <-draft.queuedBids
		if draft.paused {
			continue
		}
		log.Println("bidding", bid.Amount)
		draft.bid(bid.Amount, bid.Team)

	}
}

func DraftTimer(draft *Draft) {
	go func() {
		max_time := 12
		secondsExpired := 0
		lastBiddingTeam := ""
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			_ = now
			// Pause logic
			if draft.paused {
				continue
			}

			sameBidder := draft.Current.Team == lastBiddingTeam

			if draft.Current.Team == "" {
				secondsExpired = 0
			} else if sameBidder {
				secondsExpired += 1
			} else {
				secondsExpired = 0
			}

			lastBiddingTeam = draft.Current.Team

			if secondsExpired == max_time {
				draft.Pause()
				break
			} else if secondsExpired > 8 {
				draft.History.Add(strconv.Itoa(max_time-secondsExpired) + " left to outbid " + draft.Current.Team + " on " +
					draft.Current.Ign + " for " + strconv.Itoa(draft.Current.Bid))
			}
		}

		draft.Finalize()
	}()
}

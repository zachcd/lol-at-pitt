package draft

import (
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"time"
)

type DraftPlayer struct {
	Id   int64
	Ign  string
	Bid  int
	Team string
	Done bool // Whether the auction is done for that player
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
	Auctioners map[string]Auctioner
	paused     bool
	dao        DraftDAO
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

func InitNewDraft(db *mgo.Database) Draft {
	draftees := []DraftPlayer{}
	auctioners := map[string]Auctioner{}
	allPlayers := dao.GetPlayersDAO().All()

	captains := allPlayers.Filter(func(player ols.Player) bool {
		return player.Captain
	})
	players := allPlayers.Filter(func(player ols.Player) bool {
		return !player.Captain && player.Team == ""
	})

	for _, captain := range captains {
		auctioners[captain.Team] = Auctioner{Id: captain.Id, Team: captain.Team, Points: captain.Score}
	}

	for _, player := range players {
		draftPlay := DraftPlayer{Id: player.Id, Done: false, Ign: player.Ign}
		draftees = append(draftees, draftPlay)
	}

	var current DraftPlayer
	current, draftees = draftees[len(draftees)-1], draftees[:len(draftees)-1]
	draft := Draft{
		Current:    current,
		Unassigned: draftees,
		Auctioners: auctioners,
		History:    InitHistory(20),
		paused:     true,
	}

	draft.History.Add("Starting Draft..")
	return draft
}

func Load(db *mgo.Database) *Draft {
	dao := InitDraftDAO(db)
	draft := dao.Load()
	draft.dao = dao
	draft.History = InitHistory(20)
	return draft
}

func (d *Draft) Pause() {
	d.paused = true
}

func (d *Draft) Resume() {
	d.paused = false
}

// Returns: true if the bid went through, false otherwise
func (d *Draft) Bid(amount int, team string) bool {
	auctioner, ok := d.Auctioners[team]

	if d.Current.Team == team || d.Current.Bid >= amount || !ok || auctioner.Points < amount {
		return false
	}

	d.Current.Bid = amount
	d.Current.Team = team
	d.History.Add(team + " bid " + string(amount) + " points for " + d.Current.Ign)
	return true

}

func (d *Draft) ArePlayersLeft() bool {
	return len(d.Unassigned) > 0
}

func (d *Draft) Finalize() {
	d.Current.Done = true
	auctioner, _ := d.Auctioners[d.Current.Team]
	auctioner.Points -= d.Current.Bid
	d.dao.Save(d)

	player := dao.GetPlayersDAO().Load(d.Current.Id)
	player.Team = d.Current.Team
	dao.GetPlayersDAO().Save(player)

	d.History.Add(d.Current.Team + " won " + d.Current.Ign + " for " + string(d.Current.Bid))
}

func (d *Draft) NextPlayer() {
	d.Current, d.Unassigned = d.Unassigned[len(d.Unassigned)-1], d.Unassigned[:len(d.Unassigned)-1]
	d.History.Add("Now bidding on " + d.Current.Ign)
}

func DraftRunner(draft *Draft, bids chan Bid) {
	for {

		bid := <-bids
		draft.Bid(bid.Amount, bid.Team)
	}
}

// Returns a channel that will blah
func DraftTimer(draft *Draft) {
	go func() {
		secondsExpired := 0
		recentHistory := ""
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			_ = now
			if recentHistory != "" && recentHistory == draft.History.Values[0] {
				secondsExpired += 1
			} else {
				secondsExpired += 1
			}

			recentHistory = draft.History.Values[0]
		}
	}()
}

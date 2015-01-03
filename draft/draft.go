package draft

import (
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
)

type DraftPlayer struct {
	Id   int64
	Bid  int
	Team string
	Done bool // Whether the auction is done for that player
}

type Auctioner struct {
	Id     int64
	Points int
	Team   string
}

type Draft struct {
	Current    DraftPlayer
	Unassigned []DraftPlayer
	Auctioners map[string]Auctioner
	paused     bool
	dao        DraftDAO
}

func InitNewDraft(db *mgo.Database) Draft {
	herd := []DraftPlayer{}
	playerDAO := dao.NewPlayerContext(db)
	players := playerDAO.All()
	players = players.Filter(func(player ols.Player) bool {
		return !player.Captain
	})

	for _, player := range players {
		draftPlay := DraftPlayer{Id: player.Id, Done: false}
		herd = append(herd, draftPlay)
	}

	var current DraftPlayer
	current, herd = herd[len(herd)-1], herd[:len(herd)-1]
	draft := Draft{
		Current:    current,
		Unassigned: herd,
		paused:     true,
	}

	return draft
}

func Load(db *mgo.Database) *Draft {
	dao := InitDraftDAO(db)
	draft := dao.Load()
	draft.dao = dao
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
}

func (d *Draft) NextPlayer() {
	d.Current, d.Unassigned = d.Unassigned[len(d.Unassigned)-1], d.Unassigned[:len(d.Unassigned)-1]
}

package ols

import (
	"labix.org/v2/mgo"
)

type DraftedPlayer struct {
	Player    *Player
	Bid       int
	BuyerTeam string
}

type DraftedPlayerDAO struct {
	db *mgo.Database
}
type DraftedPlayers []*DraftedPlayers

func (d *DraftedPlayerDAO) QueryIgn(ign string) DraftedPlayer {
	var player Player
	d.db.C("players").Find(map[string]string{"Ign": ign}).One(&player)
	var draftedPlayer = DraftedPlayer{Player: &player}
	return draftedPlayer
}

func (d *DraftedPlayerDAO) QueryId(id string) DraftedPlayer {
	var player Player
	d.db.C("players").Find(map[string]string{"Id": id}).One(&player)
	return DraftedPlayer{Player: &player}
}

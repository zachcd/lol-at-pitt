package ols

import (
	"fmt"
	"labix.org/v2/mgo"
)

type Player struct {
	Ign           string
	Id            int64
	Name          string
	NormalizedIgn string
	Roles         []string
	Score         int
	Team          string
	Captain       bool
	Tier          string
}

type PlayersDAO struct {
	db         *mgo.Database
	collection *mgo.Collection
}
type Players []*Player

func (p *Players) Filter(filter func(Player) bool) Players {
	players := Players{}
	for _, player := range *p {
		if filter(*player) {
			players = append(players, player)
		}
	}

	return players
}

func (p *Players) Print() {
	for _, player := range *p {
		fmt.Println(player)
	}
}

func NewPlayerContext(db *mgo.Database) *PlayersDAO {
	dao := PlayersDAO{db, db.C("players")}
	return &dao
}

func (p *PlayersDAO) loadIGN(ign string) Player {
	var player Player
	p.collection.Find(map[string]string{"Ign": ign}).One(&player)
	return player
}

func (p *PlayersDAO) load(id string) Player {
	var player Player
	p.collection.Find(map[string]string{"Id": id}).One(&player)
	return player
}

func (p *PlayersDAO) All() Players {
	var players Players
	p.collection.Find(map[string]string{}).All(&players)
	return players
}

func (p *PlayersDAO) Save(player Player) {
	count, _ := p.collection.Find(player).Count()
	if count > 0 {
		p.collection.Update(map[string]int64{"Id": player.Id}, player)
	} else {
		p.collection.Insert(player)
	}
}

func (p *PlayersDAO) Delete(player Player) {
	p.collection.Remove(player)
}

/// DEPRECATED. Remove when you can.
func QueryIgn(db *mgo.Database, ign string) Player {
	var player Player
	db.C("players").Find(map[string]string{"Ign": ign}).One(&player)
	return player
}

func QueryId(db *mgo.Database, id string) Player {
	var player Player
	db.C("players").Find(map[string]string{"Id": id}).One(&player)
	return player
}

func QueryAllPlayers(db *mgo.Database) Players {
	var players Players
	db.C("players").Find(map[string]string{}).All(&players)
	return players
}

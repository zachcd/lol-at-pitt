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

func QueryAllPlayers(db *mgo.Database) Players {
	var players Players
	db.C("players").Find(map[string]string{}).All(&players)
	return players
}

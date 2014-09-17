package ols

import "fmt"

type Player struct {
	Ign           string
	Id            int64
	Name          string
	NormalizedIgn string
	Roles         []string
	Score         int
	Team          string
	Captain       bool
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

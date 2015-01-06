package ols

import (
	"fmt"
)

type Player struct {
	Ign             string
	Id              int64
	Name            string
	NormalizedIgn   string
	Roles           []string
	Score           int
	Team            string
	Captain         bool
	Tier            string
	Lolking         int
	RoleDescription string
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

func (p Players) Len() int {
	return len(p)
}

func (p Players) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Players) Less(i, j int) bool {
	return p[i].Lolking > p[j].Lolking
}

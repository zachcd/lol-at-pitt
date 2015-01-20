package main

import (
	"github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
)

func CreateTeamsFromPlayers() {
	var teams map[string]*ols.Team
	for _, player := range db.GetPlayersDAO().All() {
		_, ok := teams[player.Team]

		if !ok {
			team := ols.Team{}
			team.Players = []int64{}
			team.Name = player.Name
			teams[player.Team] = &team
		}

		team := teams[player.Team]
		team.Players = append(team.Players, player.Id)

		if player.Captain {
			team.Captain = player.Id
		}

	}

	for _, team := range teams {
		db.GetTeamsDAO().Save(*team)
	}

}

func UpdateTeamScore(name string, win bool) {
	team := db.GetTeamsDAO().Load(name)

	if win {
		team.Wins++
	} else {
		team.Losses++
	}

	db.GetTeamsDAO().Save(team)
}

func NewTeamScore(name string, wins int, losses int) {
	team := db.GetTeamsDAO().Load(name)
	team.Wins = wins
	team.Losses = losses
	db.GetTeamsDAO().Save(team)
}

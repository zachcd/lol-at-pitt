package main

import (
	"fmt"
	"github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
)

func ShowTeams() {
	teams := db.GetTeamsDAO().All()

	for _, team := range teams {
		fmt.Print("Team (", team.Name, "): ")
		for _, playerId := range team.Players {
			player := db.GetPlayersDAO().Load(playerId)
			fmt.Print(player.Ign, " ")
		}
		fmt.Println()
	}
}

func UpdateTeamName(teamName, updatedName string) {
	team := db.GetTeamsDAO().Load(teamName)
	updatedTeam := db.GetTeamsDAO().Load(teamName)
	updatedTeam.Name = updatedName
	db.GetTeamsDAO().Update(team, updatedTeam)
	matches := db.GetMatchesDAO().LoadTeamMatches(team.Name)

	for _, match := range matches {
		var newMatch ols.Match = *match

		if newMatch.BlueTeam == team.Name {
			newMatch.BlueTeam = updatedName

		}

		if newMatch.RedTeam == teamName {
			newMatch.RedTeam = updatedName
		}

		if newMatch.Winner == teamName {
			newMatch.Winner = updatedName
		}

		db.GetMatchesDAO().Update(*match, newMatch)
	}

	players := db.GetPlayersDAO().All()
	players = players.Filter(func(player ols.Player) bool {
		return player.Team == teamName
	})

	for _, player := range players {
		player.Team = updatedName
		db.GetPlayersDAO().Save(*player)
	}

}

func CreateTeamsFromPlayers() {
	teams := map[string]ols.Team{}
	for _, player := range db.GetPlayersDAO().All() {
		_, ok := teams[player.Team]

		if !ok {
			team := ols.Team{}
			team.Players = []int64{}
			team.Name = player.Team
			teams[player.Team] = team
		}

		team := teams[player.Team]
		team.Players = append(team.Players, player.Id)
		teams[player.Team] = team
		if player.Captain {
			team.Captain = player.Id
		}

	}

	for _, team := range teams {
		db.GetTeamsDAO().Save(team)
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

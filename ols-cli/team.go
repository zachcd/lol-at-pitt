package main

import (
	"fmt"

	"github.com/lab-d8/lol-at-pitt/ols"
)

func ShowTeams() {
	teams := ols.GetTeamsDAO().All()

	for _, team := range teams {
		fmt.Print("Team (", team.Name, "): ")
		for _, playerId := range team.Players {
			player := ols.GetPlayersDAO().Load(playerId)
			fmt.Print(player.Ign, " ")
		}
		fmt.Println()
	}
}

func UpdateTeamName(teamName, updatedName string) {
	team := ols.GetTeamsDAO().Load(teamName)
	updatedTeam := ols.GetTeamsDAO().Load(teamName)
	updatedTeam.Name = updatedName
	ols.GetTeamsDAO().Update(team, updatedTeam)
	matches := ols.GetMatchesDAO().LoadTeamMatches(team.Name)

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

		ols.GetMatchesDAO().Update(*match, newMatch)
	}
}

func AddTeamPlayer(teamName string, playerId int64) {
	team := ols.GetTeamsDAO().Load(teamName)
	team.Players = append(team.Players, playerId)
	ols.GetTeamsDAO().Save(team)
}

func UpdateTeamCaptain(teamName string, playerId int64) {
	team := ols.GetTeamsDAO().Load(teamName)
	team.Captain = playerId
	ols.GetTeamsDAO().Save(team)
}

func RemoveTeamPlayer(teamName string, removePlayerId int64) {
	team := ols.GetTeamsDAO().Load(teamName)
	newPlayers := []int64{}

	for _, playerId := range team.Players {
		if playerId != removePlayerId {
			newPlayers = append(newPlayers, playerId)
		}
	}
	team.Players = newPlayers
	ols.GetTeamsDAO().Save(team)

}

func ReplaceTeamPlayer(teamName string, newPlayer, oldPlayer int64) {
	RemoveTeamPlayer(teamName, oldPlayer)
	AddTeamPlayer(teamName, newPlayer)
}

func UpdateTeamScore(name string, win bool) {
	team := ols.GetTeamsDAO().Load(name)

	if win {
		team.Wins++
	} else {
		team.Losses++
	}

	ols.GetTeamsDAO().Save(team)
}

func NewTeamScore(name string, wins int, losses int) {
	team := ols.GetTeamsDAO().Load(name)
	team.Wins = wins
	team.Losses = losses
	ols.GetTeamsDAO().Save(team)
}

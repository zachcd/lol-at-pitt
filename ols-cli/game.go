package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"log"
)

func CheckGames() {
	players := db.GetPlayersDAO().All()
	var usedSet map[int64]bool
	var potentialGames []int64
	for _, player := range players {
		summonerId := player.Id
		games, err := goriot.RecentGameBySummoner("na", summonerId)
		if err != nil {
			log.Printf("Error: ", err.Error())
		}

		for _, game := range games {
				potentialGames = append(potentialGames, game.GameID)
				usedSet[game.GameID] = true
			}
		}
	}

	for _, gameId := range potentialGames {
		game, err := goriot.MatchByMatchID("na", false, gameId)
		if err != nil {
			log.Printf("Error: ", err.Error())
			continue
		}

	}

}

func CorrectTeams(match goriot.MatchDetail) bool {
	teamMatch := 3 // At least 3 members of the original team need to be together.

	match.ParticipantIdentities[0].ParticipantId
}

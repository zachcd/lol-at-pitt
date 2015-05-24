package draft

import (
	"encoding/csv"
	"os"

	"github.com/lab-d8/lol-at-pitt/ols"
)

type DraftPlayer struct {
	ols.Player
	HighestBidder string
	Team          string
}

type DraftCaptain struct {
	ols.Player
	FacebookID string
	TeamName   string
	Points     int
}

type DraftPlayers []*DraftPlayer
type DraftCaptains []*DraftCaptain

func GetPlayers() DraftPlayers {
	players := ols.GetPlayersDAO().All()
	draftPlayers := []*DraftPlayer{}
	for _, player := range players {
		team := ols.GetTeamsDAO().LoadPlayer(player.Id)
		if player.Score != 0 && team.Captain != player.Id {
			draftPlayers = append(draftPlayers, &DraftPlayer{Player: *player})
		}
	}

	return draftPlayers
}

func GetCaptains() DraftCaptains {
	captains := ols.GetPlayersDAO().All()

	draftCaptains := []*DraftCaptain{}
	for _, player := range captains {
		team := ols.GetTeamsDAO().LoadPlayer(player.Id)
		if team.Captain == player.Id {
			user := ols.GetUserDAO().GetUserLeague(player.Id)
			draftCaptains = append(draftCaptains, &DraftCaptain{Player: *player, FacebookID: user.FacebookId, Points: player.Score, TeamName: team.Name})
		}
	}

	return draftCaptains
}

func getCSVData(filename string) [][]string {
	r, _ := os.Open(filename)
	csvReader := csv.NewReader(r)
	allData, _ := csvReader.ReadAll()
	return allData
}

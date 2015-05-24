package draft

import (
	"sync"

	"github.com/lab-d8/lol-at-pitt/ols"
)

type DraftPlayer struct {
	ols.Player
	HighestBid int
	Team       string
}

type DraftCaptain struct {
	ols.Player
	FacebookID string
	TeamName   string
	Points     int
}

type DraftPlayers []*DraftPlayer
type DraftCaptains []*DraftCaptain

var (
	Paused          bool         = false
	lock            sync.Mutex   = sync.Mutex{}
	previous        DraftPlayers = DraftPlayers{}
	current         *DraftPlayer
	upcomingPlayers DraftPlayers
	captains        map[string]*DraftCaptain = map[string]*DraftCaptain{}
)

func Init() {
	upcomingPlayers = getPlayers()
	if len(upcomingPlayers) > 0 {
		current, upcomingPlayers = upcomingPlayers[0], upcomingPlayers[1:]
	} else {
		current = &DraftPlayer{}
	}

	allCaptains := getCaptains()

	for _, captain := range allCaptains {
		captains[captain.FacebookID] = captain
	}
}

func GetCurrentPlayer() *DraftPlayer {
	return current
}

func GetPlayers() DraftPlayers {
	return upcomingPlayers
}

func GetCaptains() map[string]*DraftCaptain {
	return captains
}

func GetAuctionerByTeam(team string) *DraftCaptain {
	for _, captain := range captains {
		if captain.TeamName == team {
			return captain
		}
	}

	return nil
}

func GetAuctioner(id string) *DraftCaptain {
	captain, ok := captains[id]

	if ok {
		return captain
	} else {
		return nil
	}
}

func Bid(id string, amount int) bool {
	captain := GetAuctioner(id)
	bidSuccessful := false
	if captain != nil {
		lock.Lock()
		if captain.TeamName != current.Team && amount > current.HighestBid && amount <= captain.Points && !Paused {
			current.Team = captain.TeamName
			current.HighestBid = amount
			bidSuccessful = true
		}
		lock.Unlock()
	}

	return bidSuccessful
}

func Win() {
	lock.Lock()
	captain := GetAuctionerByTeam(current.Team)
	if captain != nil {
		captain.Points -= current.HighestBid
	}
	Paused = true

	previous = append(previous, current)
	if len(upcomingPlayers) != 0 {
		current = upcomingPlayers[0]
		upcomingPlayers = upcomingPlayers[1:]
	}
	lock.Unlock()

}

func TogglePause() {
	Paused = !Paused
}

func Next() {
	lock.Lock()
	previous = append(previous, current)
	if len(upcomingPlayers) != 0 {

		current = upcomingPlayers[0]
		upcomingPlayers = upcomingPlayers[1:]
	}
	lock.Unlock()
}

func Previous() {
	lock.Lock()
	currentArr := DraftPlayers{}
	currentArr = append(currentArr, current)
	upcomingPlayers = append(currentArr, upcomingPlayers...)
	current = previous[len(previous)-1]
	previous = previous[:len(previous)-1]

	// Refund logic.
	captain := GetAuctionerByTeam(current.Team)

	if captain != nil {
		captain.Points += current.HighestBid
	}
	current.HighestBid = 0
	current.Team = ""
	lock.Unlock()

}

// Setup stuff
func getPlayers() DraftPlayers {
	players := ols.GetPlayersDAO().All()
	draftPlayers := []*DraftPlayer{}
	for _, player := range players {
		team := ols.GetTeamsDAO().LoadPlayer(player.Id)
		if team.Captain != player.Id {
			draftPlayers = append(draftPlayers, &DraftPlayer{Player: *player})
		}
	}

	return draftPlayers
}

func getCaptains() DraftCaptains {
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

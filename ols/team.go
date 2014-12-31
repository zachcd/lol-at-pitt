package ols

import (
	"labix.org/v2/mgo"
	"sort"
	"strings"
)

type Team struct {
	Name           string
	NormalizedName string
	Players        Players
	Captain        *Player
	Wins           int
	Losses         int
	//	Games   Games
}

type Teams []*Team

// Sorting functions
func (p Teams) Len() int {
	return len(p)
}

func (p Teams) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Teams) Less(i, j int) bool {
	if p[i].Wins == p[j].Wins {
		return p[i].Losses < p[j].Losses
	} else {
		return p[i].Wins > p[j].Wins
	}
}

////////////////////// DAOS

func QueryAllTeams(db *mgo.Database) Teams {
	var teams Teams
	db.C("teams").Find(map[string]string{}).All(&teams)
	realTeams := Teams{}
	for _, team := range teams {
		realTeams = append(realTeams, QueryTeam(db, team.Name))
	}
	sort.Sort(realTeams)

	return realTeams
}

func QueryTeam(db *mgo.Database, teamName string) *Team {
	normalizedTeamName := NormalizedName(teamName)
	var team Team
	var players Players
	db.C("teams").Find(map[string]string{"normalizedname": normalizedTeamName}).One(&team)
	db.C("players").Find(map[string]string{"team": team.Name}).All(&players)
	team.Players = players
	team.Captain = GetCaptain(players)
	return &team
}

func GetCaptain(players Players) *Player {
	for _, player := range players {
		if player.Captain {
			return player
		}
	}
	return nil
}

func NormalizedName(name string) string {
	lowercase_name := strings.ToLower(name)
	normalized_name := strings.Replace(lowercase_name, " ", "", -1)
	return normalized_name
}

package db

import (
	"github.com/lab-d8/lol-at-pitt/ols"
	"testing"
)

func TestTeamsSave(t *testing.T) {
	TeamCollectionName = "testteams"
	team := ols.Team{Name: "pew",
		Players: []int64{1, 2, 3, 4, 5},
		Wins:    1}

	GetTeamsDAO().Save(team)

	count, _ := GetTeamsDAO().Collection.Count()
	var savedTeam ols.Team
	GetTeamsDAO().Collection.Find(map[string]string{"name": team.Name}).One(&savedTeam)

	if count != 1 || savedTeam.Wins != 1 {
		t.Error("Failed saving the team")
	}

	GetTeamsDAO().Delete(team)
}

func TestTeamsFindPlayer(t *testing.T) {
	TeamCollectionName = "testteams"
	team := ols.Team{Name: "pew",
		Players: []int64{1, 2, 3},
		Wins:    1}

	GetTeamsDAO().Save(team)

	team = ols.Team{Name: "rawr",
		Players: []int64{4},
		Wins:    1}

	GetTeamsDAO().Save(team)
	savedTeam := GetTeamsDAO().LoadPlayer(1)

	if savedTeam.Name != "pew" {
		t.Error("Failed saving the team", savedTeam)
	}

	GetTeamsDAO().DeleteAll()
}

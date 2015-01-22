package ols

import (
	"testing"
)

func TestMatchesDAOSave(t *testing.T) {
	MatchesCollectionName = "testMatches"
	match := Match{Id: 2, BlueTeam: "blue"}
	GetMatchesDAO().Save(match)
	c := GetMatchesDAO().Collection

	count, err := c.Find(map[string]int64{"id": match.Id}).Count()
	var savedMatch Match
	c.Find(map[string]int64{"id": match.Id}).One(&savedMatch)
	if count != 1 || err != nil || savedMatch.BlueTeam != "blue" {
		t.Error("Didn't save properly", savedMatch, err)
	}

	GetMatchesDAO().Delete(match)
}

func TestMatchesDAOLoad(t *testing.T) {
	MatchesCollectionName = "testMatches"
	match := Match{Id: 2, BlueTeam: "blue"}
	GetMatchesDAO().Save(match)

	loadedMatch := GetMatchesDAO().Load(2)

	if loadedMatch.BlueTeam != "blue" {
		t.Error("Load not working properly")
	}

	GetMatchesDAO().Collection.DropCollection()

}

func TestMatchesDAOLoadTeam(t *testing.T) {
	MatchesCollectionName = "testMatches"
	match := Match{Id: 0, BlueTeam: "blue", RedTeam: "redteam1", Week: 1}
	GetMatchesDAO().Save(match)
	match = Match{Id: 0, BlueTeam: "blue", RedTeam: "redteam2", Week: 2}
	GetMatchesDAO().Save(match)
	loadedMatches := GetMatchesDAO().LoadTeamMatches("blue")

	if len(loadedMatches) != 2 {
		t.Error("Match length incorrect", loadedMatches)
	}

	if loadedMatches[0].RedTeam != "redteam1" {
		t.Error("Match sorting is incorrect", loadedMatches[0])
	}

	if loadedMatches[1].RedTeam != "redteam2" {
		t.Error("Match sorting is incorrect", loadedMatches[1])
	}

	GetMatchesDAO().Collection.DropCollection()
}

func TestMatchesDAOWinningLoadTeam(t *testing.T) {
	MatchesCollectionName = "testMatches"
	match := Match{Id: 0, BlueTeam: "blue", RedTeam: "redteam1", Week: 1}
	GetMatchesDAO().Save(match)
	match = Match{Id: 0, BlueTeam: "blue", RedTeam: "redteam2", Week: 2, Played: true, Winner: "blue"}
	GetMatchesDAO().Save(match)
	loadedMatches := GetMatchesDAO().LoadWinningMatches("blue")

	if len(loadedMatches) != 1 {
		t.Error("Match length incorrect", loadedMatches)
	}

	if loadedMatches[0].RedTeam != "redteam2" {
		t.Error("Match sorting is incorrect", loadedMatches[0])
	}

	GetMatchesDAO().Collection.DropCollection()
}

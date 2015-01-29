package ols

import (
	"github.com/TrevorSStone/goriot"
	"labix.org/v2/mgo"
	"sort"
)

var MatchesCollectionName string = "olsmatches"
var LeagueMatchesCollectionName string = "leaguematches"

type MatchesDAO struct {
	DAO
}

func NewMatchesContext(db *mgo.Database) *MatchesDAO {
	dao := MatchesDAO{DAO{db, db.C(MatchesCollectionName)}}
	return &dao
}

func (m *MatchesDAO) IsSaved(id int64) bool {
	val, _ := m.Collection.Find(map[string]int64{"id": id}).Count()
	return val > 0
}

func (m *MatchesDAO) Load(id int64) Match {
	var match Match
	m.Collection.Find(map[string]int64{"id": id}).One(&match)
	return match
}

func (m *MatchesDAO) Save(match Match) {
	m.DAO.Save(map[string]interface{}{"week": match.Week, "blueteam": match.BlueTeam, "redteam": match.RedTeam}, match)
}

func (m *MatchesDAO) Update(oldMatch, match Match) {
	m.Collection.Update(oldMatch, match)
}

func (m *MatchesDAO) LoadWeekForMatch(blueTeam string, redTeam string) int {
	var matches []Match
	m.Collection.Find(map[string]interface{}{"blueteam": blueTeam, "redteam": redTeam, "played": false}).All(&matches)

	closestWeek := 100 // Only goes to 8, lol
	for _, match := range matches {
		if closestWeek > match.Week {
			closestWeek = match.Week
		}
	}

	return closestWeek

}

func (m *MatchesDAO) LoadTeamMatches(team string) []*Match {
	var matches []*Match
	var matchesRed []*Match
	m.Collection.Find(map[string]string{"blueteam": team}).All(&matches)
	m.Collection.Find(map[string]string{"redteam": team}).All(&matchesRed)

	allMatches := append(matches, matchesRed...)
	sort.Sort(Matches(allMatches))
	return allMatches
}

func (m *MatchesDAO) SaveLeagueGame(match goriot.MatchDetail) {
	m.db.C(LeagueMatchesCollectionName).Insert(match)
}

func (m *MatchesDAO) IsLeagueGameSaved(match goriot.MatchDetail) bool {
	count, _ := m.db.C(LeagueMatchesCollectionName).Find(map[string]int64{"matchid": match.MatchID}).Count()
	return count > 0
}

func (m *MatchesDAO) LoadLeagueGame(matchId int64) goriot.MatchDetail {
	var match goriot.MatchDetail
	m.db.C(LeagueMatchesCollectionName).Find(map[string]int64{"matchid": matchId}).One(&match)
	return match
}

func (m *MatchesDAO) LoadWinningMatches(team string) []*Match {
	var matches []*Match
	m.Collection.Find(map[string]string{"winner": team}).All(&matches)

	sort.Sort(Matches(matches))
	return matches
}

func (m *MatchesDAO) All() []*Match {
	var matches []*Match
	m.Collection.Find(map[string]string{}).All(&matches)

	sort.Sort(Matches(matches))
	return matches
}

func (m *MatchesDAO) Delete(match Match) {
	m.Collection.Remove(match)
}

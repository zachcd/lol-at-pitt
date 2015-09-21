package ols

import (
	"fmt"

	"labix.org/v2/mgo"
)

type TeamsDAO struct {
	DAO
}

var TeamCollectionName string = "teams"

func NewTeamsContext(db *mgo.Database) *TeamsDAO {
	dao := TeamsDAO{DAO{db, db.C(TeamCollectionName)}}
	return &dao
}

func (t *TeamsDAO) Load(name string) Team {
	var team Team
	t.Collection.Find(map[string]string{"name": name}).One(&team)
	return team
}
func (t *TeamsDAO) LoadPlayerDisplay(summonerId int64) TeamDisplay {
	team := t.LoadPlayer(summonerId)
	if team.Name == "" {
		team = t.LoadPlayerByCaptain(summonerId)
	}
	fmt.Println(summonerId)
	fmt.Println(team)
	teamDisplay := TeamDisplay{Name: team.Name, Wins: team.Wins, Losses: team.Losses}
	players := []Player{}
	for _, playerId := range team.Players {
		player := GetPlayersDAO().Load(playerId)
		players = append(players, player)
	}

	captainPlayer := GetPlayersDAO().Load(team.Captain)
	teamDisplay.Captain = captainPlayer
	teamDisplay.Players = players
	return teamDisplay
}
func (t *TeamsDAO) LoadPlayer(summonerId int64) Team {
	var team Team
	t.Collection.Find(map[string]int64{"players": summonerId}).One(&team)
	return team
}

func (t *TeamsDAO) LoadPlayerByCaptain(summonerId int64) Team {
	var team Team
	t.Collection.Find(map[string]int64{"captain": summonerId}).One(&team)
	return team

}

func (t *TeamsDAO) All() Teams {
	var teams Teams
	t.Collection.Find(map[string]string{}).All(&teams)
	return teams
}

func (t *TeamsDAO) Update(team, updatedTeam Team) {
	t.Collection.Update(team, updatedTeam)
}

func (t *TeamsDAO) Save(team Team) {
	t.DAO.Save(map[string]string{"name": team.Name}, team)
}

func (t *TeamsDAO) DeleteAll() {
	t.DAO.Collection.DropCollection()
}

func (t *TeamsDAO) Delete(team Team) {
	t.DAO.Collection.Remove(team)
}

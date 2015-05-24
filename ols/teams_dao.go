package ols

import (
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

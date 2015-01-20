package db

import (
	"github.com/lab-d8/lol-at-pitt/ols"
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

func (t *TeamsDAO) Load(name string) ols.Team {
	var team ols.Team
	t.Collection.Find(map[string]string{"name": name}).One(&team)
	return team
}

func (t *TeamsDAO) LoadPlayer(summonerId int64) ols.Team {
	var team ols.Team
	t.Collection.Find(map[string]int64{"players": summonerId}).One(&team)
	return team
}

func (t *TeamsDAO) All() ols.Teams {
	var teams ols.Teams
	t.Collection.Find(map[string]string{}).All(&teams)
	return teams
}

func (t *TeamsDAO) Save(team ols.Team) {
	t.DAO.Save(map[string]string{"name": team.Name}, team)
}

func (t *TeamsDAO) DeleteAll() {
	t.DAO.Collection.DropCollection()
}

func (t *TeamsDAO) Delete(team ols.Team) {
	t.DAO.Collection.Remove(team)
}

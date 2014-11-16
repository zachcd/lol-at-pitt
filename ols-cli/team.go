package main

import (
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
)

func UpdateTeamScore(name string, win bool) {
	session, _ := mgo.Dial(MongoLocation)
	db := session.DB(DatabaseName)
	team := ols.QueryTeam(db, name)
	session.Close()
	if win {
		NewTeamScore(name, team.Wins+1, team.Losses)
	} else {
		NewTeamScore(name, team.Wins, team.Losses+1)
	}
}

func NewTeamScore(name string, wins int, losses int) {
	session, _ := mgo.Dial(MongoLocation)
	db := session.DB(DatabaseName)
	team := ols.Team{Name: name, Wins: wins, Losses: losses}
	selector := ols.Team{Name: name}
	db.C("teams").Update(selector, team)
	session.Close()
}

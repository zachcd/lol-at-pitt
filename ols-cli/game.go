package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
)

func CheckGames() {
	session, _ := mgo.Dial(MongoLocation)
	db := session.DB(DatabaseName)

}

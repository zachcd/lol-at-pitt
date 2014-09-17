package main

import (
	"encoding/json"
	"fmt"
	"github.com/lab-d8/lol-at-pitt/ols"
	"io/ioutil"
	"labix.org/v2/mgo"
	"os"
)

const DatabaseName string = "lolpitt"
const MongoLocation = "mongodb://localhost"
const InputJson string = "resources/ols_players.json"
const InputTeamJson string = "resources/teams.json"

func main() {
	initDbPlayers()
	initDbTeams()
}

func initDbPlayers() {
	file, _ := os.Open(InputJson)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	var players ols.Players

	json.Unmarshal(data, &players)
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("players").DropCollection()
	for _, player := range players {
		db.C("players").Insert(player)
		fmt.Println("here..")
	}
	session.Close()

}

func initDbTeams() {
	file, _ := os.Open(InputTeamJson)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	var teams ols.Teams

	json.Unmarshal(data, &teams)

	session, err := mgo.Dial(MongoLocation)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	db := session.DB(DatabaseName)
	db.C("teams").DropCollection()

	for _, team := range teams {
		db.C("teams").Insert(team)
	}
}

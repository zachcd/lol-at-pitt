package main

import (
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"time"
)

// This is used for chron jobs. Reads in the
const MongoLocation = "mongodb://localhost"
const DatabaseName string = "lolpitt"

func main() {
	goriot.SetAPIKey(ApiKey)
	goriot.SetLongRateLimit(500, 10*time.Minute)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	session, err := mgo.Dial(MongoLocation)
	if err != nil {
		panic(err)
	}

	db := session.DB(DatabaseName)
	var players ols.Players
	db.C("players").Find(map[string]string{}).All(&players)

	for _, player := range players {
		id := player.Id
		leagues_by_id, err := goriot.LeagueBySummoner("na", id)
		if err != nil {
			fmt.Println("wat: ", err.Error())
			player.Tier = "None"
		}
		league, ok := leagues_by_id[id]
		if ok {
			player.Tier = getBestLeague(league, *player)
		}
		fmt.Println("player: ", player)
		db.C("players").Update(map[string]int64{"id": player.Id}, player)
	}

}

func getBestLeague(leagues []goriot.League, player ols.Player) string {
	standings := map[string]int{
		"BRONZE":     0,
		"SILVER":     1,
		"GOLD":       2,
		"PLATINUM":   3,
		"DIAMOND":    4,
		"MASTER":     5,
		"CHALLENGER": 6,
	}

	division_standings := map[string]int{"V": 5, "IV": 4, "III": 3, "II": 2, "I": 1}

	currentTier := "BRONZE"
	currentDivision := "V" // Bronze 5 pleb. Get better
	for _, league := range leagues {
		if standings[currentTier] <= standings[league.Tier] {
			currentTier = league.Tier
			for _, entry := range league.Entries {
				if entry.PlayerOrTeamName == player.Ign && division_standings[currentDivision] > division_standings[entry.Division] {
					currentDivision = entry.Division
				}
			}
		}
	}

	return currentTier + " " + currentDivision
}

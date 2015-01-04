package main

// The idea of this package is to provide a CLI to edit the database for Mongodb.
import (
	"encoding/csv"
	"fmt"
	"github.com/TrevorSStone/goriot"
	"github.com/docopt/docopt-go"
	dao "github.com/lab-d8/lol-at-pitt/db"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
	"os"
	"strconv"
	"strings"
	"time"
)

const ApiKey string = "a3c96054-e21f-4238-a842-28caa10943a0"

type CmdArgs map[string]interface{}
type Runnable func(map[string]interface{}) bool
type Command struct {
	Runnable               // used for testing whether a command is to be run
	Cmd      func(CmdArgs) // The actual function to run
}

type DB struct {
	Players ols.Players
	Teams   ols.Teams
}

const DatabaseName string = "lolpitt"
const MongoLocation = "mongodb://localhost"

// All possible Command line commands.
var cmds []Command = []Command{
	Command{Runnable: runnableGenerator("db", "dump"), Cmd: func(m CmdArgs) {
		dumpDb(m["<olsfile>"].(string))
	}},
	Command{Runnable: runnableGenerator("db", "upload"), Cmd: func(m CmdArgs) {
		upload(m["<olsfile>"].(string))
	}},
	Command{Runnable: runnableGenerator("db", "atomic_delete"), Cmd: func(m CmdArgs) {
		deleteDb()
	}},
	Command{Runnable: runnableGenerator("player", "file"), Cmd: func(m CmdArgs) {
		filePlayer(m["<file_name>"].(string))
	}},
	Command{Runnable: runnableGenerator("captain", "file"), Cmd: func(m CmdArgs) {
		newCaptain(m["<file_name>"].(string))
	}},
	Command{Runnable: runnableGenerator("user", "new"), Cmd: func(m CmdArgs) {
	}},
	Command{Runnable: runnableGenerator("user", "update"), Cmd: func(m CmdArgs) {
		player := dao.GetPlayersDAO().LoadIGN(m["<ign>"].(string))
		fmt.Println(m)
		_, ok := m["--captain"]
		if ok && m["--captain"].(string) == "true" {
			player.Captain = true
		} else if ok && m["--captain"].(string) == "false" {
			player.Captain = false
		}

		dao.GetPlayersDAO().Save(player)
	}},
	Command{Runnable: runnableGenerator("team", "score", "--win"), Cmd: func(m CmdArgs) {
		UpdateTeamScore(m["<name>"].(string), true)
	}},
	Command{Runnable: runnableGenerator("team", "score", "--lose"), Cmd: func(m CmdArgs) {
		UpdateTeamScore(m["<name>"].(string), false)
	}},
	Command{Runnable: runnableGenerator("team", "new_score"), Cmd: func(m CmdArgs) {
		wins, _ := strconv.Atoi(m["<wins>"].(string))
		losses, _ := strconv.Atoi(m["<losses>"].(string))
		NewTeamScore(m["<name>"].(string), wins, losses)
	}},
	Command{Runnable: runnableGenerator("team", "update"), Cmd: func(m CmdArgs) {

	}},
	Command{Runnable: runnableGenerator("tiers"), Cmd: func(m CmdArgs) {
		tiers()
	}},
}

func main() {
	goriot.SetAPIKey(ApiKey)
	goriot.SetLongRateLimit(500, 10*time.Minute)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	usage := `OLS CLI

Usage:
   ols-cli captain file <file_name>
   ols-cli player file <file_name>
   ols-cli user new <name> <ign> <email>
   ols-cli user update <ign> [--team=<newteam>|--captain=<bool>|--email=<email>|--ign=<newign>]
   ols-cli team score <name> [--win|--lose]
   ols-cli team new_score <wins> <losses>
   ols-cli team update <name> [--name=<newname>]
   ols-cli db dump <olsfile>
   ols-cli db upload <olsfile>
   ols-cli db atomic_delete
   ols-cli tiers
`
	arguments, _ := docopt.Parse(usage, nil, true, "ols-cli 1.0", false)

	for _, cmd := range cmds {
		if cmd.Runnable(arguments) {
			cmd.Cmd(arguments)
		}
	}

}

// Makes an easy to use runnable function
func runnableGenerator(args ...string) Runnable {
	return func(sys_args map[string]interface{}) bool {
		for _, arg := range args {
			if !sys_args[arg].(bool) {
				return false
			}
		}

		return true
	}
}

func update_user_name(ign string, updated_ign string) {

}

func tiers() {

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
func filePlayer(fileName string) {
	csvfile, err := os.Open(fileName)
	defer csvfile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := csv.NewReader(csvfile)
	rawCSVData, err := reader.ReadAll()
	reader.FieldsPerRecord = -1
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, val := range rawCSVData {
		// Prints the headers...dumb thing.
		if i == 0 {
			continue
		}

		name := strings.TrimSpace(val[1])
		ign := strings.TrimSpace(val[2])
		lolking, _ := strconv.Atoi(strings.TrimSpace(val[3]))

		player := dao.GetPlayersDAO().LoadNormalizedIGN(ign)

		if player.Id != 0 {
			player.Name = name
			player.Lolking = lolking
			fmt.Println("Player saved:", player)
		} else {
			normalizedSummonerName := goriot.NormalizeSummonerName(ign)[0]
			result, err := goriot.SummonerByName("na", normalizedSummonerName)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			summonerProfile := result[normalizedSummonerName]
			player = ols.Player{Id: summonerProfile.ID, Ign: summonerProfile.Name, NormalizedIgn: normalizedSummonerName, Name: name, Lolking: lolking}
			fmt.Println("Player new:", player)
		}

		dao.GetPlayersDAO().Save(player)

	}
}
func newCaptain(fileName string) {
	csvfile, err := os.Open(fileName)
	defer csvfile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := csv.NewReader(csvfile)
	rawCSVData, err := reader.ReadAll()
	reader.FieldsPerRecord = -1
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(rawCSVData, fileName)
	for i, val := range rawCSVData {

		// Prints the headers...dumb thing.
		if i == 0 {
			continue
		}

		ign := strings.TrimSpace(val[0])
		team := strings.TrimSpace(val[1])
		points, _ := strconv.Atoi(strings.TrimSpace(val[2]))
		normalizedSummonerName := goriot.NormalizeSummonerName(ign)[0]
		result, err := goriot.SummonerByName("na", normalizedSummonerName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		summonerProfile := result[normalizedSummonerName]
		player := ols.Player{Id: summonerProfile.ID, Ign: summonerProfile.Name, Score: points, Captain: true, Team: team, NormalizedIgn: normalizedSummonerName}
		fmt.Println("Success")
		dao.GetPlayersDAO().Save(player)

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

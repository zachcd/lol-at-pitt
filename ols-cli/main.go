package main

// The idea of this package is to provide a CLI to edit the database for Mongodb.
import (
	"fmt"
	"strconv"
	"time"

	"github.com/TrevorSStone/goriot"
	"github.com/docopt/docopt-go"
	"github.com/lab-d8/lol-at-pitt/ols"
	"labix.org/v2/mgo"
)

const ApiKey string = "98b86371-9734-4313-b252-79bffc69e804"

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
	Command{Runnable: runnableGenerator("user", "new"), Cmd: func(m CmdArgs) {
		NewPlayer(m["<name>"].(string), m["<ign>"].(string))
	}},
	Command{Runnable: runnableGenerator("team", "score", "--win"), Cmd: func(m CmdArgs) {
		UpdateTeamScore(m["<name>"].(string), true)
	}},
	Command{Runnable: runnableGenerator("team", "score", "--lose"), Cmd: func(m CmdArgs) {
		UpdateTeamScore(m["<name>"].(string), false)
	}},

	Command{Runnable: runnableGenerator("team", "name"), Cmd: func(m CmdArgs) {
		UpdateTeamName(m["<name>"].(string), m["<newname>"].(string))
	}},
	Command{Runnable: runnableGenerator("team", "new_score"), Cmd: func(m CmdArgs) {
		wins, _ := strconv.Atoi(m["<wins>"].(string))
		losses, _ := strconv.Atoi(m["<losses>"].(string))
		NewTeamScore(m["<name>"].(string), wins, losses)
	}},
	Command{Runnable: runnableGenerator("error", "names"), Cmd: func(m CmdArgs) {
		nameErrors()
	}},
	Command{Runnable: runnableGenerator("team", "stats"), Cmd: func(m CmdArgs) {
		ShowTeams()
	}},
	Command{Runnable: runnableGenerator("team", "add"), Cmd: func(m CmdArgs) {
		id, _ := strconv.ParseInt(m["<id>"].(string), 10, 64)
		AddTeamPlayer(m["<team>"].(string), id)
	}},
	Command{Runnable: runnableGenerator("team", "remove"), Cmd: func(m CmdArgs) {
		id, _ := strconv.ParseInt(m["<id>"].(string), 10, 64)
		RemoveTeamPlayer(m["<team>"].(string), id)
	}},
	Command{Runnable: runnableGenerator("update", "tiers"), Cmd: func(m CmdArgs) {
		tiers()
	}},
	Command{Runnable: runnableGenerator("update", "names"), Cmd: func(m CmdArgs) {
		nameUpdates()
	}},
	Command{Runnable: runnableGenerator("matches"), Cmd: func(m CmdArgs) {
		CheckGames()
		UpdateMatches()
	}},
	Command{Runnable: runnableGenerator("player"), Cmd: func(m CmdArgs) {
		UploadPlayers(m["<upload>"].(string))
	}},
	Command{Runnable: runnableGenerator("captain"), Cmd: func(m CmdArgs) {
		UploadCaptains(m["<upload>"].(string))
	}},
}

func main() {
	goriot.SetAPIKey(ApiKey)
	goriot.SetLongRateLimit(500, 10*time.Minute)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	usage := `OLS CLI

Usage:
   ols-cli user new <name> <ign>
   ols-cli team score <name> [--win|--lose]
   ols-cli team new_score <wins> <losses>
   ols-cli team name <name> <newname>
   ols-cli team stats
   ols-cli team remove <team> <id>
   ols-cli team add <team> <id>
   ols-cli db dump <olsfile>
   ols-cli db upload <olsfile>
   ols-cli db atomic_delete
   ols-cli update names
   ols-cli update tiers
   ols-cli update teams
   ols-cli matches
   ols-cli error names
   ols-cli player <upload>
   ols-cli captain <upload>
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

func NewPlayer(name, ign string) *ols.Player {
	player := ols.Player{Name: name}

	leaguePlayerMap, err := goriot.SummonerByName("na", goriot.NormalizeSummonerName(ign)[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}
	leaguePlayer := leaguePlayerMap[goriot.NormalizeSummonerName(ign)[0]]
	player.Ign = leaguePlayer.Name
	player.Id = leaguePlayer.ID
	player.NormalizedIgn = goriot.NormalizeSummonerName(leaguePlayer.Name)[0]
	id := player.Id
	leagues_by_id, err := goriot.LeagueBySummoner("na", id)
	if err != nil {
		fmt.Println("wat: ", err.Error())
		player.Tier = "None"
	}
	league, ok := leagues_by_id[id]
	if ok {
		player.Tier = getBestLeague(league, player)
	}
	fmt.Println("New Player added: ", player)
	ols.GetPlayersDAO().Save(player)
	return &player
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

func nameErrors() {
	players := ols.GetPlayersDAO().All()
	for _, player := range players {
		_, err := goriot.SummonerByName("na", goriot.NormalizeSummonerName(player.NormalizedIgn)...)
		if err != nil {
			fmt.Println("Error with: ", player.Ign, " : ", err)
		}

	}
}

func nameUpdates() {
	players := ols.GetPlayersDAO().All()

	for _, player := range players {
		summoner, err := goriot.SummonerByID("na", player.Id)
		if err != nil {
			fmt.Println("Error with: ", player.Ign, " : ", err, player)
			continue
		}

		player.Ign = summoner[player.Id].Name
		ols.GetPlayersDAO().Save(*player)
	}
}

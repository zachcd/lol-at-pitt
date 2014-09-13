package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// const exist in go?
var PlayersCsvFile string = "players.csv"
var OutputFile string = "ols_players.csv"

type Player struct {
	Ign           string
	Id            string
	Name          string
	NormalizedIgn string
	Roles         []string
	Email         string
	Score         int
	Team          string
	Captain       bool
}

type Players []Player

func main() {
	//	players := Initialize() // magic.
	derp := LolGetSummoners("https://na.api.pvp.net/api/lol/na/v1.4/summoner/by-name/iph?api_key=a3c96054-e21f-4238-a842-28caa10943a0")
	fmt.Println("Summoner data: ", derp)

}

// Use this to rewrite the json ols player blobs
func WriteJsonPlayers(players []Player) {
	data, _ := json.MarshalIndent(players, "", "    ") //json pretty printer.
	os.Stdout.Write(data)                              // lazy redirect it for now.
}

func Initialize() []Player {
	ols_csv_file, err := os.Open(PlayersCsvFile)

	// Row number that the info is defined in. 0 indexed, so if it is player_name,lolkingscore, player_name_row = 0
	/*
		player_name_row := 0
		player_score_row := 0
		player_roles_row := 0
		player_ign_row := 0
		player_email_row := 0
	*/
	// This is pulled from the league API. Currently I just manually made a shell script
	// that creates the call and use CURL to populate it. Could be done in go..should be done in go.

	// Do it properly or explode.
	defer ols_csv_file.Close()

	// do it properly or really explode here.
	if err != nil {
		panic("wtf. supply the correct files asshole")
	}

	reader := csv.NewReader(ols_csv_file)

	// Inefficient for large files, but 140 players can fit into memory.
	// TODO: Figure out the "go" way for handling files in memory.
	all_data, _ := reader.ReadAll()

	players := []Player{}
	for _, row := range all_data {
		score, _ := strconv.Atoi(row[3]) // am I programming in C?
		player := Player{Ign: row[2], NormalizedIgn: NormalizedPlayerName(row[2]), Roles: PlayerRoles(row[4]), Score: score, Email: row[6], Name: row[1]}
		players = append(players, player)
	}

	/*
		var json_blob interface{}
		json_reader := json.NewDecoder(json_league_data)
		json_reader.UseNumber()
		json_reader.Decode(&json_blob)

		// Concatenate useful data from league data. Currently just IDs, but there are some other useful things..
		for i, player := range players {
			player_json := json_blob.(map[string]interface{})[player.NormalizedIgn]
			if player_json == nil {
				continue
			}
			id := player_json.(map[string]interface{})["id"]

			players[i].Id = id.(json.Number).String()
		}

		// Why does go not have programming specific filters? fuck.
		players = FilterPlayers(players)
		return players
	*/
	return Players{}
}

// FilterPlayers: Removes dirty players with no ids. There was an error in the league API calls, so only first 120 got recorded, otherwise I wouldnt have to.
// TODO: Create a tool to add players to the json blob on the fly, automatically.
func FilterPlayers(players []Player) []Player {
	new_players := []Player{}
	for _, player := range players {
		if len(player.Id) > 0 {
			new_players = append(new_players, player)
		}
	}
	return new_players
}

// NormalizedPlayerName:  Apparently na.op.gg and rito uses normalized ign names. This "normalizes" it.
func NormalizedPlayerName(name string) string {
	lowercase_name := strings.ToLower(name)
	normalized_name := strings.Replace(lowercase_name, " ", "", -1)
	return normalized_name

}

// Because functions.
func PlayerRoles(roles string) []string {
	split := strings.Split(roles, ", ")
	return split
}

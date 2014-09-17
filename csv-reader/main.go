package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const PlayersCsvFile string = "resources/players.csv"
const OutputFile string = "resources/ols_players.json"
const TeamCsvFile string = "resources/teams.csv"

type Player struct {
	Ign           string
	Id            int64
	Name          string
	NormalizedIgn string
	Roles         []string
	Score         int
	Team          string
	Captain       bool
}

type Players []*Player

func (p *Players) Filter(filter func(Player) bool) Players {
	players := Players{}
	for _, player := range *p {
		if filter(*player) {
			players = append(players, player)
		}
	}

	return players
}

func (p *Players) Print() {
	for _, player := range *p {
		fmt.Println(player)
	}
}

func main() {
	//	players := Initialize() // magic.
	players := Initialize()
	WriteJsonPlayers(players)

}

// Use this to rewrite the json ols player blobs
func WriteJsonPlayers(players Players) {
	data, _ := json.MarshalIndent(players, "", "    ") //json pretty printer.
	ioutil.WriteFile(OutputFile, data, 0644)

}

func Initialize() Players {
	ols_csv_file, err := os.Open(PlayersCsvFile)

	// Row number that the info is defined in. 0 indexed, so if it is player_name,lolkingscore, player_name_row = 0
	player_name_row := 0
	player_score_row := 2
	player_roles_row := 3
	player_ign_row := 1

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

	players := Players{}
	for i, row := range all_data {
		if i == 0 {
			continue
		}
		score, _ := strconv.Atoi(row[player_score_row]) // am I programming in C?

		player := Player{Ign: row[player_ign_row], NormalizedIgn: NormalizedPlayerName(row[player_ign_row]), Roles: PlayerRoles(row[player_roles_row]), Score: score, Name: row[player_name_row]}
		players = append(players, &player)
	}

	// filter players out that have no ign (probably blank in the csv
	players = players.Filter(func(player Player) bool {
		return player.Ign != ""
	})
	// Pull in league data.
	players = BuildSummonerData(&players)

	// Set up captains
	for _, player := range players {
		_, ok := Captains[player.NormalizedIgn]
		if ok {
			player.Captain = true
		}
	}
	// Pull in team data.
	LoadInTeam(players)

	return players
}

func LoadInTeam(players Players) {
	team_csv_file, err := os.Open(TeamCsvFile)

	// Do it properly or explode.
	defer team_csv_file.Close()

	// do it properly or really explode here.
	if err != nil {
		panic("wtf. supply the correct files asshole")
	}

	reader := csv.NewReader(team_csv_file)

	// Inefficient for large files, but 140 players can fit into memory.
	// TODO: Figure out the "go" way for handling files in memory.
	all_data, _ := reader.ReadAll()

	for _, team := range all_data {
		teamPlayers := team[2:]
		teamName := team[0]
		for _, teamPlayer := range teamPlayers {
			player := find(teamPlayer, players)
			if player != nil {
				player.Team = teamName
			}
		}

	}

	players.Print()

}

func find(teamPlayer string, players Players) *Player {
	for _, player := range players {
		if player.NormalizedIgn == NormalizedPlayerName(teamPlayer) {
			return player
		}
	}

	return nil
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

// Holds extra stuff that makes Go easier to use.
package utils

import (
	"encoding/base64"
	"fmt"
)

// Map types
const (
	SUMMONERS_RIFT   = "map11"
	HOWING_ABYSS     = "map12"
	CULL_STUFF       = "map8"
	TWISTED_TREELINE = "map10"
)

// Pick modes
const (
	ALL_RANDOM       = "pick4"
	TOURNAMENT_DRAFT = "pick6"
	DRAFT_MODE       = "pick2"
	BLIND_PICK       = "pick1"
)

// Spectator modes
const (
	SPECTATOR_ALL   = "specALL"
	SPECTATOR_LOBBY = "specLOBBYONLY"
	SPECTATOR_NONE  = "specNONE"
)

const pvp_format = "pvpnet://lol/customgame/joinorcreate/%s/%s/team%d/%s/%s"
const game_data = `{"name":"%s","password":"%s"}`

// Required input to create a Tournament code.
type TournamentCodeInput struct {
	SpectatorMode    string
	Mode             string // Refer to pick modes
	Map              string // Refer to Map Types comment
	Name             string
	OptionalPassword string
	TeamSize         int
}

/*
Takes a filled TournamentCodeInput and produces a tournament code for league.
*/
func GenerateTournamentCode(settings TournamentCodeInput) string {
	json_data := fmt.Sprintf(game_data, settings.Name, settings.OptionalPassword)
	output := base64.StdEncoding.EncodeToString([]byte(json_data))
	ret := fmt.Sprintf(pvp_format, settings.Map, settings.Mode, settings.TeamSize, settings.SpectatorMode, output)
	return ret
}

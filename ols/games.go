package ols

import (
	"github.com/TrevorSStone/goriot"
	"time"
)

type Game struct {
	GameId      int64
	PlayerStats map[int64]GameStat
	Blue        *Team
	Purple      *Team
	Winner      *Team
	Date        time.Time
	Official    bool // Some games might not be official
}

type GameStat struct {
	LeagueStats goriot.GameStat
	ChampionId  int
	Spell1      int
	Spell2      int
}

// Assumes you decided you applied the rules for a game.
func BuildGame(game goriot.Game) {
	ols_game := Game{}

	ols_game.GameId = game.GameID
	players := game.FellowPlayers
	_ = players
}

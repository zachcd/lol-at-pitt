package utils

import "testing"

func TestEncode(tst *testing.T) {
	tourneyData := TournamentCodeInput{
		SpectatorMode: SPECTATOR_ALL,
		Mode:          TOURNAMENT_DRAFT,
		Map:           SUMMONERS_RIFT,
		Name:          "Summoner hoes r us",
		TeamSize:      5,
	}

	correctData := "pvpnet://lol/customgame/joinorcreate/map11/pick6/team5/specALL/eyJuYW1lIjoiU3VtbW9uZXIgaG9lcyByIHVzIiwicGFzc3dvcmQiOiIifQ=="
	if GenerateTournamentCode(tourneyData) != correctData {
		tst.Error("Invalid Pvp code generated")
	}

}

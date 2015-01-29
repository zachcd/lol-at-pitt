package main

import (
	"github.com/TrevorSStone/goriot"
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/ols"
	"github.com/martini-contrib/render"
)

var fuck int = 310001
var championMap map[int]string

type PlayerStats struct {
	Division       string
	Team           string
	Summoner       string
	TotalKills     int
	TotalDeaths    int
	TotalAssists   int
	TotalGold      int
	GamesPlayed    int
	KillsPerGame   float32
	DeathsPerGame  float32
	AssistsPerGame float32
	KDA            float32
}

type MatchPlayer struct {
	Champion string
	Ign      string
}

type MatchDisplay struct {
	Blue []MatchPlayer
	Red  []MatchPlayer
}

var stats []PlayerStats

func initFunnyRouter(m *martini.ClassicMartini) {
	stats := initStats()
	m.Get("/player/stats", func(renderer render.Render) {
		renderer.JSON(200, stats)
	})

	m.Get("/stats", func(renderer render.Render) {
		renderer.HTML(200, "stats", stats)
	})

	m.Get("/fuck/smegs", func(renderer render.Render) {
		renderer.HTML(200, "fuck", fuck)
	})

	m.Get("/fuck/smeg/count", func(renderer render.Render) {
		fuck += 1
		renderer.Redirect("/fuck/smegs", 302)
	})

	m.Get("/matches/:team", func(renderer render.Render, params martini.Params) {
		matches := ols.GetMatchesDAO().LoadTeamMatches(params["team"])
		for _, match := range matches {
			player := map[int]MatchPlayer{}
			leagueGame := ols.GetMatchesDAO().LoadLeagueGame(match.Id)
			for _, participant := range leagueGame.Participants {
				player[participant.ParticipantID] = participant.ChampionID
			}
		}

	})

}
func initChampionMap() {
	champions, err := goriot.ChampionList("na", false)
	if err != nil {
		panic(err)
	}

	chMap := map[int]string{}
	for _, champion := range champions {
		chMap[champion.ID]
	}
}

func initStats() []PlayerStats {
	pStats := []PlayerStats{}

	for _, team := range ols.GetTeamsDAO().All() {
		matches := ols.GetMatchesDAO().LoadTeamMatches(team.Name)
		for _, playerId := range team.Players {
			pStat := PlayerStats{}
			player := ols.GetPlayersDAO().Load(playerId)
			pStat.Summoner = player.Ign
			pStat.Team = team.Name
			pStat.Division = team.League
			for _, match := range matches {
				leagueMatch := ols.GetMatchesDAO().LoadLeagueGame(match.Id)
				index := getParticipantIndex(leagueMatch, *match, player.Id)
				if index == -1 {
					continue
				}

				pleague := leagueMatch.Participants[index].Stats
				pStat.AssistsPerGame += float32(pleague.Assists)
				pStat.DeathsPerGame += float32(pleague.Deaths)
				pStat.GamesPlayed += 1
				pStat.KillsPerGame += float32(pleague.Kills)
				pStat.TotalAssists += int(pleague.Assists)
				pStat.TotalDeaths += int(pleague.Deaths)
				pStat.TotalGold += int(pleague.GoldEarned)
				pStat.TotalKills += int(pleague.Kills)
			}
			if pStat.GamesPlayed > 0 {
				pStat.AssistsPerGame /= float32(pStat.GamesPlayed)
				pStat.DeathsPerGame /= float32(pStat.GamesPlayed)
				pStat.KillsPerGame /= float32(pStat.GamesPlayed)
				pStat.KDA = float32(pStat.TotalAssists+pStat.TotalKills) / float32(pStat.TotalDeaths)
			}
			pStats = append(pStats, pStat)
		}

	}

	return pStats
}

func getParticipantIndex(match goriot.MatchDetail, olsMatch ols.Match, id int64) int {
	pId := 0
	for _, participant := range olsMatch.Participants {
		if participant.Id == id {
			pId = participant.ParticipantId
		}
	}

	for i, participant := range match.Participants {
		if participant.ParticipantID == pId {
			return i
		}
	}

	return -1
}

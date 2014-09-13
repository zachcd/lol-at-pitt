package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Summoner struct {
	Id            int64
	Name          string
	ProfileIconId int
	// rest of the shit I dont care about like revisionDate and summonerLevel
}

type SummonerRequest map[string]Summoner

// Adds league data to players asynchronously. Sends them back through output
// A lot of this code is to prevent league from killing me for sending over the max
// requests (10 per second)
func BuildSummonerData(players *Players) Players {
	// you can build 20 summoner queries at a time.

	playerUrlChan := make(chan string)
	summonersChan := make(chan SummonerRequest)

	go DataDaemon(playerUrlChan, summonersChan, ApiKey) // ApiKey is stored in a .go file that JUST has the API key.
	playerUrl := ""
	for i, player := range *players {
		if i%20 == 0 && i != 0 {
			playerUrlChan <- strings.TrimRight(playerUrl, ",") // removes ending comma
			playerUrl = ""
			summoners := <-summonersChan
			OrganizeData(players, summoners)
		}

		playerUrl = playerUrl + player.NormalizedIgn + "," // inefficient, but whatever
	}

	playerUrlChan <- strings.TrimRight(playerUrl, ",")
	summoners := <-summonersChan
	OrganizeData(players, summoners)
	close(playerUrlChan)
	return *players

}

func OrganizeData(players *Players, summoners SummonerRequest) {
	for _, player := range *players {
		summoner, ok := summoners[player.NormalizedIgn]
		if ok {
			player.Id = summoner.Id
		}
	}
}

// Why not..
func DataDaemon(playerUrlChan chan string, summoners chan SummonerRequest, apikey string) {
	for {
		playerUrl, ok := <-playerUrlChan

		if !ok {
			return // We are done here, wrap it up!
		}

		lolUrl := "https://na.api.pvp.net/api/lol/na/v1.4/summoner/by-name/" + playerUrl + "?api_key=" + apikey
		summons := LolGetSummoners(lolUrl)
		summoners <- summons
		time.Sleep(time.Second / 10.0) // Can make 10 requests per second
	}

}
func LolGetSummoners(url string) SummonerRequest {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		var data SummonerRequest
		json.Unmarshal(body, &data)
		return data
	}
}

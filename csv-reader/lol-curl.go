package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
func BuildSummonerData(player Players, output chan Players) {
	// you can build 20 summoner queries at a time.

	playerUrlChan := make(chan string)
	summoners := make(chan SummonerRequest)

	go DataDaemon(playerUrlChan, summoners, ApiKey) // ApiKey is stored in a gofile that JUST has the API key.

}

// Why not..
func DataDaemon(playerUrlChan chan string, summoners chan SummonerRequest) {
	for {
		playerUrl, ok := <-playerUrlChan

		if !ok {
			return // We are done here, wrap it up!
		}

		lolUrl := "https://na.api.pvp.net/api/lol/na/v1.4/summoner/by-name/" + playerUrl
		time.Sleep(time.Second)
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

package main

import (
	"fmt"
	"log"
)

var (
	_REGION       = "na"
	_GET_SUMMONER = "https://na.api.pvp.net/api/lol/%s/v1.4/summoner/by-name/%s"
	_GET_RECENT   = "https://na.api.pvp.net/api/lol/%s/v1.3/game/by-summoner/%d/recent"
)

var _SEED_SUMMONERS = []string{
	"colonelxc",
	"jeneraljam",
	"jon9890",
	"slipmthgoose",
}

var (
	toCrawl = make([]int64)
)

func doGet(url, token string, dst *interface{}) error {
	// TODO: Implement, handle rate limiting

	return nil
}

func main() {
	// TODO: Get token from command line

	// Find the summoner IDs for the seed summoners
	var summoners = make(map[string]Summoner)
	url := fmt.Sprintf(_GET_SUMMONER, _REGION, strings.join(_SEED_SUMMONERS, ","))
	err := doGet(url, "", summoners)
	if err != nil {
		log.Printf("Unable to fetch seed summoners: %s", err.Error())
		return
	}

	// Add seed summoners to toCrawl
	for _, summoner := range summoners {
		toCrawl = append(toCrawl, summoner.ID)
	}

	// Loop until done crawling
	for len(toCrawl) > 0 {
		// Pop a summoner off the list of people to crawl
		summoner, toCrawl = toCrawl[0], toCrawl[1:]

		// Look up the summoner's recent games
		var games = RecentGames{}
		url := fmt.Sprintf(_GET_RECENT, _REGION, summoner)
		err := doGet(url, "", games)
		if err != nil {
			log.Printf("Unable to fetch recent games for summoner: %d -- %s", summoner, err.Error())
			return
		}

		for _, game := range games {
			// TODO: Store the game

			for _, player := range game.FellowPlayers {
				// TODO: Check if we have already crawled summoner ID

				// Haven't crawled summoner, add to toCrawl
				toCrawl = append(toCrawl, player.SummonerID)
			}
		}
	}

	log.Printf("Done crawling for now")
}

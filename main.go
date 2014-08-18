package main

import "log"

var _SEED_SUMMONERS = []string{
	"colonelxc",
	"jeneraljam",
	"jon9890",
	"slipmthgoose",
}

var (
	toCrawl = make([]int64, 100)
)

func main() {
	// TODO: Get token, rate limits from command line
	c := crawler{Token: "",
		RateLimitPerMinute: 60,
		RateLimitPerHour:   500}

	// Find the summoner IDs for the seed summoners
	summoners, err := c.getSummoners(_SEED_SUMMONERS)
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
		summoner, toCrawl := toCrawl[0], toCrawl[1:]

		// Look up the summoner's recent games
		games, err := c.getRecentGames(summoner)
		if err != nil {
			log.Printf("Unable to fetch recent games for summoner: %d -- %s", summoner, err.Error())
			return
		}

		for _, game := range games.Games {
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

package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	_REGION           = "na"
	_GET_SUMMONER     = "https://na.api.pvp.net/api/lol/%s/v1.4/summoner/by-name/%s"
	_GET_RECENT_GAMES = "https://na.api.pvp.net/api/lol/%s/v1.3/game/by-summoner/%d/recent"
)

type crawler struct {
	Token              string
	Requests           []int64
	RateLimitPerMinute uint
	RateLimitPerHour   uint
}

// Block until rate limit is not exceeded
func (c *crawler) rateLimit() {
	// TODO: handle rate limiting
}

func (c *crawler) fetchResource(url string, dst interface{}) error {
	// Block until we are able to make a request
	c.rateLimit()

	// TODO: Implement fetch and decode the JSON

	return nil
}

// Lookup summoners by their summoner name. A maximum of 40 summoners is
// allowed at one time.
func (c *crawler) getSummoners(summoners []string) (map[string]Summoner, error) {
	if len(summoners) > 40 {
		return nil, errors.New("Maximum of 40 summoners per query")
	}

	var res = make(map[string]Summoner)

	url := fmt.Sprintf(_GET_SUMMONER, _REGION, strings.Join(_SEED_SUMMONERS, ","))
	err := c.fetchResource(url, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup the recent games for a given summoner ID. Returns at most 10 games.
func (c *crawler) getRecentGames(id int64) (*RecentGames, error) {
	games := &RecentGames{}

	url := fmt.Sprintf(_GET_RECENT_GAMES, _REGION, id)
	err := c.fetchResource(url, games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

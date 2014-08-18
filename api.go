package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// North America, TODO: make configurable
	_Region = "na"

	// HTTP Status code for Rate Limit Exceeded.
	_StatusRateLimitExceeded = 429

	// API URLs, ready for fmt.Sprintf
	_GetSummoner    = "https://na.api.pvp.net/api/lol/%s/v1.4/summoner/by-name/%s"
	_GetRecentGames = "https://na.api.pvp.net/api/lol/%s/v1.3/game/by-summoner/%d/recent"
)

// Struct for maintaining crawler state including the timestamps of recent
// requests so that we can perform rate limiting.
type crawler struct {
	Token                  string       // API Key for authentication
	RateLimitPerTenSeconds int          // Maximum number of requests per ten seconds
	RateLimitPerHour       int          // Maximum number of requests per hour
	MaxRetries             int          // Maximum number of times to retry a request
	Client                 *http.Client // Client for making requests
	Requests               []time.Time  // Timestamps for the most recent requests
}

// Block until rate limit is not exceeded
func (c *crawler) rateLimit() {
	// Get the current time in
	now := time.Now()
	minusHour := now.Add(-1 * time.Hour)
	minusSeconds := now.Add(-10 * time.Second)

	// Time to sleep for, max of rate limited for per-ten-second and per-hour
	// requests.
	var sleep time.Duration

	// Prune off requests that are more than one hour old
	for i, t := range c.Requests {
		if minusHour.Before(t) {
			c.Requests = c.Requests[i:]
			break
		}
	}

	// Count the number of requests in the last hour
	if len(c.Requests) >= c.RateLimitPerHour {
		// Sleep until oldest request plus one hour has passed
		sleep = c.Requests[0].Add(time.Hour).Sub(now)
	}

	// Count the number of requests in the last second
	for i, t := range c.Requests {
		if minusSeconds.Before(t) {
			// Count the number of requests in the last ten seconds
			if len(c.Requests)-i >= c.RateLimitPerTenSeconds {
				// Sleep until oldest request in the last second plus one second has
				// passed. Only update sleep if it's longer than the previously
				// computed sleep time.
				d := c.Requests[i].Add(10 * time.Second).Sub(now)
				if d > sleep {
					sleep = d
				}
			}

			break
		}
	}

	// Sleep for the predetermined amount of time
	time.Sleep(sleep)

	// Add current time to list of requests
	c.Requests = append(c.Requests, time.Now())
}

// Take a base URL and add query parameters from map to the end of the URL. If
// a query parameter is already set in the base URL, it will be overwritten.
func buildURL(base string, params map[string]string) (string, error) {
	// Parse the URL
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	// Add the query parameters
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// Return the updated URL
	return u.String(), nil
}

func (c *crawler) fetchResource(url string, dst interface{}) error {
	var err error

	// Add API key to the base URL
	u, err := buildURL(url, map[string]string{"api_key": c.Token})
	if err != nil {
		return err
	}

	retries := 0

	for ; retries < c.MaxRetries; retries++ {
		if retries != 0 {
			// Didn't succeed on previous attempt, sleep for a bit (in addition to the
			// rate limiting) before trying again.
			time.Sleep(10 * time.Second)
		}

		// Block until we are able to make a request
		c.rateLimit()

		resp, err := c.Client.Get(u)
		if err != nil {
			log.Printf("Failed to request resource. Sleeping and then retrying.")
			continue
		}

		if resp.StatusCode == _StatusRateLimitExceeded {
			log.Printf("Rate limit exceeded. Sleeping and then retrying.")
			continue
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code from server: %s", resp.Status)
		} else {
			// Got a valid response
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response body. Sleeping and then retrying.")
				continue
			}

			if err := json.Unmarshal(body, dst); err != nil {
				log.Printf("Failed to unmarshal response body. Sleeping and then retrying.")
				continue
			} else {
				// Success!
				return nil
			}
		}
	}

	return fmt.Errorf("max retries exceeded for API request (last error: %s)", err.Error())
}

// Lookup summoners by their summoner name. A maximum of 40 summoners is
// allowed at one time.
func (c *crawler) getSummoners(summoners []string) (map[string]Summoner, error) {
	if len(summoners) > 40 {
		return nil, errors.New("maximum of 40 summoners per query")
	}

	log.Printf("Fetching summoners: %v", summoners)

	var res = make(map[string]Summoner)

	url := fmt.Sprintf(_GetSummoner, _Region, strings.Join(_SEED_SUMMONERS, ","))
	err := c.fetchResource(url, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup the recent games for a given summoner ID. Returns at most 10 games.
func (c *crawler) getRecentGames(id int64) (*RecentGames, error) {
	games := &RecentGames{}

	log.Printf("Fetching recent games for summoner: %d", id)

	url := fmt.Sprintf(_GetRecentGames, _Region, id)
	err := c.fetchResource(url, games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

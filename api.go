package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// North America, TODO: make configurable
	_Region = "na"

	// HTTP Status code for Rate Limit Exceeded.
	_StatusRateLimitExceeded = 429

	// API URLs, ready for fmt.Sprintf
	_GetSummoner     = "https://na.api.pvp.net/api/lol/%s/v1.4/summoner/by-name/%s"
	_GetSummonerByID = "https://na.api.pvp.net/api/lol/%s/v1.4/summoner/%s"
	_GetRecentGames  = "https://na.api.pvp.net/api/lol/%s/v1.3/game/by-summoner/%d/recent"

	_MaxSummonersPerQuery = 40
)

// Struct for maintaining crawler state including the timestamps of recent
// requests so that we can perform rate limiting.
type crawler struct {
	Token                  string       // API Key for authentication
	RateLimitPerTenSeconds int          // Maximum number of requests per ten seconds
	RateLimitPerTenMinutes int          // Maximum number of requests per ten minutes
	MaxRetries             int          // Maximum number of times to retry a request
	Client                 *http.Client // Client for making requests
	Requests               []time.Time  // Timestamps for the most recent requests
}

func newCrawler(token string, rateLimitPerTenSeconds, rateLimitPerTenMinutes, maxRetries int) *crawler {
	return &crawler{
		Token: token,
		RateLimitPerTenSeconds: rateLimitPerTenSeconds,
		RateLimitPerTenMinutes: rateLimitPerTenMinutes,
		MaxRetries:             maxRetries,
		Client:                 &http.Client{},
		Requests:               make([]time.Time, rateLimitPerTenMinutes)}

}

// Block until rate limit is not exceeded
func (c *crawler) rateLimit() {
	// Get the current time
	now := time.Now()

	// Calculate rate limit intervals
	minusTenMins := now.Add(-10 * time.Minute)
	minusTenSecs := now.Add(-10 * time.Second)

	// Time to sleep for, max of rate limited for per-ten-second and
	// per-ten-minutes requests.
	var sleep time.Duration = 0

	// Prune off requests that are more than ten minutes old
	for i, t := range c.Requests {
		if minusTenMins.Before(t) {
			c.Requests = c.Requests[i:]
			break
		}
	}

	// Count the number of requests in the last ten minutes
	if len(c.Requests) >= c.RateLimitPerTenMinutes {
		// Sleep until oldest request plus ten minutes has passed
		sleep = c.Requests[0].Add(10*time.Minute + 30*time.Second).Sub(now)
	}

	// Count the number of requests in the last second
	for i, t := range c.Requests {
		if minusTenSecs.Before(t) {
			// Count the number of requests in the last ten seconds
			if len(c.Requests)-i >= c.RateLimitPerTenSeconds {
				// Sleep until oldest request in the last second plus one second has
				// passed. Only update sleep if it's longer than the previously
				// computed sleep time.
				d := c.Requests[i].Add(11 * time.Second).Sub(now)
				if d > sleep {
					sleep = d
				}
			}

			break
		}
	}

	//log.Printf("Rate limiting, sleeping for %f seconds", sleep.Seconds())

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

	for retries := 0; retries < c.MaxRetries; retries++ {
		//log.Printf("Attempting to get URL: %s, retries = %d", u, retries)
		if retries != 0 {
			// Didn't succeed on previous attempt, sleep for a bit (in addition to the
			// rate limiting) before trying again.
			time.Sleep(10 * time.Second)
		}

		// Block until we are able to make a request
		c.rateLimit()

		resp, err2 := c.Client.Get(u)
		if err2 != nil {
			err = err2
			log.Printf("Failed to request resource. Sleeping and then retrying.")
			continue
		}

		if resp.StatusCode == _StatusRateLimitExceeded {
			log.Printf("Rate limit exceeded. Sleeping and then retrying.")
		} else if resp.StatusCode == http.StatusOK {
			// Got a valid response
			defer resp.Body.Close()

			body, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				err = err2
				log.Printf("Failed to read response body. Sleeping and then retrying.")
				continue
			}

			if err2 := json.Unmarshal(body, dst); err2 != nil {
				err = err2
				log.Printf("Failed to unmarshal response body. Sleeping and then retrying.")
			} else {
				return nil // Success!
			}
		} else {
			err = fmt.Errorf("unexpected status code from server: %s", resp.Status)
			log.Printf(err.Error())
		}
	}

	if err != nil {
		return fmt.Errorf("max retries exceeded for API request (last error: %s)", err.Error())
	}

	return errors.New("unknown error occured while fetching resource")
}

// Lookup summoners by their summoner name. A maximum of _MaxSummonersPerQuery
// summoners is allowed at one time.
func (c *crawler) getSummoners(summoners []string) (map[string]Summoner, error) {
	return c.getSummonersHelper(_GetSummoner, strings.Join(summoners, ","))
}

// Lookup summoners by their summoner ID. A maximum of _MaxSummonersPerQuery summoners is allowed
// at one time.
func (c *crawler) getSummonersByID(ids []int64) (map[string]Summoner, error) {
	// Convert ids to strings so we can concat them together
	s := make([]string, 0, len(ids))
	for _, id := range ids {
		s = append(s, strconv.FormatInt(id, 10))
	}

	return c.getSummonersHelper(_GetSummonerByID, strings.Join(s, ","))
}

func (c *crawler) getSummonersHelper(url, summoners string) (map[string]Summoner, error) {
	if strings.Count(summoners, ",") > _MaxSummonersPerQuery {
		return nil, errors.New("exceeded maximum number of summoners per query")
	}

	//log.Printf("Fetching summoners: %s", summoners)

	var res = make(map[string]Summoner)

	url = fmt.Sprintf(url, _Region, summoners)
	if err := c.fetchResource(url, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup the recent games for a given summoner ID. Returns at most 10 games.
func (c *crawler) getRecentGames(id int64) (*RecentGames, error) {
	games := &RecentGames{}

	//log.Printf("Fetching recent games for summoner: %d", id)

	url := fmt.Sprintf(_GetRecentGames, _Region, id)
	err := c.fetchResource(url, games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

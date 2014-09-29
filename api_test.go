package main

import (
	"os"
	"testing"
)

var c *crawler

var (
	_TestSummoners   = []string{"Turtle the Cat", "Pobelter", "AAAltec"}
	_TestSummonerIDs = []int64{18991200, 2648, 42762376}
)

func init() {
	t := os.Getenv("TOKEN")
	if t != "" {
		c = newCrawler(t, int(*rateLimitPerTenSeconds),
			int(*rateLimitPerTenMinutes), int(*maxRetries))
	}
}

func TestGetSummoners(t *testing.T) {
	if c == nil {
		t.Fatal("TOKEN is not set")
	}

	res, err := c.getSummoners(_TestSummoners)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	checkTestSummoners(t, res)
}

func TestGetSummonersByID(t *testing.T) {
	if c == nil {
		t.Fatal("TOKEN is not set")
	}

	res, err := c.getSummonersByID(_TestSummonerIDs)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	checkTestSummoners(t, res)
}

func checkTestSummoners(t *testing.T, res map[string]Summoner) {
	// Make sure we got the right number
	if len(res) != len(_TestSummonerIDs) {
		t.Fatalf("Expected %d summoners, got %d", len(_TestSummonerIDs), len(res))
	}

	// Make sure that all the summoners are the ones we queried for and that the
	// mapping from summoner name -> ID is right.
	for _, summoner := range res {
		found := false

		for i, name := range _TestSummoners {
			if name == summoner.Name {
				found = true

				if _TestSummonerIDs[i] != summoner.Id {
					t.Errorf("Summoner ID does not match expect ID (%d != %d) for '%s'",
						summoner.Id, _TestSummonerIDs[i], name)
				}
			}
		}

		if !found {
			t.Errorf("Unxpected summoner '%s' in results", summoner.Name)
		}
	}
}

func TestGetMatch(t *testing.T) {
	if c == nil {
		t.Fatal("TOKEN is not set")
	}

	res, err := c.getMatch(int64(1560034527))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	t.Log("%#v", res)
}

func TestGetMatchHistory(t *testing.T) {
	if c == nil {
		t.Fatal("TOKEN is not set")
	}

	res, err := c.getMatchHistory(int64(18991200), 0)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	t.Log("%#v", res)
}

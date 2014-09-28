package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jcrussell/crawlol/external/github.com/jinzhu/gorm"

	// Import for side-effect
	_ "github.com/jcrussell/crawlol/external/github.com/mattn/go-sqlite3"
)

var (
	rateLimitPerTenSeconds = flag.Uint("rate-per-ten-seconds", 10, "API Rate limit per ten seconds")
	rateLimitPerTenMinutes = flag.Uint("rate-per-ten-minutes", 500, "API Rate limit per ten minutes")
	maxRetries             = flag.Uint("max-retries", 2, "Maximum number of times to retry a request")
	dbPath                 = flag.String("db", "crawlol.db", "Location for the SQLite database")
	seedSummoners          = flag.String("seed", "", "List of summoner names (separated by ',') to use to seed the database")
	newDatabase            = flag.Bool("init-db", false, "Specify that the database is new and tables should be created")
)

func crawl(db gorm.DB, c *crawler) {
	for {
		// Only recrawl summoners once every 12 hours
		lastCrawled := time.Now().Add(-12 * time.Hour)

		var summoners []Summoner
		db.Where("last_crawled < ?", lastCrawled.UnixNano()).Order("last_crawled asc").Limit(1000).Find(&summoners)

		if len(summoners) == 0 {
			log.Printf("All known summoners have been crawled in the last 12 hours")
			break
		}

		log.Printf("Crawling recent games for %d summoners", len(summoners))

		newSummoners := make(map[int64]bool, 100)
		newMatches := 0

		for _, summoner := range summoners {
			// Look up the summoner's match history, keep trying to get more until we
			// find no new matches.
			start := int64(0)

			for {
				foundNewMatches := false

				// Query for the summoner's match history
				matches, err := c.getMatchHistory(summoner.Id, start)
				if err != nil {
					log.Printf("Unable to fetch recent matches for summoner: %d (start = %d) -- %s", summoner, start, err.Error())
					break
				}

				// Process the matches found
				for _, match := range matches {
					// Check if we've seen this match before
					if db.Where(&MatchDetail{MatchID: match.MatchID}).First(&MatchDetail{}).RecordNotFound() {
						// We haven't so we should note to try to get more matches for this
						// summoner.
						foundNewMatches = true
						newMatches++

						// Get the actual details for the match
						details, err := c.getMatch(match.MatchID)
						if err != nil {
							log.Printf("Unable to fetch match details: %d -- %s", match.MatchID, err.Error())
							continue
						}

						// Save the match details
						if err := db.Create(&details).Error; err != nil {
							log.Printf("Unable to save match details: %d -- summoner: %d", match.MatchID, summoner.Id)
						}

						// Finally, process the players to find new summoners
						for _, identity := range details.ParticipantIdentities {
							summonerID := identity.Player.SummonerID

							// Check if we have already crawled summoner ID
							if db.Where(&Summoner{Id: summonerID}).First(&Summoner{}).RecordNotFound() {
								// Haven't crawled summoner, add to newSummoners
								newSummoners[summonerID] = true
							}
						}
					}
				}

				// Finished one batch of 15, move on to the next if there were new
				// matches.
				start += 15

				if !foundNewMatches {
					break
				}
			}

			// Finished with summoner, for now. Update the last crawled field and save
			// to the database.
			summoner.LastCrawled = time.Now().UnixNano()
			if err := db.Save(&summoner).Error; err != nil {
				log.Printf("Unable to update last crawled for summoner: %d -- %s", summoner.Id, err.Error)
			}
		}

		log.Printf("Crawled %d new matches and found %d new summoners", newMatches, len(newSummoners))
		lookupSummoners(db, c, newSummoners)
	}
}

func initDB(db gorm.DB) {
	// TODO: Error checking?

	/*
	TODO: Figure out how to store MatchDetail
	// Create summoner table
	db.CreateTable(Summoner{})
	db.Model(Summoner{}).AddUniqueIndex("idx_name", "name")
	db.Model(Summoner{}).AddIndex("idx_last_crawled", "last_crawled")

	// Create game table
	db.CreateTable(Game{})
	db.Model(Game{}).AddUniqueIndex("idx_game_summoner_id", "summoner_id", "game_id")
	db.Model(Game{}).AddIndex("idx_game_id", "game_id")

	// Create raw stats table
	db.CreateTable(RawStats{})
	*/
}

// Seed the database with summoners looked up by name.
func seedDatabase(db gorm.DB, c *crawler) {
	names := strings.Split(*seedSummoners, ",")

	// Find the summoner IDs for the seed summoners
	if summoners, err := c.getSummoners(names); err != nil {
		log.Fatalf("Unable to fetch seed summoners: %s", err.Error())
	} else {
		saveSummoners(db, summoners)
	}
}

// Save new summoners to the database, set the last crawled time as never.
func saveSummoners(db gorm.DB, summoners map[string]Summoner) {
	for _, summoner := range summoners {
		summoner.LastCrawled = 0
		if err := db.Create(&summoner).Error; err != nil {
			log.Printf("Unable to save summoner: %d -- %s", summoner.Id, err.Error())
		}
	}
}

func lookupSummoners(db gorm.DB, c *crawler, summoners map[int64]bool) {
	ids := make([]int64, 0, len(summoners))
	for k := range summoners {
		ids = append(ids, k)
	}

	for i := 0; i*_MaxSummonersPerQuery < len(ids); i++ {
		// Create slice of length _MaxSummonersPerQuery or less if there aren't
		// enough elements remaining.
		ub := (i + 1) * _MaxSummonersPerQuery
		if ub >= len(ids) {
			ub = len(ids) - 1
		}
		slice := ids[i*_MaxSummonersPerQuery : ub]

		if res, err := c.getSummonersByID(slice); err != nil {
			log.Printf("Unable to fetch summoners: %s", err.Error())
		} else {
			saveSummoners(db, res)
		}
	}
}

type fakeLogger struct{}

func (fakeLogger) Print(v ...interface{}) {}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("USAGE: %s [OPTIONS] TOKEN\n", os.Args[0])
		os.Exit(1)
	}

	c := newCrawler(flag.Arg(0), int(*rateLimitPerTenSeconds),
		int(*rateLimitPerTenMinutes), int(*maxRetries))

	db, err := gorm.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Unable to open database: %s", err.Error())
	}

	// TODO: Remove, only here to test that we can fetch a match
	match, err := c.getMatch(int64(1560034527))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("%#v", match)
	}

	os.Exit(0)

	// Disable Gorm's logging
	db.SetLogger(fakeLogger{})

	if *newDatabase {
		initDB(db)
	}

	if *seedSummoners != "" {
		seedDatabase(db, c)
	}

	crawl(db, c)

	log.Printf("Done crawling for now")
}

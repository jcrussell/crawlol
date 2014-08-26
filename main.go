package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jcrussell/crawlol/external/github.com/jinzhu/gorm"
)

var (
	rateLimitPerTenSeconds = flag.Uint("rate-per-ten-seconds", 10, "API Rate limit per ten seconds")
	rateLimitPerHour       = flag.Uint("rate-per-hour", 500, "API Rate limit per hour")
	maxRetries             = flag.Uint("max-retries", 2, "Maximum number of times to retry a request")
	dbPath                 = flag.String("db", "crawlol.db", "Location for the SQLite database")
	seedSummoners          = flag.String("seed", "", "List of summoner names (separated by ':') to use to seed the database")
	newDatabase            = flag.Bool("init-db", false, "Specify that the database is new and tables should be created")
)

func crawl(db gorm.DB, c crawler) {
	for {
		var summoners []Summoner

		// Only recrawl summoners once every 12 hours
		lastCrawled := time.Now().Add(-12 * time.Hour)

		db.Where("last_crawled < ?", lastCrawled.UnixNano()).Limit(100).Find(&summoners)

		if len(summoners) == 0 {
			log.Printf("All known summoners have been crawled in the last 12 hours")
			break
		}

		newSummoners := make([]int64, 100)

		for _, summoner := range summoners {
			// Look up the summoner's recent games
			recent, err := c.getRecentGames(summoner.ID)
			if err != nil {
				log.Printf("Unable to fetch recent games for summoner: %d -- %s", summoner, err.Error())
				continue
			}

			for _, game := range recent.Games {
				game.SummonerID = summoner.ID
				// TODO: Check errors. Need to save Stats or done automatically?
				db.Save(&game)
				db.Save(&game.Stats)

				for _, player := range game.FellowPlayers {
					// Check if we have already crawled summoner ID
					count := 0
					db.Where(&Summoner{ID: player.SummonerID}).Count(&count)

					// Haven't crawled summoner, add to newSummoners
					if count == 0 {
						newSummoners = append(newSummoners, player.SummonerID)
					}
				}
			}
		}

		for i := 0; i*_MaxSummonersPerQuery < len(newSummoners); i++ {
			ids := newSummoners[i*_MaxSummonersPerQuery : (i+1)*_MaxSummonersPerQuery]

			if summoners, err := c.getSummonersByID(ids); err != nil {
				log.Printf("Unable to fetch summoners: %s", err.Error())
			} else {
				saveSummoners(db, summoners)
			}
		}
	}
}

func initDB(db gorm.DB) {
	// TODO: Error checking?

	// Create summoner table
	db.CreateTable(Summoner{})
	db.Model(Summoner{}).AddUniqueIndex("idx_name", "name")
	db.Model(Summoner{}).AddIndex("idx_last_crawled", "last_crawled")

	// Create game table
	db.CreateTable(Game{})
	db.Model(Game{}).AddUniqueIndex("idx_game_summoner_id", "summoner", "game_id")
	db.Model(Game{}).AddIndex("idx_game_id", "game_id")

	// Create raw stats table
	db.CreateTable(RawStats{})
}

// Seed the database with summoners looked up by name.
func seedDatabase(db gorm.DB, c crawler) {
	names := strings.Split(*seedSummoners, ":")

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
		if err := db.Save(&summoner).Error; err != nil {
			log.Printf("Unable to save summoner: %d -- %s", summoner.ID, err.Error())
		}
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("USAGE: %s [OPTIONS] TOKEN", os.Args[0])
	}

	c := crawler{Token: flag.Arg(0),
		RateLimitPerTenSeconds: int(*rateLimitPerTenSeconds),
		RateLimitPerHour:       int(*rateLimitPerHour)}

	db, err := gorm.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Unable to open database: %s", err.Error())
	}

	if *newDatabase {
		initDB(db)
	}

	if *seedSummoners != "" {
		seedDatabase(db, c)
	}

	crawl(db, c)

	log.Printf("Done crawling for now")
}

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
		db.Where("last_crawled < ?", lastCrawled.UnixNano()).Limit(1000).Find(&summoners)

		if len(summoners) == 0 {
			log.Printf("All known summoners have been crawled in the last 12 hours")
			break
		}

		log.Printf("Crawling recent games for %d summoners", len(summoners))

		newSummoners := make(map[int64]bool, 100)

		for _, summoner := range summoners {
			// Look up the summoner's recent games
			recent, err := c.getRecentGames(summoner.Id)
			if err != nil {
				log.Printf("Unable to fetch recent games for summoner: %d -- %s", summoner, err.Error())
				continue
			}

			for _, game := range recent.Games {
				// TODO: Check errors.
				//log.Printf("Processing Game: %d", game.GameId)

				// Check if we have already stored this game before
				if db.Where(&Game{GameId: game.GameId, SummonerId: summoner.Id}).First(&Game{}).RecordNotFound() {
					// Save new game
					game.SummonerId = summoner.Id
					db.Create(&game)

					for _, player := range game.FellowPlayers {
						//log.Printf("Processing Player: %d", player.SummonerId)
						// Check if we have already crawled summoner ID
						if db.Where(&Summoner{Id: player.SummonerId}).First(&Summoner{}).RecordNotFound() {
							// Haven't crawled summoner, add to newSummoners
							newSummoners[player.SummonerId] = true
						}
					}
				}
			}
		}

		log.Printf("Found %d new summoners to crawl", len(newSummoners))

		ids := make([]int64, 0)

		for id := range newSummoners {
			ids = append(ids, id)

			if len(ids) == _MaxSummonersPerQuery {
				lookupSummoners(db, c, ids)
				ids = make([]int64, 0)
			}
		}

		if len(ids) > 0 {
			lookupSummoners(db, c, ids)
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
	db.Model(Game{}).AddUniqueIndex("idx_game_summoner_id", "summoner_id", "game_id")
	db.Model(Game{}).AddIndex("idx_game_id", "game_id")

	// Create raw stats table
	db.CreateTable(RawStats{})
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
		//log.Printf("Save summoner: %s (Id = %d)", summoner.Name, summoner.Id)
		summoner.LastCrawled = 0
		if err := db.Create(&summoner).Error; err != nil {
			log.Printf("Unable to save summoner: %d -- %s", summoner.Id, err.Error())
		}
	}
}

func lookupSummoners(db gorm.DB, c *crawler, ids []int64) {
	if summoners, err := c.getSummonersByID(ids); err != nil {
		log.Printf("Unable to fetch summoners: %s", err.Error())
	} else {
		saveSummoners(db, summoners)
	}
}

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

	if *newDatabase {
		initDB(db)
	}

	if *seedSummoners != "" {
		seedDatabase(db, c)
	}

	crawl(db, c)

	log.Printf("Done crawling for now")
}

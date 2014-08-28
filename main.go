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
		newGames := 0

		for _, summoner := range summoners {
			// Look up the summoner's recent games
			recent, err := c.getRecentGames(summoner.Id)
			if err != nil {
				log.Printf("Unable to fetch recent games for summoner: %d -- %s", summoner, err.Error())
				continue
			}

			for _, game := range recent.Games {
				// Check if we have already stored this game before
				if db.Where(&Game{GameId: game.GameId, SummonerId: summoner.Id}).First(&Game{}).RecordNotFound() {
					newGames++

					// Save new game
					game.SummonerId = summoner.Id
					if err := db.Create(&game).Error; err != nil {
						log.Printf("Unable to save game: %d -- summoner: %d", game.GameId, summoner.Id)
					}

					// Process the players that also played in this game
					for _, player := range game.FellowPlayers {
						// Check if we have already crawled summoner ID
						if db.Where(&Summoner{Id: player.SummonerId}).First(&Summoner{}).RecordNotFound() {
							// Haven't crawled summoner, add to newSummoners
							newSummoners[player.SummonerId] = true
						}
					}
				}
			}

			// Update the last crawled time to now for summoner
			summoner.LastCrawled = time.Now().UnixNano()
			if err := db.Save(&summoner).Error; err != nil {
				log.Printf("Unable to update last crawled for summoner: %d -- %s", summoner.Id, err.Error)
			}
		}

		log.Printf("Crawled %d games and found %d new summoners", newGames, len(newSummoners))
		lookupSummoners(db, c, newSummoners)
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

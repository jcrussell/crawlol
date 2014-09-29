package main

import (
	"fmt"
	"log"

	"github.com/jcrussell/crawlol/external/github.com/jinzhu/gorm"
)

type fakeLogger struct{}

func (fakeLogger) Print(v ...interface{}) {}

func openDB(path string) (gorm.DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return gorm.DB{}, err
	}

	// Disable Gorm's logging
	db.SetLogger(fakeLogger{})

	// Create tables if necessary
	db.AutoMigrate(&Summoner{}, &MarshaledMatchDetail{})

	// Add indices to summoner table
	db.Model(Summoner{}).AddUniqueIndex("idx_name", "name")
	db.Model(Summoner{}).AddIndex("idx_last_crawled", "last_crawled")

	// Add indices to match table
	db.Model(MarshaledMatchDetail{}).AddIndex("idx_match_mode", "match_mode")
	db.Model(MarshaledMatchDetail{}).AddIndex("idx_match_type", "match_mode")
	db.Model(MarshaledMatchDetail{}).AddIndex("idx_match_queue_type", "queue_type")
	db.Model(MarshaledMatchDetail{}).AddIndex("idx_match_season", "season")

	return db, nil
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

// Save a match in the database. Returns any errors that occurred.
func saveMatch(db gorm.DB, match *MatchDetail) error {
	// Convert the MatchDetail so that it can be saved to the database
	res, err := match.Marshal()
	if err != nil {
		return fmt.Errorf("Unable to marshal match details: %d -- %s", match, err.Error())
	}

	// Save the match details
	if err := db.Create(res).Error; err != nil {
		return fmt.Errorf("Unable to save match details: %d -- %s", match, err.Error())
	}

	return nil
}

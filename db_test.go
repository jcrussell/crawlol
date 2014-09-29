package main

import (
	"log"
	"testing"

	"github.com/jcrussell/crawlol/external/github.com/jinzhu/gorm"

	// Import for side-effect
	_ "github.com/jcrussell/crawlol/external/github.com/mattn/go-sqlite3"
)

var db gorm.DB

func init() {
	var err error

	if db, err = openDB(":memory:"); err != nil {
		log.Fatalf("Unable to open database: %s", err.Error())
	}
}

func TestSaveMatch(t *testing.T) {
	details := &MatchDetail{}
	details.Id = 0

	saveMatch(db, details)

	if db.Where(details).First(&MatchDetail{}).RecordNotFound() {
		t.Fatal("Unable to find saved match details")
	}
}

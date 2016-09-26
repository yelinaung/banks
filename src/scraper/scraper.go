package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/yelinaung/banks/pkg/scraper"
	"os"
)

var dbName = "test"
var tableName = "currency"

func main() {
	var s = scraper.NewScraper(dbName, tableName)

	// Do jobs without params
	if os.Getenv("BANKS_MODE") == "release" {
		gocron.Every(2).Hours().Do(scraper.RunScraper, s)
	} else {
		gocron.Every(20).Seconds().Do(scraper.RunScraper, s)
	}

	// Run the job
	<-gocron.Start()
}

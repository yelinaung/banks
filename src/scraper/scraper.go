package main

import (
	"github.com/yelinaung/banks/pkg/scraper"
	"github.com/jasonlvhit/gocron"
)

var dbName = "test"
var tableName = "currency"

func main() {
	var s = scraper.NewScraper(dbName, tableName)

	// Do jobs without params
	gocron.Every(10).Seconds().Do(scraper.RunScraper, s)
	// gocron.Every(1).Day().At("05:30").Do(RunScraper, scraper)

	// Run the job
	<-gocron.Start()
}

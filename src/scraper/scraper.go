package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/yelinaung/banks/pkg/scraper"
)

var dbName = "test"
var tableName = "currency"

func main() {
	var s = scraper.NewScraper(dbName, tableName)

	// Do jobs without params
	// gocron.Every(10).Seconds().Do(scraper.RunScraper, s)
	gocron.Every(1).Day().At("05:30").Do(scraper.RunScraper, s)

	// Run the job
	<-gocron.Start()
}

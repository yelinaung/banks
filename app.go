package main

import (
	"os"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/jasonlvhit/gocron"
)

var session *r.Session
var dbName = "test"
var tableName = "currency"

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
		MaxOpen:  10,
	})

	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	_, err1 := r.DB(dbName).TableCreate(tableName).RunWrite(session)

	if err1 == nil {
		fmt.Printf("Error creating table: %s", err1)
	} else {
		r.DB(dbName).TableCreate(tableName).RunWrite(session)
	}
}

func main() {
	var scraper = NewScraper(dbName, tableName)

	// Do jobs without params
	gocron.Every(10).Seconds().Do(RunScraper, scraper)
	// gocron.Every(1).Day().At("05:30").Do(RunScraper, scraper)

	// Run the job
	<-gocron.Start()

	server := NewAPIServer(os.Getenv("PORT"), tableName, session)
	StartAPIServer(server)
}

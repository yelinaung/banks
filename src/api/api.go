package main

import (
	"fmt"
	"github.com/yelinaung/banks/pkg/api"
	r "github.com/dancannon/gorethink"
	"os"
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
	server := api.NewAPIServer(os.Getenv("PORT"), dbName, tableName, session)
	api.StartAPIServer(server)
}

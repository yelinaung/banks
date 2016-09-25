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
	var connectError error

	session, connectError = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
		MaxOpen:  10,
	})

	if connectError != nil {
		fmt.Errorf("failed to connect to database: %v", connectError)
	}

	_, tableCreateError := r.DB(dbName).TableCreate(tableName).RunWrite(session)

	if tableCreateError == nil {
		fmt.Errorf("Error creating table : %v", tableCreateError)
	} else {
		r.DB(dbName).TableCreate(tableName).RunWrite(session)
	}
}

func main() {
	server := api.NewAPIServer(os.Getenv("PORT"), dbName, tableName, session)
	api.StartAPIServer(server)
}

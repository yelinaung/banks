package main

import (
	"encoding/json"
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/jasonlvhit/gocron"
)

var dbName = "test"
var tableName = "currency"
var s *r.Session

// TODO Add route to get all the currencies of one bank
// TODO Add route to get all the "latest" currencies
// TODO Add route to get all the currencies by "date"
// TODO Add route to get one latest currency of a bank

func init() {
	var err error

	s, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
		MaxOpen:  10,
	})

	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	_, err1 := r.DB(dbName).TableCreate(tableName).RunWrite(s)

	if err1 == nil {
		fmt.Printf("Error creating table: %s", err1)
	} else {
		r.DB(dbName).TableCreate(tableName).RunWrite(s)
	}
}

func main() {
	fmt.Println("Starting..")
	// Do jobs without params
	gocron.Every(2).Minutes().Do(Run)

	<-gocron.Start()

	//ginRoute := gin.New()
	//
	//// Base
	//ginRoute.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK,
	//		"Nothing to see here.Check https://github.com/yelinaung/banks")
	//})
	//
	//ginRoute.GET("/all", func(c *gin.Context) {
	//	currencies, err := getAll()
	//	var response Response
	//	var data Data
	//	data.Currencies = currencies
	//	response.Data = data
	//	if err == nil {
	//		c.JSON(http.StatusOK, response)
	//	} else {
	//		c.JSON(http.StatusInternalServerError,
	//			gin.H{
	//				"message": "Something went wrong!",
	//			})
	//	}
	//
	//})
	//
	//ginRoute.GET("/b/:bank", func(c *gin.Context) {
	//	bankName := c.Params.ByName("bank")
	//	currencies, err := filterByBankName(bankName)
	//	var response Response
	//	var data Data
	//	data.Currencies = currencies
	//	response.Data = data
	//	if err == nil {
	//		c.JSON(http.StatusOK, response)
	//	} else {
	//		c.JSON(http.StatusInternalServerError,
	//			gin.H{
	//				"message": "Something went wrong!",
	//			})
	//	}
	//
	//})
	//ginRoute.Run(":" + os.Getenv("PORT"))
}

func getAll() ([]Currency, error) {
	query := r.Table(tableName)
	row, err := query.Run(s)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var currencies = []Currency{}
	err2 := row.All(&currencies)

	if err2 != nil {
		return nil, err2
	}

	_, err3 := json.Marshal(currencies)

	fmt.Println("currencies ", len(currencies))

	return currencies, err3
}

func filterByBankName(name string) ([]Currency, error) {
	query := r.Table(tableName).Filter(r.Row.Field("bank_name").Eq(name))
	row, err := query.Run(s)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var currencies = []Currency{}
	err2 := row.All(&currencies)

	if err2 != nil {
		return nil, err2
	}

	_, err3 := json.Marshal(currencies)

	fmt.Println("currencies ", len(currencies))

	return currencies, err3
}

package main

import (
	"encoding/json"
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/jasonlvhit/gocron"
	str "strings"
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"os"
)

var dbName = "test"
var tableName = "currency"
var s *r.Session

// DONE Add route to get all the currencies of one bank
// DONE Add route to get all the "latest" currencies
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
	// gocron.Every(1).Day().At("00:30").Do(Run)

	// Run the job
	<-gocron.Start()

	ginRoute := gin.New()

	// Base
	ginRoute.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,
			"Nothing to see here.Check https://github.com/yelinaung/banks")
	})

	ginRoute.GET("/all", func(c *gin.Context) {
		currencies, err := getAll()
		var response Response
		var data Data
		data.Currencies = currencies
		response.Data = data
		if err == nil {
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"message": "Something went wrong!",
				})
		}

	})

	ginRoute.GET("/b/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		currencies, err := getAllCurrenciesByBankName(bankName)
		var response Response
		var data Data
		data.Currencies = currencies
		response.Data = data
		if err == nil {
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"message": "Something went wrong!",
				})
		}

	})

	ginRoute.GET("/latest1", func(c *gin.Context) {
		start := time.Now()
		currencies, err := getAllLatestCurrencies()
		var response Response
		var data Data
		data.Currencies = currencies
		response.Data = data
		if err == nil {
			elapsed := time.Since(start)
			fmt.Printf("latest one took %s\n", elapsed)
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"message": "Something went wrong!",
				})
		}
	})

	ginRoute.Run(":" + os.Getenv("PORT"))
}

func getAll() ([]Currency, error) {
	query := r.Table(tableName)
	return resolveCursorToValue(query)
}

func getAllLatestCurrencies() ([]Currency, error) {
	// a bit hacky way to do it
	query := r.Table(tableName).OrderBy("time").Limit(6)

	// another butt ugly way and super slow way is
	//query := filterLatest("KBZ").
	//	Union(filterLatest("CBB")).
	//	Union(filterLatest("MAB")).
	//	Union(filterLatest("AGD")).
	//	Union(filterLatest("AYA")).
	//	Union(filterLatest("UAB"))

	return resolveCursorToValue(query)
}

func filterLatest(name string) r.Term {
	return r.Table(tableName).OrderBy("time").
		Filter(r.Row.Field("bank_name").Eq(str.ToUpper(name))).
		Limit(1)
}

func getAllCurrenciesByBankName(name string) ([]Currency, error) {
	query := r.Table(tableName).Filter(r.Row.Field("bank_name").Eq(str.ToUpper(name)))
	return resolveCursorToValue(query)
}

func resolveCursorToValue(t r.Term) ([]Currency, error) {
	var currencies = []Currency{}

	row, err := t.Run(s)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	err2 := row.All(&currencies)
	if err2 != nil {
		return nil, err2
	}

	_, err3 := json.Marshal(currencies)

	// fmt.Println("currencies ", len(currencies))
	return currencies, err3
}

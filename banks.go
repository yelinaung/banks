package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	r "github.com/dancannon/gorethink"
	"github.com/gin-gonic/gin"
)

var dbName = "test"
var tableName = "currency"
var s *r.Session

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

	//c := gron.New()
	//c.AddFunc(gron.Every(2 * time.Minute), func() {
	//	fmt.Println("runs every two minute.")

	//})
	//c.Start()

	//
	// r.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK,
	// 		"Nothing to see here.Check https://github.com/yelinaung/banks")
	// })
	//

	r := gin.New()
	r.GET("/run", func(c *gin.Context) {
		Run()
	})

	r.GET("/all", func(c *gin.Context) {
		err, currencies := getAll()
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

	r.GET("/b/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		err, currencies := filterByBankName(bankName)
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
	r.Run(":" + os.Getenv("PORT"))
}

func getAll() (error, [] Currency) {
	query := r.Table(tableName)
	row, err := query.Run(s)
	if err != nil {
		fmt.Print(err)
		return err, nil
	}

	var currencies = []Currency{}
	err2 := row.All(&currencies)

	if err2 != nil {
		return err2, nil
	}

	_, err3 := json.Marshal(currencies)

	fmt.Println("currencies ", len(currencies))

	return err3, currencies
}

func filterByBankName(name string) (error, []Currency) {
	query := r.Table(tableName).Filter(r.Row.Field("bank_name").Eq(name))
	row, err := query.Run(s)
	if err != nil {
		fmt.Print(err)
		return err, nil
	}

	var currencies = []Currency{}
	err2 := row.All(&currencies)

	if err2 != nil {
		return err2, nil
	}

	_, err3 := json.Marshal(currencies)

	fmt.Println("currencies ", len(currencies))

	return err3, currencies
}

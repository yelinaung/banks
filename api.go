package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	str "strings"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/gin-gonic/gin"
)

var api API

func NewAPIServer(port string, tableName string, session *r.Session) API {
	api.port = port
	api.tableName = tableName
	api.session = session
	return api
}

type API struct {
	port      string
	tableName string
	session   *r.Session
}

func StartAPIServer(api API) {
	ginRoute := gin.New()

	// Base
	ginRoute.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,
			"Nothing to see here.Check https://github.com/yelinaung/banks")
	})

	ginRoute.GET("/all", func(c *gin.Context) {
		currencies, err := getAll(api.tableName)
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
		currencies, err := getAllCurrenciesByBankName(api.tableName, bankName)
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

	ginRoute.GET("/latest", func(c *gin.Context) {
		start := time.Now()
		currencies, err := getAllLatestCurrencies(api.tableName)
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

	ginRoute.Run(":" + api.port)
}

func getAll(tableName string) ([]Currency, error) {
	query := r.Table(tableName)
	return resolveCursorToValue(query)
}

func getAllLatestCurrencies(tableName string) ([]Currency, error) {
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

func filterLatest(tableName string, name string) r.Term {
	return r.Table(tableName).OrderBy("time").
		Filter(r.Row.Field("bank_name").Eq(str.ToUpper(name))).
		Limit(1)
}

func getAllCurrenciesByBankName(tableName string, name string) ([]Currency, error) {
	query := r.Table(tableName).Filter(r.Row.Field("bank_name").Eq(str.ToUpper(name)))
	return resolveCursorToValue(query)
}

func resolveCursorToValue(t r.Term) ([]Currency, error) {
	var currencies = []Currency{}

	row, err := t.Run(api.session)
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

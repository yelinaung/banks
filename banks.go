package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	str "strings"
	"time"

	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

var (
	kbz = "http://www.kbzbank.com"
	cb  = "http://www.cbbank.com.mm/exchange_rate.aspx"
	aya = "http://ayabank.com"
	mab = "http://www.mabbank.com"
	uab = "http://www.unitedamarabank.com"

	// Turns out AGD was loading data through ajax
	agd = "https://ibanking.agdbank.com.mm/RateInfo?id=ALFKI&callback=?"
)

func scrapKBZ() []string {
	tmp := []string{}

	// Using with file
	// f, err := os.Open("agd.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocument(agd)

	doc, err := goquery.NewDocument(kbz)
	PanicIf(err)

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp
}

func scrapAGD() []string {
	tmp := []string{}

	response, err := http.Get(agd)
	PanicIf(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	PanicIf(err)

	// contents has extra characters which causes
	// invalid json structure
	st := string(contents)
	st = str.Replace(st, "?", "", -1)
	st = str.Replace(st, "(", "", -1)
	st = str.Replace(st, ")", "", -1)
	st = str.Replace(st, ";", "", -1)

	a := new(AGD)
	json.Unmarshal([]byte(st), a)

	tmp = append(tmp, "EURO")
	tmp = append(tmp, floatToString(a.Exchangerates[1].Rate))
	tmp = append(tmp, floatToString(a.Exchangerates[0].Rate))
	tmp = append(tmp, "SGD")
	tmp = append(tmp, floatToString(a.Exchangerates[3].Rate))
	tmp = append(tmp, floatToString(a.Exchangerates[2].Rate))
	tmp = append(tmp, "USD")
	tmp = append(tmp, floatToString(a.Exchangerates[5].Rate))
	tmp = append(tmp, floatToString(a.Exchangerates[4].Rate))

	return tmp
}

func scrapCB() []string {
	tmp := []string{}

	doc, err := goquery.NewDocument(cb)
	PanicIf(err)

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp
}

func scrapAYA() []string {
	tmp := []string{}

	doc, err := goquery.NewDocument(aya)
	PanicIf(err)

	doc.Find("#tablepress-2 tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp
}

func scrapMAB() []string {
	tmp := []string{}

	// Using with file
	f, err := os.Open("mab.html")
	PanicIf(err)
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)

	PanicIf(err)

	doc.Find("#block-block-5 tbody tr").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp
}

func process(tmp []string) Bank {
	bank := Bank{}

	bank.Base = "MMK"
	bank.Time = time.Now().String()

	currencies := []string{tmp[0], tmp[3], tmp[6]}
	buy := []string{tmp[1], tmp[4], tmp[7]}
	sell := []string{tmp[2], tmp[5], tmp[8]}

	for x, _ := range currencies {
		bank.Rates = append(bank.Rates, map[string]BuySell{
			currencies[x]: BuySell{buy[x], sell[x]}})
	}

	return bank
}

func main() {

	fmt.Println(scrapMAB())

	r := gin.Default()
	//
	var bank Bank
	r.GET("/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		if bankName == "kbz" {
			bank = process(scrapKBZ())
			bank.Name = "KBZ"
		} else if bankName == "mab" {
			bank = process(scrapMAB())
			bank.Name = "AYA"
		}
		//	} else if bankName == "cb" {
		//		bank = process(scrapCB())
		//		bank.Name = "CB"
		//	} else if bankName == "agd" {
		//		bank = process(scrapAGD())
		//		bank.Name = "AGD"
		//	} else if bankName == "aya" {
		//		bank = process(scrapAYA())
		//		bank.Name = "AYA"
		//	}
		c.JSON(200, bank)
	})
	r.Run(":" + os.Getenv("PORT"))
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Bank struct {
	Name  string               `json:"name"`
	Base  string               `json:"base"`
	Time  string               `json:"time"`
	Rates []map[string]BuySell `json:"rates"`
}

type AGD struct {
	Exchangerates []struct {
		From string  `json:"From"`
		To   string  `json:"To"`
		Rate float64 `json:"Rate"`
	} `json:"ExchangeRates"`
}

type BuySell struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

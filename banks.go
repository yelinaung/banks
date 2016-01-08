package main

import (
	"fmt"
	"os"
	str "strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

var (
	kbz = "http://www.kbzbank.com"
	cb  = "http://www.cbbank.com.mm/exchange_rate.aspx"
	agd = "http://myanmarbank.asia/category/exchange-rate"
	aya = "http://ayabank.com"
)

func scrapKBZ() []string {
	tmp := []string{}

	doc, err := goquery.NewDocument(kbz)
	PanicIf(err)

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp
}

// Scraping from the url NOTworking (yet)
// Found the cause but no solution yet
func scrapAGD() []string {
	tmp := []string{}

	// Using with file
	// f, err := os.Open("agd.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocumentFromReader(f)
	doc, err := goquery.NewDocument(agd)
	PanicIf(err)

	// doc.Find("#curency-table").Slice(1, 4).Find("tr").Each(func(i int, s *goquery.Selection) {
	fmt.Println(str.TrimSpace(doc.Find("#curency-table tbody").Text()))
	// tmp = append(tmp, str.TrimSpace(s.Text()))
	//s.Find("td").Each(func(u int, t *goquery.Selection) {
	//	tmp = append(tmp, str.TrimSpace(t.Text()))
	//})
	// })

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

	// fmt.Println("agd ", scrapAGD())

	r := gin.Default()

	var bank Bank
	r.GET("/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		if bankName == "kbz" {
			bank = process(scrapKBZ())
			bank.Name = "KBZ"
		} else if bankName == "cb" {
			bank.Name = "CB"
			bank = process(scrapCB())
			//		} else if bankName == "agd" {
			//			bank.Name = "AGD"
			//			bank = process(scrapAGD())
		} else if bankName == "aya" {
			bank = process(scrapAYA())
			bank.Name = "AYA"
		}
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

type BuySell struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

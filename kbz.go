package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	str "strings"
	"time"
)

var (
	kbz = "http://www.kbzbank.com"
	cb  = "http://www.cbbank.com.mm/exchange_rate.aspx"
)

func ScrapKBZ(url string) ([]string, string) {
	tmp := []string{}

	// Using with file
	// f, err := os.Open("kbz.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocumentFromReader(f)

	doc, err := goquery.NewDocument(url)
	PanicIf(err)

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, t.Text())
		})
	})

	return tmp, "kbz"
}

func ScrapCB(url string) ([]string, string) {
	tmp := []string{}

	doc, err := goquery.NewDocument(url)
	PanicIf(err)

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp, "cb"
}

func Process(temp []string, bName string) Bank {
	currencies := []string{}
	buy := []string{}
	sell := []string{}

	k := Bank{}

	if bName == "kbz" {
		k.Name = "KBZ"
		k.Base = "MMK"
	} else if bName == "cb" {
		k.Name = "CB"
		k.Base = "MMK"
	}

	k.Time = time.Now().String()

	// I don't know why I do this lol
	for j, _ := range temp {
		if j%3 == 0 {
			currencies = append(currencies, str.TrimSpace(temp[j]))
		}
	}

	buy = append(buy, temp[1], temp[4], temp[7])
	sell = append(sell, temp[2], temp[5], temp[8])

	for x, _ := range currencies {
		k.Rates = append(k.Rates, map[string]BuySell{
			currencies[x]: BuySell{buy[x], sell[x]}})
	}

	return k
}

func main() {

	rawKBZ, x := ScrapKBZ(kbz)
	rawCB, y := ScrapCB(cb)

	router := gin.Default()

	router.GET("/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")

		if bankName == "kbz" {
			bank := Process(rawKBZ, x)
			c.JSON(200, bank)
		} else if bankName == "cb" {
			bank := Process(rawCB, y)
			c.JSON(200, bank)
		}
	})

	router.Run(":3001")
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

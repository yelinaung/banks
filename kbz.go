package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	str "strings"
	"time"
)

var (
	kbz = "http://www.kbzbank.com"
)

func ScrapWork(url string) []string {
	temp := []string{}

	// Using with file
	// f, err := os.Open("kbz.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocumentFromReader(f)

	doc, err := goquery.NewDocument(url)
	PanicIf(err)

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			temp = append(temp, t.Text())
		})
	})

	return temp
}

func Process(temp []string) Bank {
	currencies := []string{}
	buy := []string{}
	sell := []string{}

	k := Bank{}

	k.Name = "KBZ"
	k.Base = "MMK"
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

func main2() {

	raw := ScrapWork(kbz)
	bank := Process(raw)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, bank)
	})
	router.Run(":3001")
}

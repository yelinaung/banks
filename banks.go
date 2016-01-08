package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	str "strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	kbzURL = "http://www.kbzbank.com"
	cbbURL = "http://www.cbbank.com.mm/exchange_rate.aspx"
	ayaURL = "http://ayabank.com"
	mabURL = "http://www.mabbank.com"
	uabURL = "http://www.unitedamarabank.com"
	agdURL = "https://ibanking.agdbank.com.mm/RateInfo?id=ALFKI&callback=?"
)

func scrapKBZ() []string {
	tmp := []string{}

	// Using with file
	// f, err := os.Open("agd.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocument(agd)

	doc, err := goquery.NewDocument(kbzURL)
	panicIf(err)

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp
}

func scrapUAB() []string {
	tmp := []string{}
	//f, err := os.Open("")
	//PanicIf(err)
	//defer f.Close();
	//doc, err := goquery.NewDocumentFromReader(f)
	doc, err := goquery.NewDocument(uabURL)
	panicIf(err)

	doc.Find(".ex_rate .ex_body").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("ul li").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})
	return tmp
}

func scrapAGD() []string {
	tmp := []string{}

	response, err := http.Get(agdURL)
	panicIf(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	panicIf(err)

	// contents has extra characters which causes
	// invalid json structure
	st := string(contents)
	st = str.Replace(st, "?", "", -1)
	st = str.Replace(st, "(", "", -1)
	st = str.Replace(st, ")", "", -1)
	st = str.Replace(st, ";", "", -1)

	a := new(agd)
	json.Unmarshal([]byte(st), a)

	tmp = append(tmp, "EURO")
	tmp = append(tmp, floatToString(a.ExchangeRates[1].Rate))
	tmp = append(tmp, floatToString(a.ExchangeRates[0].Rate))
	tmp = append(tmp, "SGD")
	tmp = append(tmp, floatToString(a.ExchangeRates[3].Rate))
	tmp = append(tmp, floatToString(a.ExchangeRates[2].Rate))
	tmp = append(tmp, "USD")
	tmp = append(tmp, floatToString(a.ExchangeRates[5].Rate))
	tmp = append(tmp, floatToString(a.ExchangeRates[4].Rate))

	return tmp
}

func scrapCBB() []string {
	tmp := []string{}

	doc, err := goquery.NewDocument(cbbURL)
	panicIf(err)

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp
}

func scrapAYA() []string {
	tmp := []string{}

	doc, err := goquery.NewDocument(ayaURL)
	panicIf(err)

	doc.Find("#tablepress-2 tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp
}

func scrapMAB() []string {
	tmp := []string{}
	doc, err := goquery.NewDocument(mabURL)
	panicIf(err)

	doc.Find("#block-block-5 tbody tr").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp
}

func process(tmp []string) bank {
	bank := bank{}

	bank.Base = "MMK"
	bank.Time = time.Now().String()

	currencies := []string{tmp[0], tmp[3], tmp[6]}
	buy := []string{tmp[1], tmp[4], tmp[7]}
	sell := []string{tmp[2], tmp[5], tmp[8]}

	for x := range currencies {
		bank.Rates = append(bank.Rates, map[string]buySell{
			currencies[x]: buySell{buy[x], sell[x]}})
	}

	return bank
}

func main() {

	//fmt.Println("UAB ", scrapUAB())

	r := gin.Default()
	//
	var bank bank
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,
			"Nothing to see here.Check https://github.com/yelinaung/banks")
	})
	r.GET("/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		switch bankName {
		case "kbz":
			bank = process(scrapKBZ())
			bank.Name = "KBZ"
		case "mab":
			bank = process(scrapMAB())
			bank.Name = "MAB"
		case "uab":
			bank = process(scrapUAB())
			bank.Name = "UAB"
		case "cbb":
			bank = process(scrapCBB())
			bank.Name = "CBB"
		case "agd":
			bank = process(scrapAGD())
			bank.Name = "AGD"
		case "aya":
			bank = process(scrapAYA())
			bank.Name = "AYA"
		default:
		// TODO	what to reply for default
		}
		c.JSON(http.StatusOK, bank)
	})
	r.Run(":" + os.Getenv("PORT"))
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type bank struct {
	Name  string               `json:"name"`
	Base  string               `json:"base"`
	Time  string               `json:"time"`
	Rates []map[string]buySell `json:"rates"`
}

type agd struct {
	ExchangeRates []struct {
		From string  `json:"From"`
		To   string  `json:"To"`
		Rate float64 `json:"Rate"`
	} `json:"ExchangeRates"`
}

type buySell struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

func floatToString(inputName float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputName, 'f', 2, 64)
}

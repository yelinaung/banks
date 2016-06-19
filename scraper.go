package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	str "strings"
	"time"
	"fmt"

	r "github.com/dancannon/gorethink"
	"log"
)

var s *r.Session

var DB_NAME = "currency_test"
var TABLE_NAME = "currency"

// Bank Urls
var (
	kbzURL = "https://www.kbzbank.com"
	cbbURL = "http://www.cbbank.com.mm/exchange_rate.aspx"
	ayaURL = "http://ayabank.com"
	mabURL = "http://www.mabbank.com"
	uabURL = "http://www.unitedamarabank.com"
	agdURL = "https://ibanking.agdbank.com.mm/RateInfo?id=ALFKI&callback=?"
)

func initDb() {
	var session *r.Session

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	session.SetMaxOpenConns(5)

	r.TableDrop(TABLE_NAME).Run(session)

	resp, err := r.TableCreate(TABLE_NAME).RunWrite(session)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%d table created\n", resp.TablesCreated)

	s = session
}

func Run() {

	initDb()

	fmt.Println("Running...")
	kbzData := new(currency)

	// Need to remove file before extracting dat
	os.Remove("kbz")
	kbzData.Time = time.Now().String()
	kbzData.BankName = "kbz"
	kbzData.Bank, _ = process(scrapKBZ())
	response, err := r.Table(TABLE_NAME).Insert(kbzData).RunWrite(s)
	panicIf(err)

	fmt.Printf("%d row inserted\n", response.Inserted)

	//
	//var uabData currency
	//
	//uabData.Time = time.Now().String()
	//uabData.BankName = "uab"
	//uabData.Bank, _ = process(scrapUAB())
	//
	//uabResult, err := r.Table(TABLE_NAME).Insert(uabData).RunWrite(s)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Printf("UAB result %s\n", uabResult.GeneratedKeys[0])
}

func scrapKBZ() ([]string, error) {
	tmp := []string{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	fileName := "kbz"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		output, _ := os.Create(fileName)
		defer output.Close()
		response, _ := client.Get(kbzURL)
		defer response.Body.Close()

		io.Copy(output, response.Body)
	}

	f, err := os.Open(fileName)
	panicIf(err)

	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	panicIf(err)

	if err == nil {
		doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
			s.Find("td").Each(func(u int, t *goquery.Selection) {
				tmp = append(tmp, str.TrimSpace(t.Text()))
			})
		})
	}

	return tmp, err
}

func scrapUAB() ([]string, error) {
	tmp := []string{}
	doc, err := goquery.NewDocument(uabURL)

	doc.Find(".ex_rate .ex_body").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("ul li").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})
	return tmp, err
}

func scrapAGD() ([]string, error) {
	tmp := []string{}

	response, err1 := http.Get(agdURL)
	//panicIf(err)
	defer response.Body.Close()

	contents, err2 := ioutil.ReadAll(response.Body)
	panicIf(err2)

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

	return tmp, err1
}

func scrapCBB() ([]string, error) {
	tmp := []string{}

	doc, err := goquery.NewDocument(cbbURL)
	panicIf(err)

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp, err
}

func scrapAYA() ([]string, error) {
	tmp := []string{}

	doc, err := goquery.NewDocument(ayaURL)
	//panicIf(err)

	doc.Find("#tablepress-2 tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return tmp, err
}

func scrapMAB() ([]string, error) {
	tmp := []string{}
	doc, err := goquery.NewDocument(mabURL)
	//panicIf(err)

	doc.Find("#block-block-5 tbody tr").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return tmp, err
}

func process(tmp []string, err error) (*bank, error) {
	bank := new(bank)

	bank.Base = "MMK"
	bank.Time = time.Now().String()

	if len(tmp) > 0 {
		currencies := []string{tmp[0], tmp[3], tmp[6]}
		buy := []string{tmp[1], tmp[4], tmp[7]}
		sell := []string{tmp[2], tmp[5], tmp[8]}

		for x := range currencies {
			bank.Rates = append(bank.Rates, map[string]buySell{
				currencies[x]: buySell{buy[x], sell[x]}})
		}
	}

	return bank, err
}

func floatToString(inputName float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputName, 'f', 2, 64)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type currency struct {
	Id       string `gorethink:"id"`
	Time     string `gorethink:"time"`
	BankName string `gorethink:"name"`
	Bank     *bank  `gorethink:"data"`
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

func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}

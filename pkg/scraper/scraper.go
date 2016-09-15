package scraper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	str "strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	r "github.com/dancannon/gorethink"
)

// Bank Urls
var (
	kbzURL = "https://www.kbzbank.com"
	cbbURL = "http://www.cbbank.com.mm/exchange_rate.aspx"
	ayaURL = "http://ayabank.com"
	mabURL = "http://www.mabbank.com"
	uabURL = "http://www.unitedamarabank.com"
	agdURL = "https://ibanking.agdbank.com.mm/RateInfo?id=ALFKI&callback=?"

	KBZ = "KBZ"
	CBB = "CBB"
	AYA = "AYA"
	MAB = "MAB"
	UAB = "UAB"
	AGD = "AGD"

	scraper Scraper
	session *r.Session
)

func NewScraper(dbName string, tableName string) Scraper {
	scraper.dbName = dbName
	scraper.tableName = tableName
	scraper.Init()
	return scraper
}

type Scraper struct {
	dbName    string
	tableName string
}

func (scraper Scraper) Init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: scraper.dbName,
		MaxOpen:  10,
	})

	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	_, err1 := r.DB(scraper.dbName).TableCreate(scraper.tableName).RunWrite(session)

	if err1 == nil {
		fmt.Printf("Error creating table: %s", err1)
	} else {
		r.DB(scraper.dbName).TableCreate(scraper.tableName).RunWrite(session)
	}
}

func generateID() (string, error) {
	// Create a new unique ID.
	rethink, err := r.UUID().Run(session)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	// Get the value.
	var id string
	err = rethink.One(&id)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	if len(id) == 0 {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	return id, nil
}

// Run is the main function to kick of all the scraping work
// it scraps, process the data and put into the db
func RunScraper(scraper Scraper) {
	fmt.Println("Running...")

	// Need to remove file before extracting dat
	os.Remove("kbz")
	kbzData, err := process(scrapKBZ())
	panicIf(err)
	kbzResult, err := writeToDb(scraper, kbzData)
	printLog(err, KBZ, kbzResult)

	// UAB
	uabData, err := process(scrapUAB())
	panicIf(err)
	uabResult, err := writeToDb(scraper, uabData)
	printLog(err, UAB, uabResult)

	// AGD Bank
	agdData, err := process(scrapAGD())
	panicIf(err)
	agdResult, err := writeToDb(scraper, agdData)
	printLog(err, AGD, agdResult)

	// CBB
	cbbData, err := process(scrapCBB())
	panicIf(err)
	cbbResult, err := writeToDb(scraper, cbbData)
	printLog(err, CBB, cbbResult)

	// AYA
	ayaData, err := process(scrapAYA())
	panicIf(err)
	ayaResult, err := writeToDb(scraper, ayaData)
	printLog(err, AYA, ayaResult)

	// MAB
	mabData, err := process(scrapMAB())
	panicIf(err)
	mabResult, err := writeToDb(scraper, mabData)
	printLog(err, MAB, mabResult)
}

func printLog(err error, bankName string, response r.WriteResponse) {
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d inserted for %s bank \n", response.Inserted, bankName)
}

func writeToDb(scraper Scraper, data Currency) (r.WriteResponse, error) {
	return r.Table(scraper.tableName).Insert(data).RunWrite(session)
}

func scrapKBZ() (string, []string, error) {
	tmp := [][]string{}

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
		doc.Find(".exchange-rate div").Each(func(i int, s *goquery.Selection) {
			s.Find(".col-lg-2").Each(func(j int, t *goquery.Selection) {
				x := str.TrimSpace(t.After("strong").Text())
				tmp = append(tmp, str.Split(x, " "))
			})
		})
	}

	return "KBZ", flatten(tmp), err
}

func scrapUAB() (string, []string, error) {
	tmp := []string{}
	doc, err := goquery.NewDocument(uabURL)

	doc.Find(".ex_rate .ex_body").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("ul li").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})
	return "UAB", tmp, err
}

func scrapAGD() (string, []string, error) {
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
	// TODO They added Thai bhatt .. hmm mm
	// tmp = append(tmp, floatToString(a.ExchangeRates[5].Rate))
	// tmp = append(tmp, floatToString(a.ExchangeRates[4].Rate))
	tmp = append(tmp, floatToString(a.ExchangeRates[7].Rate))
	tmp = append(tmp, floatToString(a.ExchangeRates[6].Rate))

	return "AGD", tmp, err1
}

func scrapCBB() (string, []string, error) {
	tmp := []string{}

	doc, err := goquery.NewDocument(cbbURL)
	panicIf(err)

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	return "CBB", tmp, err
}

func scrapAYA() (string, []string, error) {
	tmp := []string{}

	doc, err := goquery.NewDocument(ayaURL)

	regex, _ := regexp.Compile(`\D`)
	doc.Find("#tablepress-2 tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		sText := s.Text()
		if str.Contains(sText, "/") {
			tmp = append(tmp, str.TrimSpace(regex.ReplaceAllLiteralString(sText, "")))
		}
		tmp = append(tmp, str.TrimSpace(sText))
	})

	return "AYA", tmp, err
}

func scrapMAB() (string, []string, error) {
	tmp := []string{}
	doc, err := goquery.NewDocument(mabURL)

	doc.Find("#block-block-5 tbody tr").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return "MAB", tmp, err
}

func process(bankName string, tmp []string, err error) (Currency, error) {
	bank := Currency{}
	id, err := generateID()
	bank.ID = id
	bank.Name = bankName
	bank.Base = "MMK"
	bank.Time = time.Now().String()

	if len(tmp) > 0 {
		currencies := []string{tmp[0], tmp[3], tmp[6]}
		buy := []string{tmp[1], tmp[4], tmp[7]}
		sell := []string{tmp[2], tmp[5], tmp[8]}

		for x := range currencies {
			bank.Rates = append(bank.Rates, map[string]buySell{
				currencies[x]: {buy[x], sell[x]}})
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

func flatten(input [][]string) []string {
	tmp := []string{}
	for i := 0; i < len(input); i++ {
		x := input[i]
		for j := 0; j < len(x); j++ {
			if len(x[j]) != 0 {
				tmp = append(tmp, x[j])
			}
		}
	}

	return tmp
}

// Response type for the api
// Doesn't include in the DB
type Response struct {
	Data Data `json:"data"`
}

// Data is the type which serves as "data" for the Response
// Doesn't include in the DB
type Data struct {
	Currencies []Currency `json:"currencies"`
}

// Currency is the type with all the necessary data
// Data is stored in the DB
type Currency struct {
	ID    string               `json:"id" gorethink:"id"`
	Name  string               `json:"bank_name" gorethink:"bank_name"`
	Base  string               `json:"base" gorethink:"base"`
	Time  string               `json:"time" gorethink:"time"`
	Rates []map[string]buySell `json:"rates" gorethink:"rates"`
}

type agd struct {
	ExchangeRates []struct {
		From string  `json:"From"`
		To   string  `json:"To"`
		Rate float64 `json:"Rate"`
	} `json:"ExchangeRates"`
}

type buySell struct {
	Buy  string `json:"buy" gorethink:"buy"`
	Sell string `json:"sell" gorethink:"sell"`
}

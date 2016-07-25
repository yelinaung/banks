package main

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

var dbName = "test"

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
)

func generateID() (string, error) {
	// Create a new unique ID.
	r, err := r.UUID().Run(s)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	// Get the value.
	var id string
	err = r.One(&id)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	if len(id) == 0 {
		return "", fmt.Errorf("failed to obtain a new unique ID: %v", err)
	}

	return id, nil
}

func Run() {

	fmt.Println("Running...")

	// KBZ
	kbzData := new(Currency)

	// Need to remove file before extracting dat
	os.Remove("kbz")
	id, err := generateID()
	panicIf(err)
	kbzData.ID = id
	kbzData.Time = time.Now().String()
	kbzData.BankName = KBZ
	kbzData.Bank, _ = process(scrapKBZ())
	kbzResult, err := r.Table(tableName).Insert(kbzData).RunWrite(s)
	printLog(err, KBZ, kbzResult)

	// UAB
	var uabData Currency
	uabID, err := generateID()
	panicIf(err)
	uabData.ID = uabID
	uabData.Time = time.Now().String()
	uabData.BankName = UAB
	uabData.Bank, _ = process(scrapUAB())
	uabResult, err := writeToDb(uabData)
	printLog(err, UAB, uabResult)

	// AGD Bank
	var agdData Currency
	agdID, err := generateID()
	panicIf(err)

	agdData.ID = agdID
	agdData.Time = time.Now().String()
	agdData.BankName = AGD
	agdData.Bank, _ = process(scrapAGD())
	agdResult, err := writeToDb(agdData)
	printLog(err, AGD, agdResult)

	// CBB
	var cbbData Currency
	cbbID, err := generateID()
	cbbData.ID = cbbID
	cbbData.Time = time.Now().String()
	cbbData.BankName = CBB
	cbbData.Bank, _ = process(scrapCBB())

	cbbResult, err := writeToDb(cbbData)
	printLog(err, CBB, cbbResult)

	// AYA
	var ayaData Currency
	ayaID, err := generateID()
	ayaData.ID = ayaID
	ayaData.Time = time.Now().String()
	ayaData.BankName = AYA
	ayaData.Bank, _ = process(scrapAYA())
	ayaResult, err := writeToDb(ayaData)
	printLog(err, AYA, ayaResult)

	// MAB
	var mabData Currency
	mabID, err := generateID()
	mabData.ID = mabID
	mabData.Time = time.Now().String()
	mabData.BankName = MAB
	mabData.Bank, _ = process(scrapMAB())

	mabResult, err := writeToDb(mabData)
	printLog(err, MAB, mabResult)
}

func printLog(err error, bankName string, response r.WriteResponse) {
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d inserted for %s bank \n", response.Inserted, bankName)
}

func writeToDb(data Currency) (r.WriteResponse, error) {
	return r.Table(tableName).Insert(data).RunWrite(s)
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

	return "KBZ", flattern(tmp), err
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
	// They added Thai bhatt .. hmm mm
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

	r, _ := regexp.Compile(`\D`)
	doc.Find("#tablepress-2 tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		sText := s.Text()
		if str.Contains(sText, "/") {
			tmp = append(tmp, str.TrimSpace(r.ReplaceAllLiteralString(sText, "")))
		}
		tmp = append(tmp, str.TrimSpace(sText))
	})

	return "AYA", tmp, err
}

func scrapMAB() (string, []string, error) {
	tmp := []string{}
	doc, err := goquery.NewDocument(mabURL)
	//panicIf(err)

	doc.Find("#block-block-5 tbody tr").Slice(1, 4).Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			tmp = append(tmp, str.TrimSpace(t.Text()))
		})
	})

	return "MAB", tmp, err
}

func process(bankName string, tmp []string, err error) (*bank, error) {
	bank := new(bank)

	bank.Name = bankName
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

func flattern(input [][]string) []string {
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

type Currency struct {
	ID       string `gorethink:"id"`
	Time     string `gorethink:"time"`
	BankName string `gorethink:"name"`
	Bank     *bank  `gorethink:"data"`
}

type bank struct {
	Name  string               `json:"name" gorethink:"name"`
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

func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}

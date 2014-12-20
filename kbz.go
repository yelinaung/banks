package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	str "strings"
)

var (
	kbz  = "http://www.kbzbank.com"
	temp = []string{}
)

func main() {

	f, err := os.Open("kbz.html")
	PanicIf(err)

	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		PanicIf(err)
	}

	doc.Find(".answer tbody tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(u int, t *goquery.Selection) {
			temp = append(temp, t.Text())
		})
	})

	for j := 0; j < len(temp); j++ {
		// fmt.Println(j)
		if j%3 == 0 {
			fmt.Println(str.TrimSpace(temp[j]))
		}
	}
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Bank struct {
	Name  string
	Time  string
	Base  string
	Rates []string
}

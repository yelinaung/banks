package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
)

var (
	kbz = "http://www.kbzbank.com"
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
		txt := s.Text()
		fmt.Println("------------------")
		fmt.Println(txt)
		fmt.Println("------------------")
	})
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

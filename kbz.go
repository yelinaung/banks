package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

var (
	kbz = "http://www.kbzbank.com"
)

func main() {
	doc, err := goquery.NewDocument(kbz)
	if err != nil {
		PanicIf(err)
	}

	doc.Find(".answer").Each(func(i int, s *goquery.Selection) {
		txt := s.Find("p").Text()
		fmt.Println(txt)
	})
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

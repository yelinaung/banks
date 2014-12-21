package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	str "strings"
	"time"
)

func main2() {
	f, err := os.Open("cb.html")
	defer f.Close()
	PanicIf(err)

	doc, err := goquery.NewDocumentFromReader(f)
	PanicIf(err)

	tmp := []string{}

	doc.Find("table tr").Slice(1, 4).Find("td").Each(func(i int, s *goquery.Selection) {
		tmp = append(tmp, str.TrimSpace(s.Text()))
	})

	currencies := []string{tmp[0], tmp[3], tmp[6]}
	buy := []string{tmp[1], tmp[4], tmp[7]}
	sell := []string{tmp[2], tmp[5], tmp[8]}

	k := Bank{}

	k.Name = "CB"
	k.Base = "MMK"

	// Probably use the time from cb bank ?
	k.Time = time.Now().String()

	fmt.Println(k)
}

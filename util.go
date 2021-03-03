package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getListCount(doc *goquery.Document) int {
	pageText := doc.Find(".pageNav").First().Text()
	if len(strings.Trim(pageText, "")) == 0 {
		return 0
	}
	pageText = pageText[strings.Index(pageText, "Of ")+len("Of "):]
	pageText = pageText[:strings.Index(pageText, " ")]
	pageCount, err := strconv.Atoi(pageText)
	if err != nil {
		log.Println(err)
		return 0
	}
	return pageCount
}

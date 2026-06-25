package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Fatal(err)
	}

	heading := doc.Find("h1, h2").First().Text()

	return strings.TrimSpace(heading)
}

func getFirstParagraphFromHTML(html string) string {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return ""
	}

	mainSelection := doc.Find("main")
	paragraphText := ""

	if mainSelection.Length() > 0 {
		paragraphText = mainSelection.Find("p").First().Text()
	} else {
		paragraphText = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(paragraphText)

}

package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(hmtl, pageURL string) PageData {
	heading := getHeadingFromHTML(hmtl)
	firstParagraph := getFirstParagraphFromHTML(hmtl)

	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		fmt.Printf("could not parse page url %v: %v\n", pageURL, err)
		return PageData{
			URL:            pageURL,
			Heading:        heading,
			FirstParagraph: firstParagraph,
			OutgoingLinks:  nil,
			ImageURLs:      nil,
		}
	}

	links, err := getURLsFromHTML(hmtl, parsedURL)
	if err != nil {
		fmt.Printf("could not get urls from html: %v\n", err)
		return PageData{}
	}
	images, err := getImagesFromHTML(hmtl, parsedURL)
	if err != nil {
		fmt.Printf("could not get images from html: %v\n", err)
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  links,
		ImageURLs:      images,
	}
}

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

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	var urls []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || strings.TrimSpace(href) == "" {
			return
		}

		url, err := url.Parse(href)
		if err != nil {
			fmt.Printf("could not parse href %v: %v", href, err)
			return
		}

		urls = append(urls, baseURL.ResolveReference(url).String())
	})

	return urls, err
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	var imageURLs []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		imgSrc, ok := s.Attr("src")
		if !ok || strings.TrimSpace(imgSrc) == "" {
			return
		}

		url, err := url.Parse(imgSrc)
		if err != nil {
			fmt.Printf("could not parse img src %v: %v", imgSrc, err)
			return
		}

		imageURLs = append(imageURLs, baseURL.ResolveReference(url).String())
	})

	return imageURLs, err
}

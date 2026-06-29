package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getHTML(rawURL string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")
	res, err := client.Do(req)
	if res.StatusCode > 399 {
		return "", fmt.Errorf("%v error: %v", res.StatusCode, err)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("content type not text/html: %v, actually %v", err, res.Header.Get("Content-Type"))
	}
	if err != nil {
		return "", fmt.Errorf("error getting response: %v", err)
	}

	defer res.Body.Close()

	htmlBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading request body: %v", err)
	}

	return string(htmlBody), nil
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

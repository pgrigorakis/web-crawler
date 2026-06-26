package main

import (
	"log"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Printf("could not parse url %v: %v\n", rawBaseURL, err)
		return
	}

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("could not parse url %v: %v\n", rawCurrentURL, err)
		return
	}

	if parsedBaseURL.Hostname() != parsedCurrentURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("could not normalize url %v: %v\n", rawCurrentURL, err)
		return
	}

	if _, ok := pages[normalizedCurrentURL]; ok {
		pages[normalizedCurrentURL]++
		return
	} else {
		pages[normalizedCurrentURL] = 1

		htmlBody, err := getHTML(rawCurrentURL)
		if err != nil {
			log.Printf("could not get HTML content from %v: %v\n", normalizedCurrentURL, err)
			return
		}

		urls, err := getURLsFromHTML(htmlBody, parsedCurrentURL)
		if err != nil {
			log.Printf("could not get urls from %v: %v", normalizedCurrentURL, err)
			return
		}

		for _, url := range urls {
			crawlPage(rawBaseURL, url, pages)
		}
	}
}

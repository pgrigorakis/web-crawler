package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	maxPages           int
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		cfg.wg.Done()
		<-cfg.concurrencyControl
	}()

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("could not parse url %v: %v\n", rawCurrentURL, err)
		return
	}

	if cfg.baseURL.Hostname() != parsedCurrentURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("could not normalize url %v: %v\n", rawCurrentURL, err)
		return
	}

	if !cfg.addPageVisit(normalizedCurrentURL) {
		return
	}

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("could not get HTML content from %v: %v\n", normalizedCurrentURL, err)
		return
	}

	pageData := extractPageData(htmlBody, rawCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedCurrentURL] = pageData
	cfg.mu.Unlock()

	fmt.Printf("crawling from %v\n", normalizedCurrentURL)

	for _, nextURL := range cfg.pages[normalizedCurrentURL].OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}

}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}
	cfg.pages[normalizedURL] = PageData{URL: normalizedURL}
	return true
}

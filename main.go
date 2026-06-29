package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {

	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("usage: crawler BASE_URL MAX_CONCURRENCY MAX_PAGES")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	userURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("MAX_CONCURRENCY argument was not an integer, please provide an integer")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("MAX_PAGES argument was not an integer, please provide an integer")
		os.Exit(1)
	}

	fmt.Printf("URL: %v\nMAX_PAGES: %d\nMAX_CONCURRENCY: %d\n", userURL, maxPages, maxConcurrency)

	fmt.Printf("starting crawl of: %s...\n", userURL)
	baseURL, err := url.Parse(userURL)
	if err != nil {
		fmt.Printf("could not parse url %v: %v\n", userURL, err)
		os.Exit(1)
	}

	cfg := config{
		pages:              map[string]PageData{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		maxPages:           maxPages,
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(userURL)
	cfg.wg.Wait()

	for url, _ := range cfg.pages {
		fmt.Printf("URL: %v\n", url)
	}

	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Printf("could not create JSON report: %v", err)
	}
}

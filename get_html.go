package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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

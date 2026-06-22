package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {
	urlParsed, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("could not parse url: %v", err)
	}

	url := strings.ToLower(fmt.Sprintf("%v%v", urlParsed.Host, strings.TrimSuffix(urlParsed.Path, "/")))

	return url, nil
}

package main

import (
	"fmt"
	"net/url"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
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

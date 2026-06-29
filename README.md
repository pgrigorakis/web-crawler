# web-crawler

A concurrent web crawler written in Go. It crawls all internal pages of a
single website, extracts basic page data (heading, first paragraph, outgoing
links, image URLs), and writes the result to a JSON report.

## Requirements

- Go 1.26 or newer

## Build

```sh
go build -o crawler
```

## Usage

```sh
./crawler BASE_URL MAX_CONCURRENCY MAX_PAGES
```

- `BASE_URL` — the site to crawl (e.g. `https://example.com`)
- `MAX_CONCURRENCY` — max number of pages fetched in parallel
- `MAX_PAGES` — stop after this many pages have been crawled

The crawler only follows links on the same domain as `BASE_URL`.

### Example

```sh
./crawler https://example.com 5 50
```

This writes a `report.json` file in the current directory containing one entry
per crawled page:

```json
[
  {
    "url": "https://example.com/",
    "heading": "Example Heading",
    "first_paragraph": "First paragraph text...",
    "outgoing_links": ["https://example.com/about"],
    "image_urls": ["https://example.com/logo.png"]
  }
]
```

## Tests

```sh
go test ./...
```

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	if len(pages) == 0 {
		fmt.Println("No data to write to JSON")
		return nil
	}

	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedPages := make([]PageData, 0, len(pages))
	for _, k := range keys {
		sortedPages = append(sortedPages, pages[k])
	}

	data, err := json.MarshalIndent(sortedPages, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	fmt.Printf("Report written to %s\n", filename)

	return nil
}

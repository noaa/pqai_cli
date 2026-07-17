package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// printJSON pretty-prints raw JSON bytes to stdout.
func printJSON(raw []byte) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		os.Stdout.Write(raw)
		fmt.Println()
		return
	}
	fmt.Println(buf.String())
}

// searchResult mirrors the fields of interest in a PQAI search result.
type searchResult struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Score    float64  `json:"score"`
	Date     string   `json:"publication_date"`
	WWW      string   `json:"www_link"`
	Owner    string   `json:"owner"`
	Snippet  string   `json:"snippet"`
	Index    string   `json:"index"`
	Inventor []string `json:"inventors"`
}

type searchResponse struct {
	Results []searchResult `json:"results"`
	Query   string         `json:"query"`
}

// printSearchResults renders a human-readable list of search results.
func printSearchResults(raw []byte) {
	var resp searchResponse
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Results) == 0 {
		printJSON(raw)
		return
	}
	for i, r := range resp.Results {
		fmt.Printf("%2d. %-16s score=%.4f  %s\n", i+1, r.ID, r.Score, r.Date)
		if r.Title != "" {
			fmt.Printf("    %s\n", r.Title)
		}
		if r.Owner != "" {
			fmt.Printf("    owner: %s\n", r.Owner)
		}
		if r.Snippet != "" {
			fmt.Printf("    snippet: %s\n", truncate(r.Snippet, 200))
		} else if r.Abstract != "" {
			fmt.Printf("    %s\n", truncate(r.Abstract, 200))
		}
		fmt.Println()
	}
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}

package main

import (
	"encoding/json"
	"fmt"
)

func ExtractLinks(text string) []string {
	return urlRegex.FindAllString(text, -1)
}

func ParseSearchResults(body []byte) ([]SearxResult, error) {
	var resp SearxResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse searxng response: %w", err)
	}
	return resp.Results, nil
}

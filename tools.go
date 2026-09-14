package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TODO: Not hardcode endpoint
func InternetSearch(httpClient *http.Client, query []byte, reqURL string) ([]SearxResult, error) {
	params := url.Values{}
	params.Set("q", string(query))
	params.Set("format", "json")

	request, err := http.NewRequest(http.MethodGet, reqURL+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("searxng returned status %d: %s", response.StatusCode, body)
	}

	return ParseSearchResults(body)
}
func FetchHTML(httpClient *http.Client, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	// Some sites block requests without a User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoFetcher/1.0)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

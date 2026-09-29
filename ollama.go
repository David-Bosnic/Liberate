package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CallOllama(httpClient *http.Client, url string, payload OllamaPayload) (*OllamaOut, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpRequest.Header.Set("accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", response.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	out := &OllamaOut{}
	if err := json.Unmarshal(respBody, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OllamaOut: %w", err)
	}

	return out, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
)

const (
	POOR  = "1"
	BAD   = "2"
	OK    = "3"
	GOOD  = "4"
	GREAT = "5"
)

func RunLiberate(httpClient *http.Client, userPrompt string) ([]string, error) {

	ollamaEndpoint := ENV.OllamaAPI
	searxngEndpoint := ENV.SearXNGAPI

	context := fmt.Sprintf(PromptMakeLink, userPrompt)

	request := OllamaPayload{
		Model: "qwen2.5:3b",
		Messages: []Message{
			{Role: "user", Content: context},
		},
		Options: Options{
			Temperature: 0.2,
		},
		Think:  false,
		Stream: false,
	}
	out, err := CallOllama(httpClient, ollamaEndpoint, request)
	if err != nil {
		return nil, fmt.Errorf("Failed to call Ollama: %e", err)
	}
	searchParam := strings.Trim(out.Message.Content, `"`)
	fmt.Println("Search:", searchParam)

	searxResults, err := InternetSearch(httpClient, []byte(searchParam), searxngEndpoint)
	if err != nil {
		return nil, fmt.Errorf("Failed to internetSearch: %e", err)
	}
	if len(searxResults) == 0 {
		return nil, fmt.Errorf("searxResults did not return any links")
	}

	var links []string
	for i := range searxResults {
		htmlContent, err := readability.FromURL(searxResults[i].URL, time.Second*30)
		if err != nil {
			Logger.Println("Failed to get content", err)
			continue
		}
		var buf bytes.Buffer
		if err := htmlContent.RenderText(&buf); err != nil {
			Logger.Println("Failed to render", searxResults[i].URL)
			continue
		}
		//TODO: See if this is actually an issue with performance later on
		if len(buf.String()) >= 4000 {
			continue
		}
		builtPrompt := fmt.Sprintf(PromptScale, userPrompt, buf.String())
		request = OllamaPayload{
			Model: "qwen2.5:3b",
			Messages: []Message{
				{Role: "user", Content: builtPrompt},
			},
			Options: Options{
				Temperature: 1,
				NumCtx:      8192,
			},
			Format: json.RawMessage(`{"type":"string","enum":["1","2","3","4","5"]}`),
			Think:  false,
			Stream: false,
		}
		out, err = CallOllama(httpClient, ollamaEndpoint, request)
		if err != nil {
			return nil, fmt.Errorf("Failed to call Ollama: %e", err)
		}
		rating := strings.Trim(out.Message.Content, `"`)

		switch rating {
		case GREAT, GOOD, OK:
			links = append(links, searxResults[i].URL)
			//TODO: Just making this short in the meantime. This should be streamed to the user via a channel
			// if len(links) == 3 {
			// 	return links, nil
			// }
		case BAD, POOR:
			continue
		}
	}
	return links, nil
}

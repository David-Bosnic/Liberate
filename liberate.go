package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

func RunLiberate(httpClient *http.Client, userPrompt string) ([]LinkRating, error) {

	ollamaEndpoint := ENV.OllamaAPI
	searxngEndpoint := ENV.SearXNGAPI

	context := fmt.Sprintf(PromptMakeLink, userPrompt)

	request := OllamaPayload{
		Model: "qwen3.6:35b",
		Messages: []Message{
			{Role: "user", Content: context},
		},
		Options: Options{
			//Lower Temperature to improve consistancy
			Temperature: 0.2,
		},
		// Think off for speed improvments.
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

	var links []LinkRating
	var linksMutex sync.Mutex
	var wg sync.WaitGroup

	for i := range searxResults {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			htmlContent, err := readability.FromURL(searxResults[index].URL, time.Second*30)
			if err != nil {
				Logger.Println("Failed to get content", err)
				return
			}
			var buf bytes.Buffer
			if err := htmlContent.RenderText(&buf); err != nil {
				Logger.Println("Failed to render", searxResults[index].URL)
				return
			}
			//TODO: Need to see if this is actually an issue with performance. Might be better
			// to do a partial parsing to not ignore larger sites
			if len(buf.String()) >= 4000 {
				return
			}
			builtPrompt := fmt.Sprintf(PromptScale, userPrompt, buf.String())
			request := OllamaPayload{
				Model: "qwen3.6:35b",
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
			out, err := CallOllama(httpClient, ollamaEndpoint, request)
			if err != nil {
				Logger.Printf("Failed to call Ollama for %s: %v", searxResults[index].URL, err)
				return
			}
			rating := strings.Trim(out.Message.Content, `"`)

			switch rating {
			case GREAT, GOOD, OK:
				linksMutex.Lock()
				links = append(links, LinkRating{
					URL:    searxResults[index].URL,
					Rating: rating,
				})
				linksMutex.Unlock()
			case BAD, POOR:
				return
			}
		}(i)
	}
	wg.Wait()
	return links, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
)

func RunLiberate(httpClient *http.Client) {

	//TODO: make this flexable to the .ENV and not hardcoded
	ollamaEndpoint := ENV.OllamaAPI
	searxngEndpoint := ENV.SearXNGAPI

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Invalid amount of arguments need 2 got", len(args))
		os.Exit(1)
	}

	context := fmt.Sprintf(PromptMakeLink, args[1])

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
		fmt.Println("Failed to call Ollama:", err)
		os.Exit(1)
	}
	searchParam := strings.Trim(out.Message.Content, `"`)
	fmt.Println("Search:", searchParam)

	searxResults, err := InternetSearch(httpClient, []byte(searchParam), searxngEndpoint)
	if err != nil {
		fmt.Println("Failed to internetSearch", err)
		os.Exit(1)
	}
	if len(searxResults) == 0 {
		fmt.Println("searxResults did not return any links")
		os.Exit(0)
	}

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
		if len(buf.String()) >= 4000 {
			continue
		}
		prompt := fmt.Sprintf(PromptScale, args[1], buf.String())
		request = OllamaPayload{
			Model: "qwen2.5:3b",
			Messages: []Message{
				{Role: "user", Content: prompt},
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
			fmt.Println("Failed to call Ollama:", err)
			os.Exit(1)
		}
		rating := strings.Trim(out.Message.Content, `"`)
		switch rating {
		case GREAT, GOOD, OK:
			fmt.Printf("%s/5: %s\n", rating, searxResults[i].URL)
		case BAD, POOR:
			continue
		}
	}
}

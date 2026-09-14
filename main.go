package main

import (
	"bufio"
	"bytes"

	"fmt"

	"log"
	"net/http"

	"os"
	"regexp"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
)

// NOTE: Will implement enum for the 2-3 tools calls require to make this work
// type ToolCall int
//
// const (
//
//	InternetSearch ToolCall = iota
//	ReadPage
//
// )

var ENV Env
var urlRegex = regexp.MustCompile(`https?://[^\s\)\]]+`)

func init() {
	ENV = InitEnv()
}

func InitEnv() Env {
	//Reads local file and does an error check
	file, err := os.Open("./.env")
	if err != nil {
		fmt.Printf("Err %v\n", err)
		os.Exit(1)
	}
	//Making sure the file is closed out
	defer file.Close()

	//Creating a small buffer that can be filled with the file information
	scanner := bufio.NewScanner(file)

	//Reading .env file by line
	var ollamaAPI, searcxngAPI string

	for scanner.Scan() {
		if scanner.Err() != nil {
			panic(scanner.Err().Error())
		}
		line := scanner.Text()
		if strings.Contains(line, "OLLAMA_API") {
			ollamaAPI = strings.TrimPrefix(line, "OLLAMA_API=")
		}
		if strings.Contains(line, "SEARXNG_API") {
			searcxngAPI = strings.TrimPrefix(line, "SEARXNG_API=")
		}

	}
	env := Env{
		OllamaAPI:  ollamaAPI,
		SearXNGAPI: searcxngAPI,
	}
	return env
}

func main() {
	httpClient := &http.Client{}

	//TODO: make this flexable to the .ENV and not hardcoded
	ollamaEndpoint := ENV.OllamaAPI
	searxngEndpoint := ENV.SearXNGAPI

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Invalid amount of arguments need 2 got", len(args))
		os.Exit(1)
	}
	context := args[1]

	request := OllamaPayload{
		Model: "qwen2.5:3b",
		Messages: []Message{
			{Role: "user", Content: context},
		},
		Think:  false,
		Stream: false,
		Tools: []Tool{
			{
				Type: "function",
				Function: Function{
					Name:        "internetSearch",
					Description: "Search the internet for current information.",
					Parameters: Parameters{
						Type: "object",
						Properties: map[string]Property{
							"query": {
								Type:        "string",
								Description: "The search query.",
							},
						},
						Required: []string{"query"},
					},
				},
			},
		},
	}

	out, err := CallOllama(httpClient, ollamaEndpoint, request)
	if err != nil {
		fmt.Println("Failed to call Ollama:", err)
		os.Exit(1)
	}
	context = context + out.Message.Content
	fmt.Println(out.Message.Content)
	searxResults, err := InternetSearch(httpClient, []byte(out.Message.ToolCalls[0].Function.Arguments.Query), searxngEndpoint)
	if err != nil {
		fmt.Println("Failed to internetSearch", err)
		os.Exit(1)
	}
	result := searxResults[0]

	// htmlContent, err := fetchHTML(httpClient, result.URL)
	htmlContent, err := readability.FromURL(result.URL, time.Second*30)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	if err := htmlContent.RenderHTML(&buf); err != nil {
		log.Fatal(err)
	}
	context = context + "Search results\n\n" + buf.String() + "Get me 3 sources no more no less **Just return the links no other information is required**"

	request = OllamaPayload{
		Model: "qwen2.5:3b",
		Messages: []Message{
			{Role: "user", Content: context},
		},
		Think:  false,
		Stream: false,
		Tools: []Tool{
			{
				Type: "function",
				Function: Function{
					Name:        "internetSearch",
					Description: "Search the internet for current information.",
					Parameters: Parameters{
						Type: "object",
						Properties: map[string]Property{
							"query": {
								Type:        "string",
								Description: "The search query.",
							},
						},
						Required: []string{"query"},
					},
				},
			},
		},
	}
	out, err = CallOllama(httpClient, ollamaEndpoint, request)
	if err != nil {
		fmt.Println("Failed to call Ollama:", err)
		os.Exit(1)
	}
	links := ExtractLinks(out.Message.Content)
	for _, link := range links {
		fmt.Println(link)
	}
}

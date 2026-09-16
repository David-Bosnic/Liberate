package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"

	"fmt"

	"net/http"

	"os"
	"regexp"
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

var PromptMakeLink = `
You are a search query generator. Your only job is to convert the user's question into a single, effective Google search query.
Rules:
- Output ONLY the search query, wrapped in double quotes.
- Do NOT answer the question.
- Do NOT explain anything.
- Do NOT add commentary before or after.
- Keep the query short (3-8 words), using the terms someone would actually type into Google.
Examples:
Question: How do I do a print statement in python
Query: "python print statement syntax example"
Question: What's the best way to center a div in CSS
Query: "how to center a div css"
Question: Why does my docker container keep crashing on startup
Query: "docker container crashing on startup fix"

Now convert this question:
Question: %s
Query:
`

// TODO: The best way is to make this not a prompt and fully mechanical
var PromptFragment = `
	You are a URL text-fragment generator. Given a page's raw text content and a phrase or claim the user wants highlighted, output a URL using the browser text-fragment feature.

Rules:
- Output format: {base_url}#:~:text={encoded_text}
- Find the EXACT matching phrase from the provided page text — do not paraphrase or reword it.
- Encode spaces as %20. Encode other special characters using standard URL encoding.
- For a range of text (start to end), format as: #:~:text={encoded_start},{encoded_end}
- Keep the highlighted phrase as short as possible while still uniquely identifying the target text — ideally under 15 words.
- Output ONLY the final URL. No explanation, no commentary, no markdown formatting.
- If the requested phrase does NOT appear verbatim in the provided page text, output exactly: NO_MATCH_FOUND

Examples:

Base URL: https://example.com/article
Page text: "The quick brown fox jumps over the lazy dog. It was a sunny afternoon."
Request: highlight the part about the fox jumping
Output: https://example.com/article#:~:text=The%20quick%20brown%20fox%20jumps%20over%20the%20lazy%20dog

Base URL: https://example.com/article
Page text: "Sales grew steadily in Q1. By Q2, revenue had doubled compared to last year."
Request: highlight from Q1 growth through the Q2 doubling
Output: https://example.com/article#:~:text=Sales%20grew%20steadily%20in%20Q1,revenue%20had%20doubled%20compared%20to%20last%20year

`
var Logger *log.Logger

func init() {
	file, err := os.OpenFile("liberate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "", 0)
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
		prompt := fmt.Sprintf(`Rate how well this page answers the question, using this scale:
		1 = Completely unrelated or no usable content
		2 = Barely related, doesn't address the question
		3 = Partially addresses the question but incomplete
		4 = Mostly answers the question with minor gaps
		5 = Directly and completely answers the question

		Question: %s

		Site content:
		%s

		Respond with only the number.`, args[1], buf.String())
		request = OllamaPayload{
			Model: "qwen2.5:3b",
			Messages: []Message{
				{Role: "user", Content: prompt},
			},
			Options: Options{
				Temperature: 0.2,
				NumCtx:      8192,
			},
			Format: json.RawMessage(`{"type":"string","enum":["1","2","3","4","5"]}`),
			// Format: json.RawMessage(`{"type":"string","enum":["YES","NO","NO-CONTENT"]}`),
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

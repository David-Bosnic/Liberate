package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
)

type ToolCall int

const (
	InternetSearch ToolCall = iota
)

type Env struct {
	Port string
}

//NOTE: There needs to be a memory held between prompts
//since ollama will reset with each response

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Tools []Tool

type OllamaPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Think    bool      `json:"think"`
	Stream   bool      `json:"stream"`
	Tools    Tools     `json:"tools"`
}

type OllamaOut struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role      string `json:"role"`
		Content   string `json:"content"`
		ToolCalls []struct {
			ID       string `json:"id"`
			Function struct {
				Index     int    `json:"index"`
				Name      string `json:"name"`
				Arguments struct {
					Query string `json:"query"`
				} `json:"arguments"`
			} `json:"function"`
		} `json:"tool_calls"`
	} `json:"message"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	TotalDuration      int    `json:"total_duration"`
	LoadDuration       int    `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int    `json:"eval_duration"`
}
type SearxResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type SearxResponse struct {
	Query   string        `json:"query"`
	Results []SearxResult `json:"results"`
}

func main() {
	env, err := InitEnv()
	if err != nil {
		panic(err)
	}
	httpClient := &http.Client{}

	//TODO: make this flexable to the .ENV and not hardcoded
	ollamaEndpoint := fmt.Sprintf("http://david-framework:%s/api/chat", env.Port)

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

	out, err := callOllama(httpClient, ollamaEndpoint, request)
	if err != nil {
		fmt.Println("Failed to call Ollama:", err)
		os.Exit(1)
	}
	context = context + out.Message.Content
	fmt.Println(out.Message.Content)
	searxResults, err := internetSearch(httpClient, []byte(out.Message.ToolCalls[0].Function.Arguments.Query))
	if err != nil {
		fmt.Println("Failed to internetSearch", err)
		os.Exit(1)
	}
	result := searxResults[0]
	fmt.Printf("%s: %s\n\n", result.Title, result.URL)

	// htmlContent, err := fetchHTML(httpClient, result.URL)
	htmlContent, err := readability.FromURL(result.URL, time.Second*30)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	if err := htmlContent.RenderHTML(&buf); err != nil {
		log.Fatal(err)
	}
	context = context + "Search results\n\n" + buf.String() + "Just return the links no other information is required"

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
	out, err = callOllama(httpClient, ollamaEndpoint, request)
	if err != nil {
		fmt.Println("Failed to call Ollama:", err)
		os.Exit(1)
	}
	fmt.Println(out.Message.Content)
}
func callOllama(httpClient *http.Client, url string, payload OllamaPayload) (*OllamaOut, error) {
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

// TODO: Not hardcode endpoint
func internetSearch(httpClient *http.Client, query []byte) ([]SearxResult, error) {
	params := url.Values{}
	params.Set("q", string(query))
	params.Set("format", "json")

	reqURL := "http://david-framework:8080/search?" + params.Encode()

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
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

	return parseSearchResults(body)
}
func fetchHTML(httpClient *http.Client, url string) (string, error) {
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

func InitEnv() (Env, error) {
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
	port := ""

	for scanner.Scan() {
		if scanner.Err() != nil {
			panic(scanner.Err().Error())
		}
		line := scanner.Text()
		if strings.Contains(line, "OLLAMA_PORT") {
			port = strings.TrimPrefix(line, "OLLAMA_PORT=")
		}

	}
	env := Env{
		Port: port,
	}
	return env, nil
}

func parseSearchResults(body []byte) ([]SearxResult, error) {
	var resp SearxResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse searxng response: %w", err)
	}
	return resp.Results, nil
}

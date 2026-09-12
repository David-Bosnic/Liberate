package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	url := fmt.Sprintf("http://david-framework:%s/api/chat", env.Port)

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Invalid amount of arguments need 2 got", len(args))
		os.Exit(1)
	}
	userInput := args[1]

	request := OllamaPayload{
		Model: "qwen3.6:35b",
		Messages: []Message{
			{Role: "user", Content: userInput},
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
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Failed to Marshal: ", err)
		os.Exit(1)
	}
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	httpRequest.Header.Set("accept", "application/json")
	//TODO: handle failed fetch error
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	out := OllamaOut{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		fmt.Println("Failed to Unmarshal OllamaOut: ", err)
		os.Exit(1)
	}
	fmt.Println(out)
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

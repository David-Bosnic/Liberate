package main

import "time"

type Env struct {
	OllamaAPI  string
	SearXNGAPI string
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

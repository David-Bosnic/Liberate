package main

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"os"
	"regexp"
)

// NOTE: Tool calls are under the llms discretion. More likely forcing situtation
// will be better for consistancy.

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
var Logger *log.Logger

func init() {
	//Logger for none fatal errors mostly
	file, err := os.OpenFile("liberate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "", 0)
	ENV, err = InitEnv()
	if err != nil {
		fmt.Println("Failed to init env:", err)
		os.Exit(1)
	}
}

func main() {
	httpClient := &http.Client{}

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api", liberateHandler(httpClient))

	fmt.Printf("Running Liberate on http://localhost%s\n", ENV.FrontendPort)
	log.Fatal(http.ListenAndServe(ENV.FrontendPort, nil))

}

func liberateHandler(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LiberateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Prompt == "" {
			http.Error(w, "prompt is required", http.StatusBadRequest)
			return
		}
		links, err := RunLiberate(client, req.Prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(LiberateResponse{Links: links})
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"os"
	"regexp"
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
var Logger *log.Logger

func init() {
	//Logger for none fatal errors mostly
	file, err := os.OpenFile("liberate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "", 0)
	ENV = InitEnv()
}

func main() {
	httpClient := &http.Client{}

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api", liberateHandler(httpClient))

	fmt.Println("Running Liberate on port :8081")
	http.ListenAndServe(":8081", nil)
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

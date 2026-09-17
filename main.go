package main

import (
	"log"

	"net/http"

	"os"
	"regexp"
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
	RunLiberate(httpClient)
}

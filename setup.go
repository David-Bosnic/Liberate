package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InitEnv() (Env, error) {
	//Reads local file and does an error check
	file, err := os.Open("./.env")
	if err != nil {
		return Env{}, fmt.Errorf("Err %v\n", err)
	}
	//Making sure the file is closed out
	defer file.Close()

	//Creating a small buffer that can be filled with the file information
	scanner := bufio.NewScanner(file)

	//Reading .env file by line
	var ollamaAPI, searcxngAPI, frontendPort string

	for scanner.Scan() {
		if scanner.Err() != nil {
			return Env{}, scanner.Err()
		}
		line := scanner.Text()
		if strings.Contains(line, "OLLAMA_API") {
			ollamaAPI = strings.TrimPrefix(line, "OLLAMA_API=")
		}
		if strings.Contains(line, "SEARXNG_API") {
			searcxngAPI = strings.TrimPrefix(line, "SEARXNG_API=")
		}
		if strings.Contains(line, "FRONTEND_PORT") {
			frontendPort = strings.TrimPrefix(line, "FRONTEND_PORT=")
		}

	}
	if ollamaAPI == "" || searcxngAPI == "" || frontendPort == "" {
		return Env{}, fmt.Errorf("failed to recieve all .env information")
	}
	env := Env{
		OllamaAPI:    ollamaAPI,
		SearXNGAPI:   searcxngAPI,
		FrontendPort: frontendPort,
	}
	return env, nil
}

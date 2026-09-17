package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func readInPrompt(filePath string) (*string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	contentStr := string(content)

	return &contentStr, err
}

func readInDocs(filePath string) (*map[string]interface{}, error) {
	// Read json file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Define docs map to store key, val lookup
	docs := make(map[string]interface{})

	// load json file into map
	if err := json.Unmarshal(content, &docs); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return &docs, nil
}

func loadInVars() {
	// Call readInDocs function and assign the returned value to the global variable loadedDocs
	docs, errDoc := readInDocs("./data/hw_docs.json")
	if errDoc != nil {
		fmt.Println("Error loading docs:", errDoc)
		return
	}
	loadedDocs = docs

	// Read in Prompt
	prompt, errPrompt := readInPrompt("./data/prompt.txt")
	if errPrompt != nil {
		fmt.Println("Error loading prompt:", errPrompt)
		return
	}
	loadedPrompt = *prompt

	// Read in Model
	model, errModel := constructModel("/home/ctyler/llm_models/mistral-7b-instruct-v0.2.Q3_K_S.gguf")
	if errModel != nil {
		fmt.Println("Error loading prompt:", errModel)
		return
	}
	loadedModel = model
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}

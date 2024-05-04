package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/go-skynet/go-llama.cpp"
)

// Global vars
var loadedDocs *map[string]interface{}
var loadedPrompt string
var loadedModel *llama.LLama

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

func processRequest(body []byte, docs *map[string]interface{}) string {
	if docs == nil {
		return ""
	}

	bodyStr := string(body)
	key := parseKeyFromBody(bodyStr)

	var modelInput string
	if value, ok := (*docs)[key]; ok {
		if strValue, ok := value.(string); ok {
			modelInput = loadedPrompt + bodyStr + strValue
		}
	}

	fmt.Println(modelInput)

	out, err := loadedModel.Predict(modelInput)
	if err != nil {
		fmt.Println("ERROR predicting with model:", err)
	}

	return out
}

func parseKeyFromBody(bodyStr string) string {
	for key := range *loadedDocs {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(key))

		if re.MatchString(bodyStr) {
			return key
		}
	}
	return ""
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	out := processRequest(body, loadedDocs)

	fmt.Printf("got / request. Body: %s\n", out)

	// Write response with success status code
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "This is my website!\n")

}

func loadInVars() {
	// Call readInDocs function and assign the returned value to the global variable loadedDocs
	docs, errDoc := readInDocs("./hw_docs.json")
	if errDoc != nil {
		fmt.Println("Error loading docs:", errDoc)
		return
	}
	loadedDocs = docs

	// Read in Prompt
	prompt, errPrompt := readInPrompt("./prompt.txt")
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

func main() {
	loadInVars()

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server Close\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

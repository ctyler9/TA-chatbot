package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

// Global variable to store loaded docs
var loadedDocs *map[string]interface{}

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

	if value, ok := (*docs)[key]; ok {
		if strValue, ok := value.(string); ok {
			return bodyStr + strValue
		}

	}
	return ""
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

func main() {
	// Call readInDocs function and assign the returned value to the global variable loadedDocs
	docs, errDoc := readInDocs("./hw_docs.json")
	if errDoc != nil {
		fmt.Println("Error loading docs:", errDoc)
		return
	}
	loadedDocs = docs

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

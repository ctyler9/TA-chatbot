package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/go-skynet/go-llama.cpp"
)

// Global var init
var loadedDocs *map[string]interface{}
var loadedPrompt string
var loadedModel *llama.LLama

var keyOut = make(map[string]string)
var channel = make(chan HTTPRequestPayload, 1000)

type HTTPRequestPayload struct {
	IdHash   string
	Question string
}

type HTTPResponsePayload struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type OutResponse struct {
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// PROCESS REQUEST
func filterOutput(output *string, strValue *string) {
	// TODO: Still want to find most efficient way to do this
}

func parseKeyFromQuestion(question string) string {
	for key := range *loadedDocs {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(key))

		if re.MatchString(question) {
			return key
		}
	}
	return ""
}
func processDataFromChannel() {
	for {
		select {
		case payload := <-channel:
			processPayload(payload)
		}
	}
}

func processPayload(payload HTTPRequestPayload) {
	id := payload.IdHash
	question := payload.Question

	fmt.Println(question)

	key := parseKeyFromQuestion(question)

	var modelInput string
	if value, ok := (*loadedDocs)[key]; ok {
		if strValue, ok := value.(string); ok {
			modelInput = strValue + loadedPrompt + question
		}
	}

	fmt.Println(modelInput)

	out, err := loadedModel.Predict(modelInput)
	if err != nil {
		fmt.Println("ERROR predicting with model:", err)
	}
	// Once function implemented, prevent "leakage" of answers
	//filteredOut := filterOutput(out)

	// Write to global dict
	keyOut[id] = out

	fmt.Println("Finished writing output to keymap")

}

// HTTP HANDLERS
func submitData(w http.ResponseWriter, r *http.Request) {
	// Read in request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	question := string(body)

	// Add random noise to avoid collison if somehow same question
	body = append(body, randomBytes(32)...)
	hash := md5.Sum(body)
	hashString := hex.EncodeToString(hash[:])

	payload := HTTPRequestPayload{hashString, question}
	channel <- payload

	// Construct JSON response
	response := HTTPResponsePayload{
		Message: "Data submitted successfully",
		ID:      hashString,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("could not marshal JSON response: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write response with success status code and JSON body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	fmt.Println("Request added to queue")

}

func getProcessedData(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	fmt.Println(keyOut)

	// Check if the id exists in the keyOut map
	if value, ok := keyOut[id]; ok {
		response := OutResponse{
			Message: "Success",
			Data:    value,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

		// Delete key/value pair once request sent
		delete(keyOut, id)

	} else {
		// If id not found, inform the client
		response := OutResponse{
			Message: "Error",
			Data:    fmt.Sprintf("ID %s not found", id),
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
	}
}

func main() {
	// Define flags with default values
	docsPath := flag.String("d", "./data/hw_docs.json", "Path to document map")
	promptPath := flag.String("p", "./data/prompt.txt", "Path to prompt")
	modelPath := flag.String("m", "", "Path to model")
	flag.Parse()
	loadInVars(*docsPath, *promptPath, *modelPath)

	go processDataFromChannel()

	mux := http.NewServeMux()
	mux.HandleFunc("/submit_data", submitData)
	mux.HandleFunc("/get_processed_data", getProcessedData)

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server Close\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

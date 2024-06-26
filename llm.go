package main

import (
	"fmt"

	llama "github.com/go-skynet/go-llama.cpp"
)

var (
	threads   = 32
	tokens    = 32000
	gpulayers = 64
	seed      = -1
)

func constructModel(modelPath string) (*llama.LLama, error) {
	l, err := llama.New(modelPath, llama.EnableF16Memory, llama.SetContext(tokens), llama.EnableEmbeddings, llama.SetGPULayers(gpulayers))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		return nil, err
	}
	fmt.Printf("Model loaded successfully.\n")

	return l, nil
}

func predictModel(llm *llama.LLama, input string) (*string, error) {
	out, err := llm.Predict(input)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

//func main() {
//	llmModel, err := constructModel("/home/ctyler/llm_models/mistral-7b-instruct-v0.2.Q3_K_S.gguf")
//	if err != nil {
//		fmt.Println("Construct model error:", err)
//	}
//
//	out, errPred := predictModel(llmModel, "How are you today?")
//	if errPred != nil {
//		fmt.Println("Predict model error:", errPred)
//	}
//
//	fmt.Println(*out)
//
//}

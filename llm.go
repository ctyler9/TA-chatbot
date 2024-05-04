package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	llama "github.com/go-skynet/go-llama.cpp"
)

var (
	threads   = 4
	tokens    = 128
	gpulayers = 0
	seed      = -1
)

func main() {
	var model string

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&model, "m", "/home/ctyler/llm_models/mistral-7b-instruct-v0.2.Q3_K_S.gguf", "path to q4_0.bin model file to load")
	flags.IntVar(&gpulayers, "ngl", 0, "Number of GPU layers to use")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 512, "number of tokens to predict")
	flags.IntVar(&seed, "s", -1, "predict RNG seed, -1 for random seed")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}
	l, err := llama.New(model, llama.EnableF16Memory, llama.SetContext(128), llama.EnableEmbeddings, llama.SetGPULayers(gpulayers))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Model loaded successfully.\n")

	text := "Hello model, this is a test"

	_, errPred := l.Predict(text, llama.Debug, llama.SetTokenCallback(func(token string) bool {
		fmt.Print(token)
		return true
	}), llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetTopK(90), llama.SetTopP(0.86), llama.SetStopWords("llama"), llama.SetSeed(seed))
	if errPred != nil {
		panic(errPred)
	}

}

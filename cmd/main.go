package main

import (
	"fmt"
	"os"
	"dllm"
)

func main() {
	conf := dllm.OpenAIConfig{
		Model:       "gpt-4",
		Temperature: 0.1,
		MaxTokens:   300,
	}
	authToken := os.Getenv("OPENAI_API_KEY")
	openai := dllm.NewOpenAI(authToken, conf)
	stream, err := openai.Ask("Compose a poem about a tree.")
	if err != nil {
		panic(err)
	}
	fmt.Println("Stream received")
	defer stream.Close()
	stream.Read(dllm.Print)
}

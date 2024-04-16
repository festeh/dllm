package dllm

import "os"

func main() {
	conf := OpenAIConfig{
		Model:       "gpt-4",
		Temperature: 0.1,
		MaxTokens:   100,
	}
	authToken := os.Getenv("OPENAI_API_KEY")
	openai := NewOpenAI(authToken, conf)
	openai.Ask("ping")
}

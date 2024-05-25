package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type OpenaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIData struct {
	Model       string          `json:"model"`
	Messages    []OpenaiMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	Stream      bool            `json:"stream"`
}

type OpenAI struct {
	authToken string
	client    *http.Client
}

func NewOpenAI() (*OpenAI, error) {
	authToken := os.Getenv("OPENAI_API_KEY")
	if authToken == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY not set")
	}
	return &OpenAI{authToken, &http.Client{}}, nil
}

func (o *OpenAI) Name() string {
	return "OpenAI"
}

func (o *OpenAI) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.openai.com/v1/chat/completions")
}


func (o *OpenAI) CreateData(query *Query) ([]byte, error) {
	messages := make([]OpenaiMessage, len(query.Messages))
	for i, message := range query.Messages {
		messages[i] = OpenaiMessage{message.Role, message.Content}
	}
	data := &OpenAIData{
		Model:       "gpt-4",
		Messages:    messages,
		Temperature: 0.1,
		MaxTokens:   500,
		Stream:      true,
	}
	return json.Marshal(data)
}

func (o *OpenAI) addHeaders(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+o.authToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
}

func (o *OpenAI) do(request *http.Request) (*http.Response, error) {
	return o.client.Do(request)
}

type Chunk struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index int   `json:"index"`
	Delta Delta `json:"delta"`
}

type Delta struct {
	Content string `json:"content"`
}

func (o *OpenAI) GetWriterCallback() func([]byte) ([]byte, bool) {
	return func(chunk []byte) ([]byte, bool) {
		// TODO: move to zerolog
		// log.Print("Received chunk", string(chunk))
		// skip data: prefix
		chunk = chunk[6:]
		done := []byte("[DONE]")
		if bytes.Equal(chunk, done) {
			// fmt.Println("Received DONE")
			return nil, true
		}
		deserialized := Chunk{}
		err := json.Unmarshal(chunk, &deserialized)
		if err != nil {
			fmt.Println("Error deserializing chunk", err)
			fmt.Println("chunk", string(chunk))
			return nil, true
		}
		if len(deserialized.Choices) == 0 {
			return nil, true
		}
		choice := deserialized.Choices[0]
		return []byte(choice.Delta.Content), false
	}
}

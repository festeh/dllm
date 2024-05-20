package dllm

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicData struct {
	model       string
	messages    []AnthropicMessage
	temperature float32
	maxTokens   int
	stream      bool
}

type Anthropic struct {
	key    string
	client *http.Client
}

func NewAnthropic() (*Anthropic, error) {
	keyName := "ANTHROPIC_API_KEY"
	key, ok := os.LookupEnv(keyName)
	if !ok {
		return nil, fmt.Errorf("Environment variable %s not set", keyName)
	}
	return &Anthropic{key, &http.Client{}}, nil
}

func (a *Anthropic) Name() string {
	return "Anthropic"
}

func (a *Anthropic) addHeaders(request *http.Request) {
	request.Header.Set("x-api-key", a.key)
	request.Header.Set("anthropic-version", "2023-06-01")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-Control", "no-cache")
}

func (a *Anthropic) GetStream(query *Query, writer StreamWriter) (*Stream, error) {
	return NewStream(query, writer, a)
}

func (a *Anthropic) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.anthropic.com/v1/messages")
}


func (a *Anthropic) CreateData(query *Query) *AnthropicData {
	messages := make([]AnthropicMessage, len(query.Messages))
	for i, message := range query.Messages {
		messages[i] = AnthropicMessage{message.Role, message.Content}
	}
	data := &AnthropicData{
		model:       "claude-3-opus-20240229",
		messages:    messages,
		temperature: 0.1,
		maxTokens:   300,
		stream:      true,
	}
	return data
}

func (a *Anthropic) do(request *http.Request) (*http.Response, error) {
	return a.client.Do(request)
}

func (a *Anthropic) GetWriterCallback() func([]byte) ([]byte, bool) {
	return func(body []byte) ([]byte, bool) {
		return body, false
	}
}

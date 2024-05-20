package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type OpenaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestParams struct {
	Model       string          `json:"model"`
	Messages    []OpenaiMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	Stream      bool            `json:"stream"`
}

type OpenAIConfig struct {
	Model       string
	Temperature float64
	MaxTokens   int
}

type OpenAI struct {
	authToken string
	config    OpenAIConfig
	client    *http.Client
}

func NewOpenAI(authToken string, config OpenAIConfig) *OpenAI {
	return &OpenAI{authToken, config, &http.Client{}}
}

func (o *OpenAI) AddHeaders(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+o.authToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
}

func (o *OpenAI) GetStream(body []byte, writer http.ResponseWriter) (*Stream, error) {
	request := NewPostRequest(o.CompletionURL())
	o.AddHeaders(request)
	query := OpenAIQuery{}
	err := LoadQuery(body, &query)
	if err != nil {
		return nil, err
	}
	return NewStream(request, writer)
}

func (o *OpenAI) CreateRequest(params RequestParams) (request *http.Request, err error) {
	// create request
	parsedURL, err := url.Parse("https://api.openai.com/v1/chat/completions")
	request = &http.Request{
		Method: "POST",
		URL:    parsedURL,
		Header: make(http.Header),
	}
	request.Header.Set("Authorization", "Bearer "+o.authToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	bufData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	request.Body = io.NopCloser(bytes.NewBuffer(bufData))
	request.ContentLength = int64(len(bufData))
	return request, nil
}

func (o *OpenAI) Ask(messages []OpenaiMessage) (stream *Stream, err error) {
	params := RequestParams{
		Model:       o.config.Model,
		Messages:    messages,
		Temperature: o.config.Temperature,
		MaxTokens:   o.config.MaxTokens,
		Stream:      true,
	}
	request, err := o.CreateRequest(params)
	if err != nil {
		return nil, err
	}
	response, err := o.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s. Message: %s", response.Status, response.Body)
	}
	return &Stream{response, w}, nil

}

func (o *OpenAI) Name() string {
	return "OpenAI"
}

func (o *OpenAI) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.openai.com/v1/chat/completions")
}

package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestParams struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
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

func (o *OpenAI) CreateHandler() (handler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle request
	}
}

func (o *OpenAI) CreateRequest(params RequestParams) (request *http.Request, err error) {
	// create request
	request = &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "https://api.openai.com/v1/chat/completions"},
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

type Stream struct {
	response *http.Response
}

func (o *OpenAI) Ask(query string) (stream* Stream, err error) {
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: query},
	}
	params := RequestParams{
		Model:       o.config.Model,
		Messages:    messages,
		Temperature: o.config.Temperature,
		MaxTokens:   o.config.MaxTokens,
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
	return &Stream{response}, nil

}

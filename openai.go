package dllm

import (
	"net/http"
	"net/url"
)

type RequestParams struct {
	model string `json:"model"`
}

type OpenAIConfig struct {
	model       string
	temperature float64
}

type OpenAI struct {
	authToken string
	config    OpenAIConfig
}

func NewOpenAI(authToken string, config OpenAIConfig) *OpenAI {
	return &OpenAI{authToken, config}
}

func (o *OpenAI) CreateHandler() (handler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle request
	}
}

func (o *OpenAI) CreateRequest() (request *http.Request) {
	// create request
	request = &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "https://api.openai.com/v1/chat/completions"},
	}
	request.Header.Set("Authorization", "Bearer "+o.authToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	// create request
	return request
}

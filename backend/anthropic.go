package dllm

import (
	"net/http"
	"net/url"
)

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicQuery struct {
	Messages []AnthropicMessage `json:"messages"`
}

type Anthropic struct {
	key string
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

func (a *Anthropic) GetStream(body []byte, writer http.ResponseWriter) (*Stream, error) {
	request := NewPostRequest(a.CompletionURL())
	a.addHeaders(request)
	query := AnthropicQuery{}
	err := LoadQuery(body, &query)
	if err != nil {
		return nil, err
	}
	return NewStream(body, writer)
}

func (a *Anthropic) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.anthropic.com/v1/messages")
}

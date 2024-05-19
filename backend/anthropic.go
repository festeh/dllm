package dllm

import "net/http"

type AnthropicMessage struct {
	Role		string `json:"role"`
	Content	string `json:"content"`
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

func (a *Anthropic) addHeaders(writer http.ResponseWriter) {
	AddHeaders(writer)
	writer.Header().Set("x-api-key", a.key)
}

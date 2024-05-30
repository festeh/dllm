package dllm

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Parameters struct {
	Model       string   `json:"model"`
	Temperature *float32 `json:"temperature"`
	MaxTokens   *int     `json:"max_tokens"`
}

type Query struct {
	Messages   []Message  `json:"messages"`
	Parameters Parameters `json:"parameters"`
}

func LoadQuery(b []byte, query *Query) (err error) {
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(query)
	if query.Parameters.Temperature == nil {
		query.Parameters.Temperature = new(float32)
		*query.Parameters.Temperature = 0.1
	}
	if query.Parameters.MaxTokens == nil {
		query.Parameters.MaxTokens = new(int)
		*query.Parameters.MaxTokens = 4000
	}
	return
}

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
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens  int     `json:"max_tokens"`
}

type Query struct {
	Messages   []Message  `json:"messages"`
	Parameters Parameters `json:"parameters"`
}

func LoadQuery(b []byte, query *Query) (err error) {
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(query)
	return
}

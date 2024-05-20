package dllm

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}


type Query struct {
	Messages []Message `json:"messages"`
}

func LoadQuery(b []byte, query *Query) (err error) {
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(query)
	return
}

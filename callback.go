package dllm

import (
	"encoding/json"
	"fmt"
)

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

func Print(chunk []byte) {
	// skip data: prefix
	chunk = chunk[6:]
	deserialized := Chunk{}
	err := json.Unmarshal(chunk, &deserialized)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(deserialized.Choices[0].Delta.Content)
}

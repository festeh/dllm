package dllm

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetWriterCallback(writer http.ResponseWriter) func([]byte) {
	return func(chunk []byte) {
		// skip data: prefix
		chunk = chunk[6:]
		deserialized := Chunk{}
		err := json.Unmarshal(chunk, &deserialized)
		if err != nil {
			fmt.Println(err)
		}
		writer.Write([]byte(deserialized.Choices[0].Delta.Content))
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

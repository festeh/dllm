package dllm

import (
	"bytes"
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
		done := []byte("DONE")
		if bytes.Equal(chunk, done) {
			fmt.Println("Received DONE")
			return
		}
		deserialized := Chunk{}
		err := json.Unmarshal(chunk, &deserialized)
		if err != nil {
			fmt.Println("Error deserializing chunk", err)
			fmt.Println("chunk", string(chunk))
			return
		}
		if len(deserialized.Choices) == 0 {
			return
		}
		writer.Write([]byte(deserialized.Choices[0].Delta.Content))
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

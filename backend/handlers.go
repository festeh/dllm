package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func addHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}


type OpenAIQuery struct {
	Messages []OpenaiMessage `json:"messages"`
}

func InitOpenAIHandler() http.HandlerFunc {
	authToken := os.Getenv("OPENAI_API_KEY")
	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		if req.Method != "POST" {
			writer.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("Received request")
		conf := OpenAIConfig{
			Model:       "gpt-4",
			Temperature: 0.1,
			MaxTokens:   300,
		}
		openai := NewOpenAI(authToken, conf)
		fmt.Println("Created OpenAI")
		defer req.Body.Close()
		// Body is the JSON with a field "message" that contains the message
		// Decode the JSON and send the message to OpenAI
		b, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Error reading body")
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Body is", string(b))

		query := OpenAIQuery{}
		err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&query)
		if err != nil {
			fmt.Println("Error decoding message")
			fmt.Println(err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Decoded message", query.Messages)
		stream, err := openai.Ask(query.Messages)
		fmt.Println("Sent message to OpenAI")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Stream received")
		defer stream.Close()
		stream.Read(GetWriterCallback(writer))
	}
}

func InitAnthropicHandler() http.HandlerFunc {
	authToken := os.Getenv("ANTHROPIC_API_KEY")
	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		if req.Method != "POST" {
			writer.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("Received request")
		conf := OpenAIConfig{
			Model:       "gpt-4",
			Temperature: 0.1,
			MaxTokens:   300,
		}
		openai := NewOpenAI(authToken, conf)
		fmt.Println("Created OpenAI")
		defer req.Body.Close()
		// Body is the JSON with a field "message" that contains the message
		// Decode the JSON and send the message to OpenAI
		b, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Error reading body")
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Body is", string(b))

		query := OpenAIQuery{}
		err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&query)
		if err != nil {
			fmt.Println("Error decoding message")
			fmt.Println(err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Decoded message", query.Messages)
		stream, err := openai.Ask(query.Messages)
		fmt.Println("Sent message to OpenAI")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Stream received")
		defer stream.Close()
		stream.Read(GetWriterCallback(writer))
	}
}

func DummyHandler(writer http.ResponseWriter, req *http.Request) {
	addHeaders(writer)
	defer req.Body.Close()
	for range 10 {
		writer.Write([]byte("Hello\n"))
		fmt.Println("Sent message")
		if f, ok := writer.(http.Flusher); ok {
			fmt.Println("Flushing")
			f.Flush()
		} else {
			fmt.Println("Unable to convert to flusher")
		}
		// sleep for 500ms
		time.Sleep(500 * time.Millisecond)
	}
	writer.Write([]byte("Goodbye\n"))
	return
}

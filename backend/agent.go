package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Agent interface {
	Name() string
	GetStream() (*Stream, error)
	GetWriterCallback() func([]byte)
}

type Manager struct{}

// TODO: implement a template method pattern
func (m *Manager) CreateHandler(agent Agent) (handler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		AddHeaders(w)
		if r.Method != "POST" {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("Received request")
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := OpenAIQuery{}
		err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&query)
		if err != nil {
			log.Println("Error decoding message")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stream, err := agent.GetStream()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Stream received")
		defer stream.Close()
		stream.Read(agent.GetWriterCallback(w))
	}

}

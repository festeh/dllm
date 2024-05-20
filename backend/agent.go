package dllm

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Agent interface {
	Name() string
	GetStream(body []byte, writer http.ResponseWriter) (*Stream, error)
	GetWriterCallback() func([]byte)
}

type Manager struct{}

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

		stream, err := agent.GetStream(b, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Stream received")

		defer stream.Close()
		stream.Read(agent.GetWriterCallback())
	}

}

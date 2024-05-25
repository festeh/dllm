package dllm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Agent interface {
	Name() string
	GetWriterCallback() func([]byte) ([]byte, bool)
	CreateData(query *Query) ([]byte, error)
	CompletionURL() *url.URL
	addHeaders(request *http.Request)
	do(request *http.Request) (*http.Response, error)
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
		query := &Query{}
		err = LoadQuery(b, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stream, err := NewStream(query, w, agent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Stream received")

		defer stream.Close()
		stream.Read(agent.GetWriterCallback())
	}

}

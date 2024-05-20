package dllm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Agent[Q any, D any] interface {
	Name() string
	GetStream(body []byte, writer http.ResponseWriter) (*Stream, error)
	GetWriterCallback() func([]byte)
	LoadQuery(body []byte) (Q, error)
	CreateData(query Q) D
	CompletionURL() *url.URL
	addHeaders(request *http.Request)
	do(request *http.Request) (*http.Response, error)
}

type Manager[Q any, D any] struct{}

func (m *Manager[Q, D]) CreateHandler(agent Agent[Q, D]) (handler http.HandlerFunc) {
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

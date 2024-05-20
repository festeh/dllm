package dllm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Agent[D any] interface {
	Name() string
	GetStream(query *Query, writer StreamWriter) (*Stream, error)
	GetWriterCallback() func([]byte) ([]byte, bool)
	CreateData(query *Query) D
	CompletionURL() *url.URL
	addHeaders(request *http.Request)
	do(request *http.Request) (*http.Response, error)
}

type Manager[D any] struct{}

func (m *Manager[D]) CreateHandler(agent Agent[D]) (handler http.HandlerFunc) {
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
		stream, err := agent.GetStream(query, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Stream received")

		defer stream.Close()
		stream.Read(agent.GetWriterCallback())
	}

}

package dllm

import (
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
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
	return func(w http.ResponseWriter, request *http.Request) {
		AddHeaders(w)
		if request.Method != "POST" {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Debug().Msg("Received POST request")
		defer request.Body.Close()
		b, err := io.ReadAll(request.Body)
		if err != nil {
			log.Error().Msg(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		query := &Query{}
		err = LoadQuery(b, query)
		if err != nil {
			log.Error().Msg(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stream, err := NewStream(query, w, agent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debug().Msg("Stream created")

		defer stream.Close()
		stream.Read(agent.GetWriterCallback())
	}

}

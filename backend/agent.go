package dllm

import "net/http"

type Agent interface {
	Name() string
	addHeaders(writer http.ResponseWriter)
}

type Manager struct{}

// TODO: implement a template method pattern
func (m *Manager) CreateHandler(agent Agent) (handler http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		agent.addHeaders(w)
	}

}

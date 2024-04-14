package dllm

import (
	"log"
	"net/http"
)

type Route struct {
	path    string
	handler http.HandlerFunc
}

type Server struct {
	handlers []Route
}

func (s *Server) AddRoute(path string, handler http.HandlerFunc) {
	s.handlers = append(s.handlers, Route{path, handler})
}

func (s *Server) Start() {
	for _, route := range s.handlers {
		http.HandleFunc(route.path, route.handler)
	}
	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}

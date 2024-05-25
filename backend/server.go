package dllm

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Route struct {
	path    string
	handler http.HandlerFunc
}

type Server struct {
	Port     int
	handlers []Route
}

func (s *Server) AddRoute(path string, handler http.HandlerFunc) {
	s.handlers = append(s.handlers, Route{path, handler})
}

var VERSION = "development"

func (s *Server) Start() {
	vers := flag.Bool("version", false, "Print the version number")
	v := flag.Bool("v", false, "Print the version number")
	flag.Parse()
	if *vers || *v {
		fmt.Println("Version:", VERSION)
		return
	}
	for _, route := range s.handlers {
		http.HandleFunc(route.path, route.handler)
	}
	fmt.Printf("Server listening on port %d\n", s.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))
}

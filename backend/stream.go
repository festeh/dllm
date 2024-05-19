package dllm

import (
	"bufio"
	"net/http"
)

type Callback func(body []byte)

const BUFFER_SIZE = 30

type Stream struct {
	response *http.Response
	w 			http.ResponseWriter
}

func (s *Stream) Read(callback Callback) {
	reader := bufio.NewReader(s.response.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if len(line) == 0 {
			continue
		}
		callback(line)
		if f, ok := s.w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (s *Stream) Close() {
	s.response.Body.Close()
}

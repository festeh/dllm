package dllm

import (
	"bufio"
	"net/http"
)

type Callback func(body []byte)

const BUFFER_SIZE = 30

type Stream struct {
	response *http.Response
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
	}
}

func (s *Stream) Close() {
	s.response.Body.Close()
}

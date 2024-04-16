package dllm

import "net/http"

type Callback func(body []byte)

const BUFFER_SIZE = 30

type Stream struct {
	response *http.Response
}

func (s *Stream) Read(callback Callback) {
	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := s.response.Body.Read(buf)
		if err != nil {
			break
		}
		callback(buf[:n])
	}
}

func (s *Stream) Close() {
	s.response.Body.Close()
}

package dllm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Callback func(body []byte)

const BUFFER_SIZE = 30

type Stream struct {
	response *http.Response
	w        http.ResponseWriter
}

func NewStream[Q any, D any, A Agent[Q, D]](body []byte, writer http.ResponseWriter, agent A) (*Stream, error) {
	query, err := agent.LoadQuery(body)
	if err != nil {
		return nil, err
	}
	data := agent.CreateData(query)
	dataMarshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request := NewPostRequest(agent.CompletionURL())
	agent.addHeaders(request)
	request.Body = io.NopCloser(bytes.NewBuffer(dataMarshalled))
	request.ContentLength = int64(len(dataMarshalled))
	response, err := agent.do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s. Message: %s", response.Status, response.Body)
	}
	return &Stream{response, writer}, nil
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

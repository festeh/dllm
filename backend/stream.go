package dllm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Callback func(body []byte) ([]byte, bool)

const BUFFER_SIZE = 30

type StreamWriter interface {
	Write([]byte) (int, error)
}

type Stream struct {
	response *http.Response
	writer   StreamWriter
}

func NewStream[ D any, A Agent[D]](query *Query, writer StreamWriter, agent A) (*Stream, error) {
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

func (stream *Stream) Read(callback Callback) {
	reader := bufio.NewReader(stream.response.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if len(line) == 0 {
			continue
		}
		result, isDone := callback(line)
		if isDone {
			break
		}
		stream.writer.Write(result)
		if f, ok := stream.writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (s *Stream) Close() {
	s.response.Body.Close()
}

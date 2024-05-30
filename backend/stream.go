package dllm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Callback func(body []byte) ([]byte, bool)

type StreamWriter interface {
	Write([]byte) (int, error)
}

type Stream struct {
	response *http.Response
	writer   StreamWriter
}

type ModelInfo struct {
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

func NewStream(query *Query, writer StreamWriter, agent Agent) (*Stream, error) {
	dataMarshalled, err := agent.CreateData(query)
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
		bodyString, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Error: %s. Message: %s", response.Status, bodyString)
	}
	modelInfo := ModelInfo{
		Model:       query.Parameters.Model,
		Temperature: *query.Parameters.Temperature,
		MaxTokens:   *query.Parameters.MaxTokens,
	}
	modelInfoMarshalled, err := json.Marshal(modelInfo)
	if err != nil {
		return nil, err
	}
	writer.Write([]byte("dllm:"))
	writer.Write(modelInfoMarshalled)
	writer.Write([]byte("\n"))
	return &Stream{response, writer}, nil
}

func (stream *Stream) Read(callback Callback) {
	reader := bufio.NewReader(stream.response.Body)
	for {
		line, _, err := reader.ReadLine()
		log.Debug().Msgf("Read line: %s", line)
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

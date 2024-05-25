package dllm

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Dummy struct {
}

func NewDummy() (*Dummy, error) {
	return &Dummy{}, nil
}

func (d *Dummy) Name() string {
	return "Dummy"
}

func (d *Dummy) CompletionURL() *url.URL {
	return nil
}

func (d *Dummy) CreateData(query *Query) ([]byte, error) {
	return nil, nil
}

func (d *Dummy) addHeaders(request *http.Request) {
}

func (d *Dummy) do(request *http.Request) (*http.Response, error) {
	return nil, nil
}

func (d *Dummy) GetWriterCallback() func([]byte) ([]byte, bool) {
	return nil
}

// TODO: refactor this to agent
func DummyHandler(writer http.ResponseWriter, req *http.Request) {
	AddHeaders(writer)
	defer req.Body.Close()
	for range 10 {
		writer.Write([]byte("Hello\n"))
		fmt.Println("Sent message")
		if f, ok := writer.(http.Flusher); ok {
			fmt.Println("Flushing")
			f.Flush()
		} else {
			fmt.Println("Unable to convert to flusher")
		}
		// sleep for 500ms
		time.Sleep(500 * time.Millisecond)
	}
	writer.Write([]byte("Goodbye\n"))
	return
}

package dllm

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
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

type streamingResponseWriter struct {
	http.ResponseWriter
	body chan string
	done chan bool
}

func (w *streamingResponseWriter) Write(data []byte) (int, error) {
	w.body <- string(data)
	return len(data), nil
}

func (w *streamingResponseWriter) Close() {
	close(w.done)
}

type MyBody struct {
	strchan chan string
	done    chan bool
}

func (b *MyBody) Read(p []byte) (n int, err error) {
	select {
	case <-b.done:
		log.Debug().Msg("Closing body")
		return 0, io.EOF
	case <-b.strchan:
		log.Debug().Msg("Reading")
		return copy(p, []byte("Foo\n")), nil
	}
}

func (b *MyBody) Close() error {
	return nil
}

func (d *Dummy) do(request *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
	}
	body := &MyBody{
		strchan: make(chan string),
		done:    make(chan bool),
	}
	resp.Body = body
	go func() {
		for i := 0; i < 32; i++ {
			body.strchan <- "Hello\n"
			time.Sleep(90 * time.Millisecond)
			log.Debug().Msg("Sent message")
		}
		body.done <- true
	}()
	log.Debug().Msg("Returning response")
	return resp, nil
}

func (d *Dummy) GetWriterCallback() func(input []byte) ([]byte, bool) {
	return func(input []byte) ([]byte, bool) {
		log.Debug().Msgf("Received chunk: %s", input)
		stop := false
		res := input
		if bytes.Equal(input, []byte("stop")) {
			stop = true
			res = []byte{}
		}
		return res, stop
	}
}

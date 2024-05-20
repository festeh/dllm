package dllm

import (
	"fmt"
	"net/http"
	"time"
)

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

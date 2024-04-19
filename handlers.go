package dllm

import (
	"fmt"
	"net/http"
	"time"
)

func DummyHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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

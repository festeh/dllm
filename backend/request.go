package dllm

import (
	"net/http"
	"net/url"
)

func NewPostRequest(u *url.URL) *http.Request {
	return &http.Request{
		Method: "POST",
		URL:    u,
		Header: make(http.Header),
	}
}

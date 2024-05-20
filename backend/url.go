package dllm

import (
	"net/url"
)

func ParseUrlYolo(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return res
}

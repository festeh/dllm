package dllm

import (
	"bytes"
	"encoding/json"
)

func LoadQuery(b []byte, query interface{}) (err error) {
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&query)
	return
}

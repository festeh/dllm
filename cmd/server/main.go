package main

import "dllm"

func main() {
	server := dllm.Server{Port: 4242}
	server.AddRoute("/dummy", dllm.DummyHandler)
	server.Start()
}

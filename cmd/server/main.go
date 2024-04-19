package main

import "dllm"

func main() {
	server := dllm.Server{Port: 4242}
	server.AddRoute("/dummy", dllm.DummyHandler)
	server.AddRoute("/openai", dllm.InitOpenAIHandler())
	server.Start()
}

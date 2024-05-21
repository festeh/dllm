package main

import "dllm"

func main() {
	server := dllm.Server{Port: 4242}

	manager := &dllm.Manager{}

	openai, err := dllm.NewOpenAI()
	if err == nil {
		server.AddRoute("/openai", manager.CreateHandler(openai))
	}
	anthropic, err := dllm.NewAnthropic()
	if err == nil {
		server.AddRoute("/anthropic", manager.CreateHandler(anthropic))
	}

	// TODO: dummy handler
	// server.AddRoute("/dummy", dllm.DummyHandler)
	server.Start()
}

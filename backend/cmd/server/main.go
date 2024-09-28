package main

import (
	"dllm"
	"flag"
	"fmt"

	"github.com/rs/zerolog"
)

var Version = "development"

func main() {
	vers := flag.Bool("version", false, "Print the version number")
	v := flag.Bool("v", false, "Print the version number")
	flag.Parse()
	if *vers || *v {
		fmt.Println("Version:", Version)
		return
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
	dummy, err := dllm.NewDummy()
	if err == nil {
		server.AddRoute("/dummy", manager.CreateHandler(dummy))
	}
	server.Start()
}

package main

import (
	"dllm"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "development"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

package main

import (
	"dllm"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	agentType := flag.String("agent", "openai", "Agent to use")
	message := flag.String("message", "", "Message to send")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *agentType != "openai" && *agentType != "anthropic" {
		log.Error().Msg("Invalid agent type")
		os.Exit(1)
	}
	var agent dllm.Agent
	var err error
	if *agentType == "openai" {
		agent, err = dllm.NewOpenAI()
	} else {
		agent, err = dllm.NewAnthropic()
	}
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
	systemMessage := dllm.Message{
		Role:    "system",
		Content: "Hi!",
	}
	userMessage := dllm.Message{
		Role:    "user",
		Content: *message,
	}
	query := &dllm.Query{
		Messages: []dllm.Message{systemMessage, userMessage},
	}
	stream, err := agent.GetStream(query, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stream.Close()
	log.Debug().Msg("Response begin")
	stream.Read(agent.GetWriterCallback())
	log.Debug().Msg("End of response")
}

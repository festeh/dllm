package main

import (
	"dllm"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExitIfErr(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func main() {
	agentType := flag.String("agent", "openai", "Agent to use")
	flag.StringVar(agentType, "a", "openai", "Agent to use")
	message := flag.String("message", "", "Message to send")
	flag.StringVar(message, "m", "", "Message to send")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.BoolVar(verbose, "v", false, "Verbose output")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *agentType == "a" {
		*agentType = "anthropic"
	}
	if *agentType == "o" {
		*agentType = "openai"
	}
	if *agentType == "d" {
		*agentType = "dummy"
	}
	if *agentType != "openai" && *agentType != "anthropic" && *agentType != "dummy" {
		ExitIfErr(fmt.Errorf("Invalid agent type: %s", *agentType))
	}
	var agent dllm.Agent
	var err error
	if *agentType == "openai" {
		agent, err = dllm.NewOpenAI()
	} else if *agentType == "anthropic"{
		agent, err = dllm.NewAnthropic()
	} else {
		agent, err = dllm.NewDummy()
	}
	ExitIfErr(err)
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
	stream, err := dllm.NewStream(query, os.Stdout, agent)
	log.Debug().Msgf("Error: %v", err)
	ExitIfErr(err)
	defer stream.Close()
	log.Debug().Msg("Response begin")
	stream.Read(agent.GetWriterCallback())
	log.Debug().Msg("End of response")
}

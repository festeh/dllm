package main

import (
	"dllm"
	"flag"
	"fmt"
	"os"
)

func main() {
	agentType := flag.String("agent", "openai", "Agent to use")
	message := flag.String("message", "", "Message to send")
	flag.Parse()
	if *agentType != "openai" && *agentType != "anthropic" {
		fmt.Println("Invalid agent")
		os.Exit(1)
	}
	var agent any
	var err error
	if *agentType == "openai" {
		agent, err = dllm.NewOpenAI()
	} else {
		agent, err = dllm.NewAnthropic()
	}
	if err != nil {
		fmt.Println(err)
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
	stream, err := agent.(dllm.Agent[any]).GetStream(query, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stream.Close()
	stream.Read(agent.(dllm.Agent[any]).GetWriterCallback())
}

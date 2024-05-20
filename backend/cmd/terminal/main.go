package main

import (
	"fmt"
	"os"
	"dllm"
	"flag"
)


func main() {
	agentType := flag.String("agent", "openai", "Agent to use")
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
	stream, err := agent.(dllm.Agent[any, any]).GetStream([]byte{}, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stream.Close()
	stream.Read(agent.(dllm.Agent[any, any]).GetWriterCallback())
}

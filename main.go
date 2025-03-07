package main

import (
	"fmt"
	"hts_agent/llm_provider"
)

func main() {
	prompt := "Tell me a haiku about Robin Hood."
	model := "phi4-mini:latest"

	response, err := llm_provider.GenerateOllamaResponse(prompt, model)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", response)
}
package main

import (
	"fmt"
	"hts_agent/llm_provider"
)

func main() {
	prompt := "What is 2+2?"
	model := "phi4-mini:latest"

	response, err := llm_provider.GenerateResponse(prompt, model)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", response)
}
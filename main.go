package main

import (
	"fmt"
	"hts_agent/llm_provider"
	// Uncomment to use Grok provider
	// "hts_agent/config"
)

func main() {
	prompt := "Tell me a haiku about Robin Hood."
	
	// Example using Ollama provider
	fmt.Println("Using Ollama provider:")
	ollamaModel := "phi4-mini:latest"
	
	ollamaResponse, err := llm_provider.GenerateOllamaResponse(prompt, ollamaModel)
	if err != nil {
		fmt.Printf("Ollama Error: %v\n", err)
	} else {
		fmt.Printf("Ollama Response: %s\n\n", ollamaResponse)
	}

	// Example using Grok provider
	// To use Grok:
	// 1. Copy .env.example to .env and add your Grok API key
	// 2. Uncomment the code below
	/*
	fmt.Println("Using Grok provider:")
	grokModel := "grok-1"
	
	// Get Grok API key from config (loaded from .env file)
	apiKey := config.GetGrokAPIKey()
	if apiKey == "" {
		fmt.Println("Grok API key not found. Please set it in the .env file.")
		return
	}
	
	llm_provider.SetGrokApiKey(apiKey)
	
	grokResponse, err := llm_provider.GenerateGrokResponse(prompt, grokModel)
	if err != nil {
		fmt.Printf("Grok Error: %v\n", err)
	} else {
		fmt.Printf("Grok Response: %s\n", grokResponse)
	}
	*/
}

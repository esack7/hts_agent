package llm_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GrokRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GrokResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Default Grok API endpoint
var grokApiEndpoint = "https://api.x.ai/v1/chat/completions"

// API key for Grok
var grokApiKey = ""

// SetGrokApiKey sets the API key for Grok
func SetGrokApiKey(apiKey string) {
	grokApiKey = apiKey
}

// GenerateGrokResponse sends a prompt to the Grok API and returns the generated response
func GenerateGrokResponse(prompt, model string) (string, error) {
	if grokApiKey == "" {
		return "", fmt.Errorf("Grok API key not set. Call SetGrokApiKey first")
	}

	// Prepare the request body
	reqBody := GrokRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", grokApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+grokApiKey)
	req.Header.Set("x-api-version", "1.0")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	var grokResp GrokResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Check if we have any choices
	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from Grok API")
	}

	// Return the content of the first choice
	return grokResp.Choices[0].Message.Content, nil
}

package llm_provider

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Default Ollama API endpoint
var apiEndpoint = "http://localhost:11434/api/generate"

// GenerateResponse sends a prompt to the Ollama API and returns the generated response
func GenerateResponse(prompt, model string) (string, error) {
	// Prepare the request body
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the streaming response
	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	
	for scanner.Scan() {
		line := scanner.Text()
		var streamResp OllamaResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return "", fmt.Errorf("error parsing response line: %v", err)
		}
		
		fullResponse.WriteString(streamResp.Response)
		
		if streamResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response stream: %v", err)
	}

	return fullResponse.String(), nil
}

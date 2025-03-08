package llm_provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateGrokResponse(t *testing.T) {
	// Create a mock server to simulate Grok API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", contentType)
		}

		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("Expected Authorization: Bearer test-api-key, got %s", authHeader)
		}

		// Verify API version header
		apiVersionHeader := r.Header.Get("x-api-version")
		if apiVersionHeader != "1.0" {
			t.Errorf("Expected x-api-version: 1.0, got %s", apiVersionHeader)
		}

		// Decode request body to verify it
		var request GrokRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			return
		}

		// Verify request fields
		if request.Model != "grok-1" {
			t.Errorf("Expected model: grok-1, got %s", request.Model)
		}
		if len(request.Messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(request.Messages))
		} else {
			if request.Messages[0].Role != "user" {
				t.Errorf("Expected role: user, got %s", request.Messages[0].Role)
			}
			if request.Messages[0].Content != "Hello, world!" {
				t.Errorf("Expected content: Hello, world!, got %s", request.Messages[0].Content)
			}
		}

		// Write mock response
		mockResponse := GrokResponse{
			ID:      "resp-123456",
			Object:  "chat.completion",
			Created: 1678048124,
			Model:   "grok-1",
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: "Hello! How can I assist you today?",
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 8,
				TotalTokens:      18,
			},
		}

		responseJSON, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := grokApiEndpoint
	grokApiEndpoint = server.URL
	defer func() { grokApiEndpoint = originalEndpoint }()

	// Set test API key
	SetGrokApiKey("test-api-key")

	// Call the function
	response, err := GenerateGrokResponse("Hello, world!", "grok-1")

	// Check for errors
	if err != nil {
		t.Errorf("GenerateGrokResponse returned an error: %v", err)
	}

	// Verify the response
	expected := "Hello! How can I assist you today?"
	if response != expected {
		t.Errorf("Expected response %q, got %q", expected, response)
	}
}

func TestGenerateGrokResponse_NoApiKey(t *testing.T) {
	// Reset API key
	originalApiKey := grokApiKey
	grokApiKey = ""
	defer func() { grokApiKey = originalApiKey }()

	// Call the function
	_, err := GenerateGrokResponse("Hello, world!", "grok-1")

	// Check for errors
	if err == nil {
		t.Error("Expected an error about missing API key, got nil")
	}
}

func TestGenerateGrokResponse_ServerError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := grokApiEndpoint
	grokApiEndpoint = server.URL
	defer func() { grokApiEndpoint = originalEndpoint }()

	// Set test API key
	SetGrokApiKey("test-api-key")

	// Call the function
	_, err := GenerateGrokResponse("Hello, world!", "grok-1")

	// Check for errors
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestGenerateGrokResponse_EmptyResponse(t *testing.T) {
	// Create a mock server that returns an empty response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a valid JSON response but with no choices
		emptyResponse := GrokResponse{
			ID:      "resp-123456",
			Object:  "chat.completion",
			Created: 1678048124,
			Model:   "grok-1",
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 0,
				TotalTokens:      10,
			},
		}

		responseJSON, _ := json.Marshal(emptyResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := grokApiEndpoint
	grokApiEndpoint = server.URL
	defer func() { grokApiEndpoint = originalEndpoint }()

	// Set test API key
	SetGrokApiKey("test-api-key")

	// Call the function
	_, err := GenerateGrokResponse("Hello, world!", "grok-1")

	// Check for errors
	if err == nil {
		t.Error("Expected an error about no choices, got nil")
	}
}

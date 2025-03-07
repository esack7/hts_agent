package llm_provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateResponse(t *testing.T) {
	// Create a mock server to simulate Ollama API
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

		// Decode request body to verify it
		var request OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			return
		}

		// Verify request fields
		if request.Model != "llama2" {
			t.Errorf("Expected model: llama2, got %s", request.Model)
		}
		if request.Prompt != "Hello, world!" {
			t.Errorf("Expected prompt: Hello, world!, got %s", request.Prompt)
		}

		// Write mock responses to simulate streaming
		response1, _ := json.Marshal(OllamaResponse{Response: "Hello", Done: false})
		response2, _ := json.Marshal(OllamaResponse{Response: ", how are you?", Done: true})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response1)
		w.Write([]byte("\n"))
		w.Write(response2)
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := apiEndpoint
	apiEndpoint = server.URL
	defer func() { apiEndpoint = originalEndpoint }()

	// Call the function
	response, err := GenerateResponse("Hello, world!", "llama2")

	// Check for errors
	if err != nil {
		t.Errorf("GenerateResponse returned an error: %v", err)
	}

	// Verify the response
	expected := "Hello, how are you?"
	if response != expected {
		t.Errorf("Expected response %q, got %q", expected, response)
	}
}

func TestGenerateResponse_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := apiEndpoint
	apiEndpoint = server.URL
	defer func() { apiEndpoint = originalEndpoint }()

	// Call the function
	_, err := GenerateResponse("Hello, world!", "llama2")

	// Check for errors
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestGenerateResponse_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid JSON"))
	}))
	defer server.Close()

	// Override API endpoint to use test server
	originalEndpoint := apiEndpoint
	apiEndpoint = server.URL
	defer func() { apiEndpoint = originalEndpoint }()

	// Call the function
	_, err := GenerateResponse("Hello, world!", "llama2")

	// Check for errors
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "error parsing response line") {
		t.Errorf("Expected error about parsing response, got: %v", err)
	}
}

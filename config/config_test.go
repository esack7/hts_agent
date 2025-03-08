package config

import (
	"os"
	"sync"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment variable
	originalGrokAPIKey := os.Getenv("GROK_API_KEY")
	defer os.Setenv("GROK_API_KEY", originalGrokAPIKey)

	// Set test environment variable
	testAPIKey := "test-api-key"
	os.Setenv("GROK_API_KEY", testAPIKey)

	// Reset config singleton for testing
	config = nil
	configOnce = sync.Once{}

	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}

	// Check if config loaded correctly
	if cfg.GrokAPIKey != testAPIKey {
		t.Errorf("LoadConfig() GrokAPIKey = %v, want %v", cfg.GrokAPIKey, testAPIKey)
	}
}

func TestGetGrokAPIKey(t *testing.T) {
	// Save original environment variable
	originalGrokAPIKey := os.Getenv("GROK_API_KEY")
	defer os.Setenv("GROK_API_KEY", originalGrokAPIKey)

	// Set test environment variable
	testAPIKey := "test-api-key"
	os.Setenv("GROK_API_KEY", testAPIKey)

	// Reset config singleton for testing
	config = nil
	configOnce = sync.Once{}

	// Get API key
	apiKey := GetGrokAPIKey()

	// Check if API key is correct
	if apiKey != testAPIKey {
		t.Errorf("GetGrokAPIKey() = %v, want %v", apiKey, testAPIKey)
	}
}

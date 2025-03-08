package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	GrokAPIKey string
}

var (
	config     *Config
	configOnce sync.Once
)

// LoadConfig loads the configuration from environment variables
// It attempts to load from .env file if it exists
func LoadConfig() (*Config, error) {
	configOnce.Do(func() {
		// Try to load .env file from the current directory
		// If it doesn't exist, we'll just use environment variables
		envFile := filepath.Join(".", ".env")
		if _, err := os.Stat(envFile); err == nil {
			err := godotenv.Load(envFile)
			if err != nil {
				fmt.Printf("Warning: Error loading .env file: %v\n", err)
			}
		}

		config = &Config{
			GrokAPIKey: os.Getenv("GROK_API_KEY"),
		}
	})

	return config, nil
}

// GetGrokAPIKey returns the Grok API key
func GetGrokAPIKey() string {
	cfg, err := LoadConfig()
	if err != nil {
		return ""
	}
	return cfg.GrokAPIKey
}

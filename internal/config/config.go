package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	AppEnv   string
	AppPort  string
	LogLevel string
}

// Load reads the .env file (if present) and populates a Config struct.
// Missing optional keys fall back to sensible defaults.
func Load() (*Config, error) {
	// Load .env – ignore error when the file simply doesn't exist (e.g. production).
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:   getEnv("APP_ENV", "production"),
		AppPort:  getEnv("APP_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if cfg.AppPort == "" {
		return nil, fmt.Errorf("APP_PORT must not be empty")
	}

	return cfg, nil
}

// getEnv returns the value for key, or fallback when the key is unset / empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

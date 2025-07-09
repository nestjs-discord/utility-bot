package env

import (
	"errors"
	"os"
	"strings"
)

type GeminiConfig struct {
	APIKeys []string
}

func NewGeminiConfig() (*GeminiConfig, error) {
	apiKeys := os.Getenv("GEMINI_API_KEYS")
	if apiKeys == "" {
		return nil, errors.New("GEMINI_API_KEYS environment variable not set")
	}

	return &GeminiConfig{
		APIKeys: strings.Split(apiKeys, ","),
	}, nil
}

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func NewConfig() (*Config, error) {
	const fileName = "./config.json"

	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config

	err = json.NewDecoder(bytes.NewReader(file)).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	return &cfg, nil
}

// Config is the config specification for the entire app
type Config struct {
	Service struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"service"`

	Rest struct {
		Port int `json:"port"`
	} `json:"rest" `
}

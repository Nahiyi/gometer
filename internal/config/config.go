package config

import (
	"encoding/json"
	"os"
)

// Config represents the request configuration file
type Config struct {
	Request RequestConfig `json:"request"`
	Users   []UserConfig  `json:"users"`
}

// RequestConfig represents shared request configuration
type RequestConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// UserConfig represents per-thread user configuration
type UserConfig struct {
	Headers map[string]string `json:"headers"`
}

// Load reads and parses the JSON config file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

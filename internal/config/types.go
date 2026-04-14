package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Sentinel errors for config validation
var (
	ErrVersionRequired   = errors.New("config version is required")
	ErrFrameworkRequired = errors.New("at least one framework is required")
	ErrInvalidFramework  = errors.New("invalid framework")
)

// Config represents .agent-instruction/config.json
type Config struct {
	Version    string   `json:"version"`
	Packages   []string `json:"packages"`
	Frameworks []string `json:"frameworks"`
}

// Validate checks config for required fields and valid values
func (c *Config) Validate() error {
	if c.Version == "" {
		return ErrVersionRequired
	}

	if len(c.Frameworks) == 0 {
		return ErrFrameworkRequired
	}

	validFrameworks := map[string]bool{
		"claude": true,
		"agents": true,
	}

	for _, fw := range c.Frameworks {
		if !validFrameworks[fw] {
			return fmt.Errorf("%w: %s (must be 'claude' or 'agents')", ErrInvalidFramework, fw)
		}
	}

	return nil
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config JSON: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

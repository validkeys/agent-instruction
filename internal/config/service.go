package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/validkeys/agent-instruction/internal/files"
	"github.com/validkeys/agent-instruction/internal/rules"
)

// ConfigService provides operations for configuration files
type ConfigService interface {
	// LoadConfig reads and validates config.json
	LoadConfig(path string) (*Config, error)

	// SaveConfig validates and writes config.json atomically
	SaveConfig(path string, config *Config) error

	// LoadRuleFile reads and validates a rule file
	LoadRuleFile(path string) (*rules.RuleFile, error)

	// SaveRuleFile validates and writes a rule file atomically
	SaveRuleFile(path string, rule *rules.RuleFile) error
}

// DefaultConfigService implements ConfigService using standard operations
type DefaultConfigService struct{}

// LoadConfig reads and validates config.json
func (s *DefaultConfigService) LoadConfig(path string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found at %s: run 'agent-instruction init' first", path)
		}
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config JSON in %s: %w", path, err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	return &cfg, nil
}

// SaveConfig validates and writes config.json atomically
func (s *DefaultConfigService) SaveConfig(path string, config *Config) error {
	// Validate before saving
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	// Write atomically
	if err := files.WriteAtomic(path, data); err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}

	return nil
}

// LoadRuleFile reads and validates a rule file
func (s *DefaultConfigService) LoadRuleFile(path string) (*rules.RuleFile, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("rule file not found: %s", path)
		}
		return nil, fmt.Errorf("read rule file %s: %w", path, err)
	}

	// Parse JSON
	var rule rules.RuleFile
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("parse rule JSON in %s: %w", path, err)
	}

	// Validate
	if err := rule.Validate(); err != nil {
		return nil, fmt.Errorf("invalid rule in %s: %w", path, err)
	}

	return &rule, nil
}

// SaveRuleFile validates and writes a rule file atomically
func (s *DefaultConfigService) SaveRuleFile(path string, rule *rules.RuleFile) error {
	// Validate before saving
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal rule: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	// Write atomically
	if err := files.WriteAtomic(path, data); err != nil {
		return fmt.Errorf("write rule file %s: %w", path, err)
	}

	return nil
}

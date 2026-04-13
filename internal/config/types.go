package config

import "fmt"

// Config represents .agent-instruction/config.json
type Config struct {
	Version    string   `json:"version"`
	Packages   []string `json:"packages"`
	Frameworks []string `json:"frameworks"`
}

// Validate checks config for required fields and valid values
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(c.Frameworks) == 0 {
		return fmt.Errorf("at least one framework is required")
	}

	validFrameworks := map[string]bool{
		"claude": true,
		"agents": true,
	}

	for _, fw := range c.Frameworks {
		if !validFrameworks[fw] {
			return fmt.Errorf("invalid framework: %s (must be 'claude' or 'agents')", fw)
		}
	}

	return nil
}

package builder

import (
	"fmt"
	"os"

	"github.com/validkeys/agent-instruction/internal/rules"
)

// ComposeInstructions composes final instruction set from global and package configs
func ComposeInstructions(globalRulesPath, packageConfigPath string, ruleSvc rules.RuleService) ([]rules.Instruction, error) {
	if ruleSvc == nil {
		return nil, fmt.Errorf("rule service cannot be nil")
	}

	allInstructions := make([]rules.Instruction, 0)

	// Resolve global rules if file exists
	if fileExists(globalRulesPath) {
		globalInstructions, err := ruleSvc.ResolveRules(globalRulesPath)
		if err != nil {
			return nil, fmt.Errorf("resolve global rules from %s: %w", globalRulesPath, err)
		}
		allInstructions = append(allInstructions, globalInstructions...)
	}

	// Resolve package config if file exists
	if fileExists(packageConfigPath) {
		packageInstructions, err := ruleSvc.ResolveRules(packageConfigPath)
		if err != nil {
			return nil, fmt.Errorf("resolve package config from %s: %w", packageConfigPath, err)
		}
		allInstructions = append(allInstructions, packageInstructions...)
	}

	return allInstructions, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

package rules

import (
	"fmt"
)

// RuleConfigService defines operations for loading and saving rule files
type RuleConfigService interface {
	LoadRuleFile(path string) (*RuleFile, error)
	SaveRuleFile(path string, rule *RuleFile) error
}

// RuleService provides high-level API for rule management
type RuleService interface {
	// ResolveRules resolves all imports starting from rootPath
	ResolveRules(rootPath string) ([]Instruction, error)

	// LoadRuleFile loads and validates a rule file
	LoadRuleFile(path string) (*RuleFile, error)

	// SaveRuleFile validates and saves a rule file
	SaveRuleFile(path string, rule *RuleFile) error

	// AddInstruction adds an instruction to a rule file
	AddInstruction(rulePath string, instruction Instruction) error
}

// DefaultRuleService implements RuleService using resolver and config service
type DefaultRuleService struct {
	configService RuleConfigService
}

// NewRuleService creates a new rule service instance
func NewRuleService(configService RuleConfigService) *DefaultRuleService {
	return &DefaultRuleService{
		configService: configService,
	}
}

// ResolveRules resolves all imports starting from rootPath
func (s *DefaultRuleService) ResolveRules(rootPath string) ([]Instruction, error) {
	resolver := NewResolver(s.configService)
	return resolver.ResolveImports(rootPath)
}

// LoadRuleFile loads and validates a rule file
func (s *DefaultRuleService) LoadRuleFile(path string) (*RuleFile, error) {
	return s.configService.LoadRuleFile(path)
}

// SaveRuleFile validates and saves a rule file
func (s *DefaultRuleService) SaveRuleFile(path string, rule *RuleFile) error {
	return s.configService.SaveRuleFile(path, rule)
}

// AddInstruction adds an instruction to a rule file atomically
func (s *DefaultRuleService) AddInstruction(rulePath string, instruction Instruction) error {
	// Load existing rule
	rule, err := s.LoadRuleFile(rulePath)
	if err != nil {
		return fmt.Errorf("load rule file: %w", err)
	}

	// Append instruction
	rule.Instructions = append(rule.Instructions, instruction)

	// Save atomically
	if err := s.SaveRuleFile(rulePath, rule); err != nil {
		return fmt.Errorf("save rule file: %w", err)
	}

	return nil
}

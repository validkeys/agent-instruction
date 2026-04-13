package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// FileConfigService implements RuleConfigService using filesystem
type FileConfigService struct {
	baseDir string
}

// NewFileConfigService creates a new file-based config service
func NewFileConfigService(baseDir string) *FileConfigService {
	return &FileConfigService{baseDir: baseDir}
}

// LoadRuleFile loads a rule file from the filesystem
func (s *FileConfigService) LoadRuleFile(path string) (*RuleFile, error) {
	// If path is not absolute, make it relative to baseDir
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.baseDir, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var rule RuleFile
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("parse JSON from %s: %w", path, err)
	}

	// Validate rule
	if err := rule.Validate(); err != nil {
		return nil, fmt.Errorf("invalid rule file %s: %w", path, err)
	}

	return &rule, nil
}

// SaveRuleFile saves a rule file to the filesystem
func (s *FileConfigService) SaveRuleFile(path string, rule *RuleFile) error {
	// Validate before saving
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	// If path is not absolute, make it relative to baseDir
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.baseDir, path)
	}

	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}

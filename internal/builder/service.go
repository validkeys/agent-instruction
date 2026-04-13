package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/validkeys/agent-instruction/internal/files"
	"github.com/validkeys/agent-instruction/internal/rules"
)

// BuildService provides high-level API for building instruction files
type BuildService interface {
	// BuildFile resolves imports, generates markdown, and writes output atomically
	BuildFile(rulePath, outputPath string) error

	// GenerateForPackage generates instruction files for a package directory
	GenerateForPackage(packagePath string) error
}

// DefaultBuildService implements BuildService using rule and file services
type DefaultBuildService struct {
	ruleService rules.RuleService
	fileService files.FileService
}

// NewBuildService creates a new build service instance
func NewBuildService(ruleService rules.RuleService, fileService files.FileService) *DefaultBuildService {
	return &DefaultBuildService{
		ruleService: ruleService,
		fileService: fileService,
	}
}

// BuildFile orchestrates the complete workflow:
//  1. Resolve imports from rule file
//  2. Generate markdown from instructions
//  3. Read existing file if present
//  4. Parse managed sections
//  5. Merge generated with existing
//  6. Write atomically with backup
func (s *DefaultBuildService) BuildFile(rulePath, outputPath string) error {
	// Validate inputs
	if rulePath == "" {
		return fmt.Errorf("rule path cannot be empty")
	}
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Resolve imports to get all instructions
	instructions, err := s.ruleService.ResolveRules(rulePath)
	if err != nil {
		return fmt.Errorf("resolve rules from %s: %w", rulePath, err)
	}

	// Generate markdown from instructions
	generatedContent := InstructionsToMarkdown(instructions)

	// Read existing file if present
	var existingManaged *files.ManagedContent
	existingBytes, err := s.fileService.ReadFile(outputPath)
	if err == nil {
		// File exists - parse managed sections
		existingManaged, err = s.fileService.ParseManaged(existingBytes)
		if err != nil {
			return fmt.Errorf("parse managed sections in %s: %w", outputPath, err)
		}

		// Create backup before modifying existing file
		if err := s.fileService.BackupFile(outputPath); err != nil {
			return fmt.Errorf("backup %s: %w", outputPath, err)
		}
	}
	// If file doesn't exist (err != nil), existingManaged stays nil

	// Merge generated content with existing user content
	finalContent := BuildManagedFile(generatedContent, existingManaged)

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory %s: %w", outputDir, err)
	}

	// Write atomically
	if err := s.fileService.WriteFile(outputPath, []byte(finalContent)); err != nil {
		return fmt.Errorf("write output file %s: %w", outputPath, err)
	}

	return nil
}

// GenerateForPackage generates instruction files for a package directory
// This is a placeholder for M4 implementation
func (s *DefaultBuildService) GenerateForPackage(packagePath string) error {
	return fmt.Errorf("GenerateForPackage not yet implemented")
}

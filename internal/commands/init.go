package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/validkeys/agent-instruction/internal/config"
	"github.com/validkeys/agent-instruction/internal/rules"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize agent-instruction in repository",
		Long: `Creates .agent-instruction directory structure and configuration files.

This command sets up the initial directory structure for managing AI agent
instructions in your repository. It creates config.json, directories for rules,
and optional backup of existing files.`,
		Example: `  # Initialize with interactive prompts
  agent-instruction init

  # Initialize with defaults (non-interactive)
  agent-instruction init --non-interactive

  # Initialize with specific frameworks
  agent-instruction init --non-interactive --frameworks claude

  # Initialize with specific packages
  agent-instruction init --non-interactive --packages app,lib`,
		Args: cobra.NoArgs,
		RunE: runInit,
	}

	cmd.Flags().Bool("non-interactive", false, "skip interactive prompts and use defaults")
	cmd.Flags().String("frameworks", "", "comma-separated frameworks (claude,agents)")
	cmd.Flags().String("packages", "", "comma-separated package paths or 'auto'")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	// Get current directory
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check if already initialized
	agentDir := filepath.Join(baseDir, ".agent-instruction")
	if _, err := os.Stat(agentDir); err == nil {
		return fmt.Errorf("already initialized: .agent-instruction directory exists\nUse 'agent-instruction build' to regenerate files")
	}

	// Get flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	frameworksFlag, _ := cmd.Flags().GetString("frameworks")
	packagesFlag, _ := cmd.Flags().GetString("packages")

	// Find existing instruction files
	existingFiles := findExistingInstructionFiles(baseDir)

	// Handle backups if files exist
	if len(existingFiles) > 0 {
		if nonInteractive {
			// In non-interactive mode, always create backups
			for _, file := range existingFiles {
				if err := createBackup(filepath.Join(baseDir, file)); err != nil {
					return fmt.Errorf("create backup of %s: %w", file, err)
				}
			}
		} else {
			// Interactive mode: prompt user
			fmt.Fprintf(cmd.OutOrStdout(), "Found existing instruction files:\n")
			for _, file := range existingFiles {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", file)
			}
			if promptYesNo(cmd, "Create backups?", true) {
				for _, file := range existingFiles {
					if err := createBackup(filepath.Join(baseDir, file)); err != nil {
						return fmt.Errorf("create backup of %s: %w", file, err)
					}
				}
			}
		}
	}

	// Determine frameworks
	var frameworks []string
	if frameworksFlag != "" {
		frameworks = parseCommaSeparated(frameworksFlag)
		// Validate frameworks
		validFrameworks := map[string]bool{"claude": true, "agents": true}
		for _, fw := range frameworks {
			if !validFrameworks[fw] {
				return fmt.Errorf("invalid framework: %s (must be 'claude' or 'agents')", fw)
			}
		}
	} else if nonInteractive {
		frameworks = []string{"claude", "agents"}
	} else {
		frameworks = promptFrameworks(cmd)
	}

	// Determine packages
	var packages []string
	if packagesFlag != "" {
		packages = parseCommaSeparated(packagesFlag)
	} else if nonInteractive {
		packages = []string{"auto"}
	} else {
		packages = promptPackages(cmd)
	}

	// Create directory structure
	if err := os.MkdirAll(filepath.Join(agentDir, "rules"), 0755); err != nil {
		return fmt.Errorf("create directory structure: %w", err)
	}

	// Create config.json
	cfg := createDefaultConfig(frameworks, packages)
	if err := saveConfig(filepath.Join(agentDir, "config.json"), &cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Create global.json template
	globalRule := createGlobalRuleTemplate()
	if err := saveRuleFile(filepath.Join(agentDir, "rules", "global.json"), globalRule); err != nil {
		return fmt.Errorf("save global rule: %w", err)
	}

	// Display success message
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Initialized agent-instruction in %s\n", baseDir)
	fmt.Fprintf(cmd.OutOrStdout(), "\nNext steps:\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  1. Edit .agent-instruction/rules/global.json to add rules\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  2. Run 'agent-instruction build' to generate instruction files\n")

	return nil
}

// createDefaultConfig creates a config with specified values
func createDefaultConfig(frameworks, packages []string) config.Config {
	return config.Config{
		Version:    "1.0",
		Frameworks: frameworks,
		Packages:   packages,
	}
}

// createGlobalRuleTemplate creates a template global rule file
func createGlobalRuleTemplate() *rules.RuleFile {
	return &rules.RuleFile{
		Title: "Global Instructions",
		Instructions: []rules.Instruction{
			{
				Heading: "Project Overview",
				Rule:    "This is a template rule. Replace with your project-specific instructions.",
			},
		},
	}
}

// findExistingInstructionFiles scans for CLAUDE.md and AGENTS.md in root directory
func findExistingInstructionFiles(baseDir string) []string {
	var found []string

	claudePath := filepath.Join(baseDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		found = append(found, "CLAUDE.md")
	}

	agentsPath := filepath.Join(baseDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); err == nil {
		found = append(found, "AGENTS.md")
	}

	return found
}

// createBackup creates a .backup copy of the file
func createBackup(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to backup
		}
		return fmt.Errorf("read file: %w", err)
	}

	backupPath := path + ".backup"
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("write backup: %w", err)
	}

	return nil
}

// saveConfig writes config to disk with proper formatting
func saveConfig(path string, cfg *config.Config) error {
	// Validate before saving
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Add newline at end
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// saveRuleFile writes rule file to disk with proper formatting
func saveRuleFile(path string, rule *rules.RuleFile) error {
	// Validate before saving
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal rule: %w", err)
	}

	// Add newline at end
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write rule file: %w", err)
	}

	return nil
}

// parseCommaSeparated splits comma-separated string into slice
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

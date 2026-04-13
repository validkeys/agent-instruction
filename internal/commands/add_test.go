package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/rules"
)

func TestAddCommand(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temporary test environment for tests that need it
	tempDir := t.TempDir()
	rulesDir := filepath.Join(tempDir, ".agent-instruction", "rules")

	tests := []struct {
		name        string
		setup       func()
		args        []string
		flags       map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid: rule content with title and rule flags",
			setup: func() {
				os.Chdir(tempDir)
				os.MkdirAll(rulesDir, 0755)
				createFile(t, filepath.Join(rulesDir, "global.json"), `{
  "title": "Global Rules",
  "instructions": [
    {
      "rule": "Initial rule"
    }
  ]
}`)
			},
			args:    []string{"Use explicit error handling"},
			flags:   map[string]string{"title": "Error Handling", "rule": "global"},
			wantErr: false,
		},
		{
			name: "requires initialized directory when rule flag omitted",
			setup: func() {
				os.Chdir(originalDir) // Go back to directory without .agent-instruction
			},
			args:        []string{"Always validate input"},
			flags:       map[string]string{},
			wantErr:     true,
			errContains: "not initialized",
		},
		{
			name: "invalid: no arguments",
			setup: func() {
				// No setup needed
			},
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
		{
			name: "invalid: too many arguments",
			setup: func() {
				// No setup needed
			},
			args:        []string{"first", "second"},
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
		{
			name: "invalid: empty rule content",
			setup: func() {
				// No setup needed
			},
			args:        []string{""},
			flags:       map[string]string{},
			wantErr:     true,
			errContains: "rule content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup()
			}

			cmd := newAddCmd()

			// Set flags
			for key, val := range tt.flags {
				if err := cmd.Flags().Set(key, val); err != nil {
					t.Fatalf("failed to set flag %s: %v", key, err)
				}
			}

			// Capture output
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAddCommandWithInteractiveSelection(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temporary test environment
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Set up .agent-instruction structure
	rulesDir := filepath.Join(tempDir, ".agent-instruction", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("failed to create rules directory: %v", err)
	}

	// Create test rule files
	createFile(t, filepath.Join(rulesDir, "global.json"), `{
  "title": "Global Rules",
  "instructions": [
    {
      "rule": "Initial rule"
    }
  ]
}`)
	createFile(t, filepath.Join(rulesDir, "testing.json"), `{
  "title": "Testing Rules",
  "instructions": [
    {
      "rule": "Initial rule"
    }
  ]
}`)

	cmd := newAddCmd()
	cmd.SetArgs([]string{"Test instruction"})

	// Provide interactive input (select option 1: global)
	var stdin bytes.Buffer
	stdin.WriteString("1\n")
	cmd.SetIn(&stdin)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute command
	err = cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Available rule files") {
		t.Errorf("expected interactive prompt, got: %s", output)
	}
	if !strings.Contains(output, "Added instruction to global.json") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify instruction was added to file
	ruleFilePath := filepath.Join(rulesDir, "global.json")
	data, err := os.ReadFile(ruleFilePath)
	if err != nil {
		t.Fatalf("failed to read rule file: %v", err)
	}

	var ruleFile rules.RuleFile
	if err := json.Unmarshal(data, &ruleFile); err != nil {
		t.Fatalf("failed to parse rule file: %v", err)
	}

	if len(ruleFile.Instructions) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(ruleFile.Instructions))
	}

	if ruleFile.Instructions[1].Rule != "Test instruction" {
		t.Errorf("instruction rule = %q, want 'Test instruction'", ruleFile.Instructions[1].Rule)
	}
}

func TestAddIntegration(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temporary test environment
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Set up .agent-instruction structure
	rulesDir := filepath.Join(tempDir, ".agent-instruction", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("failed to create rules directory: %v", err)
	}

	// Create test rule file with existing instruction
	createFile(t, filepath.Join(rulesDir, "testing.json"), `{
  "title": "Testing Rules",
  "instructions": [
    {
      "rule": "Existing test rule"
    }
  ]
}`)

	tests := []struct {
		name       string
		args       []string
		flags      map[string]string
		wantErr    bool
		verifyFunc func(t *testing.T, rulesDir string)
	}{
		{
			name:  "adds instruction without title",
			args:  []string{"Always test edge cases"},
			flags: map[string]string{"rule": "testing"},
			verifyFunc: func(t *testing.T, rulesDir string) {
				data, err := os.ReadFile(filepath.Join(rulesDir, "testing.json"))
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				var rule rules.RuleFile
				if err := json.Unmarshal(data, &rule); err != nil {
					t.Fatalf("failed to parse file: %v", err)
				}

				if len(rule.Instructions) != 2 {
					t.Fatalf("expected 2 instructions, got %d", len(rule.Instructions))
				}

				if rule.Instructions[1].Rule != "Always test edge cases" {
					t.Errorf("instruction = %q, want 'Always test edge cases'", rule.Instructions[1].Rule)
				}

				if rule.Instructions[1].Heading != "" {
					t.Errorf("heading = %q, want empty", rule.Instructions[1].Heading)
				}
			},
		},
		{
			name:  "adds instruction with title",
			args:  []string{"Use table-driven tests"},
			flags: map[string]string{"rule": "testing", "title": "Best Practices"},
			verifyFunc: func(t *testing.T, rulesDir string) {
				data, err := os.ReadFile(filepath.Join(rulesDir, "testing.json"))
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				var rule rules.RuleFile
				if err := json.Unmarshal(data, &rule); err != nil {
					t.Fatalf("failed to parse file: %v", err)
				}

				// Find the last instruction
				lastIdx := len(rule.Instructions) - 1
				if rule.Instructions[lastIdx].Rule != "Use table-driven tests" {
					t.Errorf("instruction = %q, want 'Use table-driven tests'", rule.Instructions[lastIdx].Rule)
				}

				if rule.Instructions[lastIdx].Heading != "Best Practices" {
					t.Errorf("heading = %q, want 'Best Practices'", rule.Instructions[lastIdx].Heading)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newAddCmd()
			cmd.SetArgs(tt.args)

			// Set flags
			for key, val := range tt.flags {
				if err := cmd.Flags().Set(key, val); err != nil {
					t.Fatalf("failed to set flag %s: %v", key, err)
				}
			}

			// Capture output
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			// Execute command
			err := cmd.Execute()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify success message
			output := stdout.String()
			if !strings.Contains(output, "Added instruction") {
				t.Errorf("expected success message, got: %s", output)
			}

			// Run custom verification
			if tt.verifyFunc != nil {
				tt.verifyFunc(t, rulesDir)
			}
		})
	}
}

func TestAddFlagParsing(t *testing.T) {
	tests := []struct {
		name      string
		flags     map[string]string
		wantTitle string
		wantRule  string
	}{
		{
			name:      "both flags provided",
			flags:     map[string]string{"title": "Test Rule", "rule": "testing"},
			wantTitle: "Test Rule",
			wantRule:  "testing",
		},
		{
			name:      "only title provided",
			flags:     map[string]string{"title": "Another Rule"},
			wantTitle: "Another Rule",
			wantRule:  "",
		},
		{
			name:      "only rule provided",
			flags:     map[string]string{"rule": "global"},
			wantTitle: "",
			wantRule:  "global",
		},
		{
			name:      "no flags provided",
			flags:     map[string]string{},
			wantTitle: "",
			wantRule:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newAddCmd()

			// Set flags
			for key, val := range tt.flags {
				if err := cmd.Flags().Set(key, val); err != nil {
					t.Fatalf("failed to set flag %s: %v", key, err)
				}
			}

			// Get flag values
			title, _ := cmd.Flags().GetString("title")
			rule, _ := cmd.Flags().GetString("rule")

			if title != tt.wantTitle {
				t.Errorf("title = %q, want %q", title, tt.wantTitle)
			}

			if rule != tt.wantRule {
				t.Errorf("rule = %q, want %q", rule, tt.wantRule)
			}
		})
	}
}

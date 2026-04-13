package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid: rule content with title and rule flags",
			args:    []string{"Use explicit error handling"},
			flags:   map[string]string{"title": "Error Handling", "rule": "global"},
			wantErr: false,
		},
		{
			name:        "requires initialized directory when rule flag omitted",
			args:        []string{"Always validate input"},
			flags:       map[string]string{},
			wantErr:     true,
			errContains: "not initialized",
		},
		{
			name:        "invalid: no arguments",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
		{
			name:        "invalid: too many arguments",
			args:        []string{"first", "second"},
			wantErr:     true,
			errContains: "accepts 1 arg(s)",
		},
		{
			name:        "invalid: empty rule content",
			args:        []string{""},
			flags:       map[string]string{},
			wantErr:     true,
			errContains: "rule content cannot be empty",
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
	createFile(t, filepath.Join(rulesDir, "global.json"), `{"instructions":[]}`)
	createFile(t, filepath.Join(rulesDir, "testing.json"), `{"instructions":[]}`)

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
	if !strings.Contains(output, "global.json") {
		t.Errorf("expected target to be global.json, got: %s", output)
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

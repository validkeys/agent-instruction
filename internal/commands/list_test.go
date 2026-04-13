package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
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
      "heading": "Error Handling",
      "rule": "Use explicit error handling"
    },
    {
      "rule": "No silent failures"
    }
  ]
}`)

	createFile(t, filepath.Join(rulesDir, "testing.json"), `{
  "title": "Testing Rules",
  "instructions": [
    {
      "heading": "Best Practices",
      "rule": "Use table-driven tests"
    }
  ]
}`)

	tests := []struct {
		name        string
		flags       map[string]bool
		wantContain []string
		wantErr     bool
	}{
		{
			name:  "default mode shows file names and instruction count",
			flags: map[string]bool{},
			wantContain: []string{
				"global.json",
				"testing.json",
				"2 instruction(s)",
				"1 instruction(s)",
			},
		},
		{
			name:  "verbose mode shows full instructions",
			flags: map[string]bool{"verbose": true},
			wantContain: []string{
				"global.json",
				"Global Rules",
				"Error Handling",
				"Use explicit error handling",
				"No silent failures",
				"testing.json",
				"Testing Rules",
				"Best Practices",
				"Use table-driven tests",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newListCmd()

			// Set flags
			for key, val := range tt.flags {
				if err := cmd.Flags().Set(key, "true"); val && err != nil {
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

			output := stdout.String()
			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestListCommandNotInitialized(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temp dir without .agent-instruction
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	cmd := newListCmd()

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute command
	err = cmd.Execute()

	if err == nil {
		t.Error("expected error for uninitialized directory, got nil")
		return
	}

	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected 'not initialized' error, got: %v", err)
	}
}

func TestListCommandEmptyRules(t *testing.T) {
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

	// Set up .agent-instruction structure with no rule files
	rulesDir := filepath.Join(tempDir, ".agent-instruction", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("failed to create rules directory: %v", err)
	}

	cmd := newListCmd()

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute command
	err = cmd.Execute()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	output := stdout.String()
	if !strings.Contains(output, "No rule files found") {
		t.Errorf("expected 'No rule files found' message, got: %s", output)
	}
}

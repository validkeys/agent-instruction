package commands

import (
	"bytes"
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
			name:    "valid: rule content without optional flags",
			args:    []string{"Always validate input"},
			flags:   map[string]string{},
			wantErr: false, // Will prompt interactively
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

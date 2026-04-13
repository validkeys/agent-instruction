package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T, output string)
	}{
		"creates root command with proper use": {
			validate: func(t *testing.T, output string) {
				cmd := NewRootCmd()
				if cmd.Use != "agent-instruction" {
					t.Errorf("Use: got %q, want %q", cmd.Use, "agent-instruction")
				}
			},
		},
		"has short description": {
			validate: func(t *testing.T, output string) {
				cmd := NewRootCmd()
				if cmd.Short == "" {
					t.Error("Short description is empty")
				}
			},
		},
		"has long description": {
			validate: func(t *testing.T, output string) {
				cmd := NewRootCmd()
				if cmd.Long == "" {
					t.Error("Long description is empty")
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t, "")
		})
	}
}

func TestRootCmdVersion(t *testing.T) {
	tests := map[string]struct {
		args       []string
		wantOutput string
	}{
		"version flag displays version": {
			args:       []string{"--version"},
			wantOutput: "agent-instruction version",
		},
		"version flag short form": {
			args:       []string{"-v"},
			wantOutput: "agent-instruction version",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := NewRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tc.wantOutput) {
				t.Errorf("output %q does not contain %q", output, tc.wantOutput)
			}
		})
	}
}

func TestRootCmdHelp(t *testing.T) {
	tests := map[string]struct {
		args       []string
		wantOutput []string
	}{
		"help flag shows usage": {
			args: []string{"--help"},
			wantOutput: []string{
				"agent-instruction",
				"Usage:",
				"Flags:",
			},
		},
		"help flag short form": {
			args: []string{"-h"},
			wantOutput: []string{
				"agent-instruction",
				"Usage:",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := NewRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()
			for _, want := range tc.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := map[string]struct {
		wantErr bool
	}{
		"executes without error": {
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Execute()

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

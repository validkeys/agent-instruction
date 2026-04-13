package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCommand(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		args     []string
		wantErr  bool
		validate func(t *testing.T, dir string, output string)
	}{
		"builds all packages successfully": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude"})
				createTestPackage(t, dir, "packages/api")
				createTestPackage(t, dir, "packages/web")
				return dir
			},
			args:    []string{},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				// Check that CLAUDE.md files were created
				checkFileExists(t, filepath.Join(dir, "packages/api/CLAUDE.md"))
				checkFileExists(t, filepath.Join(dir, "packages/web/CLAUDE.md"))

				// Check output message
				if !strings.Contains(output, "2 package") && !strings.Contains(output, "Successfully") {
					t.Errorf("expected success message in output: %s", output)
				}
			},
		},
		"handles single package": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude"})
				createTestPackage(t, dir, "api")
				return dir
			},
			args:    []string{},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				checkFileExists(t, filepath.Join(dir, "api/CLAUDE.md"))
			},
		},
		"returns error when not initialized": {
			setup: func(t *testing.T) string {
				return t.TempDir() // Empty directory
			},
			args:    []string{},
			wantErr: true,
		},
		"dry-run mode doesn't create files": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude"})
				createTestPackage(t, dir, "api")
				return dir
			},
			args:    []string{"--dry-run"},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				// File should not exist in dry-run mode
				path := filepath.Join(dir, "api/CLAUDE.md")
				if _, err := os.Stat(path); err == nil {
					t.Errorf("file should not exist in dry-run mode: %s", path)
				}

				// Should show what would be done
				if !strings.Contains(output, "DRY RUN") && !strings.Contains(output, "Would") {
					t.Errorf("expected dry-run indication in output: %s", output)
				}
			},
		},
		"verbose mode shows progress": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude"})
				createTestPackage(t, dir, "api")
				createTestPackage(t, dir, "web")
				return dir
			},
			args:    []string{"--verbose"},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				// Should show package names in verbose mode
				if !strings.Contains(output, "api") || !strings.Contains(output, "web") {
					t.Errorf("expected package names in verbose output: %s", output)
				}
			},
		},
		"handles multiple frameworks": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude", "agents"})
				createTestPackage(t, dir, "api")
				return dir
			},
			args:    []string{},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				checkFileExists(t, filepath.Join(dir, "api/CLAUDE.md"))
				checkFileExists(t, filepath.Join(dir, "api/AGENTS.md"))
			},
		},
		"handles empty monorepo": {
			setup: func(t *testing.T) string {
				return setupBuildTestRepo(t, []string{"claude"}) // No packages
			},
			args:    []string{},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				if !strings.Contains(output, "No packages") {
					t.Errorf("expected no packages message: %s", output)
				}
			},
		},
		"preserves existing user content": {
			setup: func(t *testing.T) string {
				dir := setupBuildTestRepo(t, []string{"claude"})
				createTestPackage(t, dir, "api")

				// Create existing file with user content
				existingFile := filepath.Join(dir, "api/CLAUDE.md")
				userContent := `# My Custom Instructions

This is my custom content that should be preserved.

<!-- BEGIN AGENT-INSTRUCTION MANAGED SECTION -->
old generated content
<!-- END AGENT-INSTRUCTION MANAGED SECTION -->

More custom content after.
`
				if err := os.WriteFile(existingFile, []byte(userContent), 0644); err != nil {
					t.Fatalf("write existing file: %v", err)
				}

				return dir
			},
			args:    []string{},
			wantErr: false,
			validate: func(t *testing.T, dir string, output string) {
				content, err := os.ReadFile(filepath.Join(dir, "api/CLAUDE.md"))
				if err != nil {
					t.Fatalf("read file: %v", err)
				}

				str := string(content)
				if !strings.Contains(str, "My Custom Instructions") {
					t.Error("user content before managed section was not preserved")
				}
				if !strings.Contains(str, "More custom content after") {
					t.Error("user content after managed section was not preserved")
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := tc.setup(t)

			// Change to test directory
			oldWd, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				t.Fatalf("chdir: %v", err)
			}
			defer os.Chdir(oldWd)

			// Create command and capture output
			cmd := newBuildCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			output := buf.String()

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v\nOutput: %s", err, output)
			}

			if tc.validate != nil {
				tc.validate(t, dir, output)
			}
		})
	}
}

// setupBuildTestRepo creates a test repo with specified frameworks
func setupBuildTestRepo(t *testing.T, frameworks []string) string {
	t.Helper()

	dir := t.TempDir()
	agentDir := filepath.Join(dir, ".agent-instruction")
	rulesDir := filepath.Join(agentDir, "rules")

	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("create agent dir: %v", err)
	}

	// Create config.json
	configPath := filepath.Join(agentDir, "config.json")
	configContent := `{
  "version": "1.0",
  "frameworks": ["` + strings.Join(frameworks, `","`) + `"],
  "packages": ["auto"]
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Create global rules file
	globalRulesPath := filepath.Join(rulesDir, "global.json")
	globalContent := `{
  "title": "Global Rules",
  "instructions": [
    {
      "heading": "Code Style",
      "rule": "Always use explicit error wrapping with fmt.Errorf and %w"
    }
  ]
}`
	if err := os.WriteFile(globalRulesPath, []byte(globalContent), 0644); err != nil {
		t.Fatalf("write global rules: %v", err)
	}

	return dir
}

// createTestPackage creates a test package directory with agent-instruction.json
func createTestPackage(t *testing.T, baseDir, pkgPath string) {
	t.Helper()

	fullPath := filepath.Join(baseDir, pkgPath)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		t.Fatalf("create package dir: %v", err)
	}

	configPath := filepath.Join(fullPath, "agent-instruction.json")
	configContent := `{
  "title": "` + pkgPath + ` Rules",
  "instructions": [
    {
      "heading": "Package Rules",
      "rule": "Package-specific instructions for ` + pkgPath + `"
    }
  ]
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write package config: %v", err)
	}
}

// checkFileExists fails if file doesn't exist
func checkFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

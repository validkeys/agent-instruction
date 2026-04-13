package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/config"
	"github.com/validkeys/agent-instruction/internal/rules"
)

func TestSetupTestRepo(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T, dir string)
	}{
		"creates temporary directory with .agent-instruction": {
			validate: func(t *testing.T, dir string) {
				agentDir := filepath.Join(dir, ".agent-instruction")
				if _, err := os.Stat(agentDir); os.IsNotExist(err) {
					t.Error("expected .agent-instruction directory to exist")
				}
			},
		},
		"creates rules subdirectory": {
			validate: func(t *testing.T, dir string) {
				rulesDir := filepath.Join(dir, ".agent-instruction", "rules")
				if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
					t.Error("expected rules directory to exist")
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := setupTestRepo(t)
			tc.validate(t, dir)
		})
	}
}

func TestCreateConfig(t *testing.T) {
	tests := map[string]struct {
		config   config.Config
		validate func(t *testing.T, dir string)
	}{
		"writes valid config to disk": {
			config: config.Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{},
			},
			validate: func(t *testing.T, dir string) {
				path := filepath.Join(dir, ".agent-instruction", "config.json")
				assertFileExists(t, path)

				loaded := loadConfig(t, dir)
				if loaded.Version != "1.0" {
					t.Errorf("version: got %q, want %q", loaded.Version, "1.0")
				}
			},
		},
		"writes config with packages": {
			config: config.Config{
				Version:    "1.0",
				Frameworks: []string{"claude", "agents"},
				Packages:   []string{"api", "web"},
			},
			validate: func(t *testing.T, dir string) {
				loaded := loadConfig(t, dir)
				if len(loaded.Packages) != 2 {
					t.Errorf("packages count: got %d, want %d", len(loaded.Packages), 2)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := setupTestRepo(t)
			createConfig(t, dir, tc.config)
			tc.validate(t, dir)
		})
	}
}

func TestCreateRuleFile(t *testing.T) {
	tests := map[string]struct {
		filename string
		ruleFile *rules.RuleFile
		validate func(t *testing.T, dir string)
	}{
		"writes rule file to disk": {
			filename: "test.json",
			ruleFile: &rules.RuleFile{
				Title: "Test Rules",
				Instructions: []rules.Instruction{
					{Rule: "Test rule"},
				},
			},
			validate: func(t *testing.T, dir string) {
				path := filepath.Join(dir, ".agent-instruction", "rules", "test.json")
				assertFileExists(t, path)

				loaded := loadRuleFile(t, dir, "test.json")
				if loaded.Title != "Test Rules" {
					t.Errorf("title: got %q, want %q", loaded.Title, "Test Rules")
				}
			},
		},
		"writes rule file with multiple instructions": {
			filename: "multi.json",
			ruleFile: &rules.RuleFile{
				Title: "Multiple Rules",
				Instructions: []rules.Instruction{
					{Rule: "Rule 1"},
					{Rule: "Rule 2"},
					{Rule: "Rule 3"},
				},
			},
			validate: func(t *testing.T, dir string) {
				loaded := loadRuleFile(t, dir, "multi.json")
				if len(loaded.Instructions) != 3 {
					t.Errorf("instructions count: got %d, want %d", len(loaded.Instructions), 3)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := setupTestRepo(t)
			createRuleFile(t, dir, tc.filename, tc.ruleFile)
			tc.validate(t, dir)
		})
	}
}

func TestFileAssertions(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T, dir string) string
		wantFail bool
	}{
		"assertFileExists succeeds when file exists": {
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "test.txt")
				os.WriteFile(path, []byte("test"), 0644)
				return path
			},
			wantFail: false,
		},
		"assertFileNotExists succeeds when file does not exist": {
			setup: func(t *testing.T, dir string) string {
				return filepath.Join(dir, "nonexistent.txt")
			},
			wantFail: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := tc.setup(t, dir)

			if strings.Contains(name, "assertFileExists") {
				assertFileExists(t, path)
			} else {
				assertFileNotExists(t, path)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := map[string]struct {
		content  string
		validate func(t *testing.T, content string)
	}{
		"reads file content": {
			content: "test content",
			validate: func(t *testing.T, content string) {
				if content != "test content" {
					t.Errorf("content: got %q, want %q", content, "test content")
				}
			},
		},
		"reads multiline content": {
			content: "line 1\nline 2\nline 3",
			validate: func(t *testing.T, content string) {
				if !strings.Contains(content, "line 2") {
					t.Error("expected content to contain 'line 2'")
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.txt")
			os.WriteFile(path, []byte(tc.content), 0644)

			content := readFile(t, path)
			tc.validate(t, content)
		})
	}
}

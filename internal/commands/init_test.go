package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/config"
)

func TestInitCmd(t *testing.T) {
	tests := map[string]struct {
		existingFiles map[string]string // file path -> content
		args          []string
		wantErr       string
		checkFunc     func(t *testing.T, baseDir string)
	}{
		"creates directory structure in empty repo": {
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				assertFileExists(t, filepath.Join(baseDir, ".agent-instruction"))
				assertFileExists(t, filepath.Join(baseDir, ".agent-instruction", "config.json"))
				assertFileExists(t, filepath.Join(baseDir, ".agent-instruction", "rules"))
				assertFileExists(t, filepath.Join(baseDir, ".agent-instruction", "rules", "global.json"))
			},
		},
		"creates config with default values in non-interactive mode": {
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				cfg := loadConfig(t, baseDir)
				if cfg.Version != "1.0" {
					t.Errorf("version = %q, want %q", cfg.Version, "1.0")
				}
				if len(cfg.Frameworks) != 2 {
					t.Errorf("len(frameworks) = %d, want 2", len(cfg.Frameworks))
				}
				wantFrameworks := map[string]bool{"claude": true, "agents": true}
				for _, fw := range cfg.Frameworks {
					if !wantFrameworks[fw] {
						t.Errorf("unexpected framework: %s", fw)
					}
				}
				if len(cfg.Packages) != 1 || cfg.Packages[0] != "auto" {
					t.Errorf("packages = %v, want [auto]", cfg.Packages)
				}
			},
		},
		"creates global rule file with template": {
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				rule := loadRuleFile(t, baseDir, "global.json")
				if rule.Title == "" {
					t.Error("global rule title is empty")
				}
				if len(rule.Instructions) == 0 {
					t.Error("global rule has no instructions")
				}
			},
		},
		"errors when already initialized": {
			existingFiles: map[string]string{
				".agent-instruction/config.json": `{"version":"1.0"}`,
			},
			args:    []string{"init", "--non-interactive"},
			wantErr: "already initialized",
		},
		"creates backups of existing CLAUDE.md in non-interactive mode": {
			existingFiles: map[string]string{
				"CLAUDE.md": "# Existing content",
			},
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				assertFileExists(t, filepath.Join(baseDir, "CLAUDE.md.backup"))
				content := readFile(t, filepath.Join(baseDir, "CLAUDE.md.backup"))
				if !strings.Contains(content, "Existing content") {
					t.Error("backup doesn't contain original content")
				}
			},
		},
		"creates backups of existing AGENTS.md in non-interactive mode": {
			existingFiles: map[string]string{
				"AGENTS.md": "# Agent instructions",
			},
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				assertFileExists(t, filepath.Join(baseDir, "AGENTS.md.backup"))
				content := readFile(t, filepath.Join(baseDir, "AGENTS.md.backup"))
				if !strings.Contains(content, "Agent instructions") {
					t.Error("backup doesn't contain original content")
				}
			},
		},
		"supports --frameworks flag": {
			args: []string{"init", "--non-interactive", "--frameworks", "claude"},
			checkFunc: func(t *testing.T, baseDir string) {
				cfg := loadConfig(t, baseDir)
				if len(cfg.Frameworks) != 1 || cfg.Frameworks[0] != "claude" {
					t.Errorf("frameworks = %v, want [claude]", cfg.Frameworks)
				}
			},
		},
		"supports --packages flag": {
			args: []string{"init", "--non-interactive", "--packages", "app,lib"},
			checkFunc: func(t *testing.T, baseDir string) {
				cfg := loadConfig(t, baseDir)
				if len(cfg.Packages) != 2 {
					t.Errorf("len(packages) = %d, want 2", len(cfg.Packages))
				}
				wantPackages := map[string]bool{"app": true, "lib": true}
				for _, pkg := range cfg.Packages {
					if !wantPackages[pkg] {
						t.Errorf("unexpected package: %s", pkg)
					}
				}
			},
		},
		"displays success message with next steps": {
			args: []string{"init", "--non-interactive"},
			checkFunc: func(t *testing.T, baseDir string) {
				// Success message is tested by checking output in root command test
				// Here we just verify the files were created
				assertFileExists(t, filepath.Join(baseDir, ".agent-instruction", "config.json"))
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create temp directory for test
			baseDir := t.TempDir()

			// Create existing files if specified
			for path, content := range tt.existingFiles {
				fullPath := filepath.Join(baseDir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("create dir for existing file: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("write existing file: %v", err)
				}
			}

			// Create root command and add init
			rootCmd := NewRootCmd()
			initCmd := newInitCmd()
			rootCmd.AddCommand(initCmd)

			// Change to test directory
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("get working directory: %v", err)
			}
			defer os.Chdir(oldDir)

			if err := os.Chdir(baseDir); err != nil {
				t.Fatalf("change directory: %v", err)
			}

			// Execute command
			output, err := executeCommand(rootCmd, tt.args...)

			// Check error
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v\noutput: %s", err, output)
			}

			// Run additional checks
			if tt.checkFunc != nil {
				tt.checkFunc(t, baseDir)
			}
		})
	}
}

func TestInitCmdValidation(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"rejects unexpected arguments": {
			args:    []string{"init", "extra-arg"},
			wantErr: "unknown command",
		},
		"validates frameworks flag": {
			args:    []string{"init", "--non-interactive", "--frameworks", "invalid"},
			wantErr: "invalid framework",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			baseDir := t.TempDir()

			rootCmd := NewRootCmd()
			initCmd := newInitCmd()
			rootCmd.AddCommand(initCmd)

			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("get working directory: %v", err)
			}
			defer os.Chdir(oldDir)

			if err := os.Chdir(baseDir); err != nil {
				t.Fatalf("change directory: %v", err)
			}

			_, err = executeCommand(rootCmd, tt.args...)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want error containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tests := map[string]struct {
		frameworks []string
		packages   []string
		want       config.Config
	}{
		"creates config with provided values": {
			frameworks: []string{"claude"},
			packages:   []string{"app", "lib"},
			want: config.Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{"app", "lib"},
			},
		},
		"creates config with both frameworks": {
			frameworks: []string{"claude", "agents"},
			packages:   []string{"auto"},
			want: config.Config{
				Version:    "1.0",
				Frameworks: []string{"claude", "agents"},
				Packages:   []string{"auto"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := createDefaultConfig(tt.frameworks, tt.packages)

			if got.Version != tt.want.Version {
				t.Errorf("Version = %q, want %q", got.Version, tt.want.Version)
			}
			if len(got.Frameworks) != len(tt.want.Frameworks) {
				t.Errorf("len(Frameworks) = %d, want %d", len(got.Frameworks), len(tt.want.Frameworks))
			}
			for i, fw := range got.Frameworks {
				if fw != tt.want.Frameworks[i] {
					t.Errorf("Frameworks[%d] = %q, want %q", i, fw, tt.want.Frameworks[i])
				}
			}
			if len(got.Packages) != len(tt.want.Packages) {
				t.Errorf("len(Packages) = %d, want %d", len(got.Packages), len(tt.want.Packages))
			}
			for i, pkg := range got.Packages {
				if pkg != tt.want.Packages[i] {
					t.Errorf("Packages[%d] = %q, want %q", i, pkg, tt.want.Packages[i])
				}
			}
		})
	}
}

func TestFindExistingInstructionFiles(t *testing.T) {
	tests := map[string]struct {
		files []string
		want  []string
	}{
		"finds CLAUDE.md in root": {
			files: []string{"CLAUDE.md"},
			want:  []string{"CLAUDE.md"},
		},
		"finds AGENTS.md in root": {
			files: []string{"AGENTS.md"},
			want:  []string{"AGENTS.md"},
		},
		"finds both files": {
			files: []string{"CLAUDE.md", "AGENTS.md"},
			want:  []string{"CLAUDE.md", "AGENTS.md"},
		},
		"returns empty when no files found": {
			files: []string{},
			want:  []string{},
		},
		"ignores files in subdirectories": {
			files: []string{"packages/app/CLAUDE.md"},
			want:  []string{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			baseDir := t.TempDir()

			// Create files
			for _, file := range tt.files {
				fullPath := filepath.Join(baseDir, file)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("create dir: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
					t.Fatalf("create file: %v", err)
				}
			}

			got := findExistingInstructionFiles(baseDir)

			if len(got) != len(tt.want) {
				t.Errorf("len(got) = %d, want %d", len(got), len(tt.want))
			}
			for i, path := range got {
				if i >= len(tt.want) {
					break
				}
				if path != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, path, tt.want[i])
				}
			}
		})
	}
}

package builder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/config"
	"github.com/validkeys/agent-instruction/internal/files"
	"github.com/validkeys/agent-instruction/internal/rules"
)

// TestBuildFile tests the complete BuildFile workflow end-to-end
func TestBuildFile(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) (rulePath, outputPath string)
		wantErr  bool
		errMsg   string
		validate func(t *testing.T, outputPath string)
	}{
		"creates new file with generated content": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()

				// Create rule file
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)
				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "Test Rules",
					Instructions: []rules.Instruction{
						{
							Heading: "Rule One",
							Rule:    "Always write tests first",
						},
						{
							Rule: "Use table-driven tests",
						},
					},
				})

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				assertFileExists(t, outputPath)
				content := readFileContent(t, outputPath)

				// Check for markers
				if !strings.Contains(content, files.BeginMarker) {
					t.Error("expected begin marker in output")
				}
				if !strings.Contains(content, files.EndMarker) {
					t.Error("expected end marker in output")
				}

				// Check for generated content
				if !strings.Contains(content, "## Rule One") {
					t.Error("expected heading in output")
				}
				if !strings.Contains(content, "Always write tests first") {
					t.Error("expected first rule in output")
				}
				if !strings.Contains(content, "Use table-driven tests") {
					t.Error("expected second rule in output")
				}
			},
		},
		"preserves existing user content before managed section": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()

				// Create rule file
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)
				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "New Rules",
					Instructions: []rules.Instruction{
						{Rule: "Updated rule content"},
					},
				})

				// Create existing file with user content before managed section
				outputPath := filepath.Join(dir, "CLAUDE.md")
				existingContent := "# My Custom Header\n\nUser content here.\n\n" +
					files.BeginMarker + "\nOld generated content\n" + files.EndMarker
				os.WriteFile(outputPath, []byte(existingContent), 0644)

				return rulePath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				content := readFileContent(t, outputPath)

				// User content before should be preserved
				if !strings.Contains(content, "# My Custom Header") {
					t.Error("expected user header to be preserved")
				}
				if !strings.Contains(content, "User content here.") {
					t.Error("expected user content to be preserved")
				}

				// Generated content should be updated
				if !strings.Contains(content, "Updated rule content") {
					t.Error("expected new rule content")
				}
				if strings.Contains(content, "Old generated content") {
					t.Error("old content should be replaced")
				}

				// Backup should exist
				backupPath := outputPath + ".backup"
				assertFileExists(t, backupPath)
			},
		},
		"preserves existing user content after managed section": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()

				// Create rule file
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)
				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "Rules",
					Instructions: []rules.Instruction{
						{Rule: "New content"},
					},
				})

				// Create existing file with user content after managed section
				outputPath := filepath.Join(dir, "CLAUDE.md")
				existingContent := files.BeginMarker + "\nOld content\n" + files.EndMarker +
					"\n\n# Footer Section\n\nUser notes at the end."
				os.WriteFile(outputPath, []byte(existingContent), 0644)

				return rulePath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				content := readFileContent(t, outputPath)

				// User content after should be preserved
				if !strings.Contains(content, "# Footer Section") {
					t.Error("expected user footer to be preserved")
				}
				if !strings.Contains(content, "User notes at the end.") {
					t.Error("expected user notes to be preserved")
				}

				// Generated content should be updated
				if !strings.Contains(content, "New content") {
					t.Error("expected new rule content")
				}
			},
		},
		"resolves imports and includes all instructions": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				// Create imported rule
				importedPath := filepath.Join(ruleDir, "base.json")
				createTestRuleFile(t, importedPath, &rules.RuleFile{
					Title: "Base Rules",
					Instructions: []rules.Instruction{
						{Rule: "Imported rule"},
					},
				})

				// Create main rule that imports base
				mainPath := filepath.Join(ruleDir, "main.json")
				createTestRuleFile(t, mainPath, &rules.RuleFile{
					Title: "Main Rules",
					Imports: []string{
						"base.json",
					},
					Instructions: []rules.Instruction{
						{Rule: "Main rule"},
					},
				})

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return mainPath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				content := readFileContent(t, outputPath)

				// Should contain both imported and main rules
				if !strings.Contains(content, "Imported rule") {
					t.Error("expected imported rule in output")
				}
				if !strings.Contains(content, "Main rule") {
					t.Error("expected main rule in output")
				}
			},
		},
		"handles instructions with references": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "Rules with References",
					Instructions: []rules.Instruction{
						{
							Rule: "Follow the style guide",
							References: []rules.Reference{
								{Title: "Go Style Guide", Path: "https://go.dev/style"},
								{Title: "Internal Docs", Path: "/docs/style.md"},
							},
						},
					},
				})

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				content := readFileContent(t, outputPath)

				// Check for references formatted as markdown links
				if !strings.Contains(content, "[Go Style Guide](https://go.dev/style)") {
					t.Error("expected first reference link")
				}
				if !strings.Contains(content, "[Internal Docs](/docs/style.md)") {
					t.Error("expected second reference link")
				}
			},
		},
		"creates output directory if it does not exist": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "Test",
					Instructions: []rules.Instruction{
						{Rule: "Test rule"},
					},
				})

				// Output path in non-existent nested directory
				outputPath := filepath.Join(dir, "packages", "api", "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: false,
			validate: func(t *testing.T, outputPath string) {
				assertFileExists(t, outputPath)
				content := readFileContent(t, outputPath)
				if !strings.Contains(content, "Test rule") {
					t.Error("expected rule content in output")
				}
			},
		},
		"returns error when rule path is empty": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				outputPath := filepath.Join(dir, "CLAUDE.md")
				return "", outputPath
			},
			wantErr: true,
			errMsg:  "rule path cannot be empty",
		},
		"returns error when output path is empty": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)
				rulePath := filepath.Join(ruleDir, "test.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title: "Test",
					Instructions: []rules.Instruction{
						{Rule: "Test"},
					},
				})
				return rulePath, ""
			},
			wantErr: true,
			errMsg:  "output path cannot be empty",
		},
		"returns error when rule file does not exist": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				rulePath := filepath.Join(dir, "nonexistent.json")
				outputPath := filepath.Join(dir, "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: true,
			errMsg:  "resolve rules",
		},
		"returns error when rule file has invalid JSON": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				rulePath := filepath.Join(ruleDir, "invalid.json")
				os.WriteFile(rulePath, []byte("{ invalid json }"), 0644)

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: true,
			errMsg:  "resolve rules",
		},
		"returns error when import cycle detected": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				// Create cycle: a.json -> b.json -> a.json
				aPath := filepath.Join(ruleDir, "a.json")
				createTestRuleFile(t, aPath, &rules.RuleFile{
					Title:   "A",
					Imports: []string{"b.json"},
					Instructions: []rules.Instruction{
						{Rule: "Rule A"},
					},
				})

				bPath := filepath.Join(ruleDir, "b.json")
				createTestRuleFile(t, bPath, &rules.RuleFile{
					Title:   "B",
					Imports: []string{"a.json"},
					Instructions: []rules.Instruction{
						{Rule: "Rule B"},
					},
				})

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return aPath, outputPath
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
		"returns error when rule file has no instructions": {
			setup: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				ruleDir := filepath.Join(dir, "rules")
				os.MkdirAll(ruleDir, 0755)

				rulePath := filepath.Join(ruleDir, "empty.json")
				createTestRuleFile(t, rulePath, &rules.RuleFile{
					Title:        "Empty",
					Instructions: []rules.Instruction{},
				})

				outputPath := filepath.Join(dir, "CLAUDE.md")
				return rulePath, outputPath
			},
			wantErr: true,
			errMsg:  "must contain at least one instruction",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rulePath, outputPath := tc.setup(t)

			// Create services
			configService := config.NewConfigService()
			ruleService := rules.NewRuleService(configService)
			fileService := files.NewFileService()
			buildService := NewBuildService(ruleService, fileService)

			// Execute
			err := buildService.BuildFile(rulePath, outputPath)

			// Validate error expectations
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr && tc.errMsg != "" {
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
				}
			}

			// Validate success expectations
			if tc.validate != nil {
				tc.validate(t, outputPath)
			}
		})
	}
}

// TestNewBuildService tests the constructor
func TestNewBuildService(t *testing.T) {
	configService := config.NewConfigService()
	ruleService := rules.NewRuleService(configService)
	fileService := files.NewFileService()

	service := NewBuildService(ruleService, fileService)

	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.ruleService == nil {
		t.Error("expected non-nil ruleService")
	}
	if service.fileService == nil {
		t.Error("expected non-nil fileService")
	}
}

// Helper functions

func createTestRuleFile(t *testing.T, path string, rule *rules.RuleFile) {
	t.Helper()

	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		t.Fatalf("marshal rule file: %v", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write rule file %s: %v", path, err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected file to exist: %s", path)
	}
}

func readFileContent(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(data)
}

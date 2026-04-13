package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/validkeys/agent-instruction/internal/rules"
)

func TestComposeInstructions(t *testing.T) {
	tests := map[string]struct {
		setup       func(t *testing.T) (globalPath, packagePath string, svc rules.RuleService)
		wantLen     int
		wantOrder   []string // Rule texts in expected order
		wantErr     bool
		errContains string
	}{
		"global rules only (no package config)": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title: "Global Rules",
					Instructions: []rules.Instruction{
						{Rule: "Global rule 1"},
						{Rule: "Global rule 2"},
					},
				})

				svc := createTestRuleService(t, dir)
				return globalPath, filepath.Join(dir, "nonexistent.json"), svc
			},
			wantLen:   2,
			wantOrder: []string{"Global rule 1", "Global rule 2"},
			wantErr:   false,
		},
		"package rules only (no global config)": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				packagePath := filepath.Join(dir, "package.json")
				createRuleFile(t, packagePath, &rules.RuleFile{
					Title: "Package Rules",
					Instructions: []rules.Instruction{
						{Rule: "Package rule 1"},
						{Rule: "Package rule 2"},
					},
				})

				svc := createTestRuleService(t, dir)
				return filepath.Join(dir, "nonexistent.json"), packagePath, svc
			},
			wantLen:   2,
			wantOrder: []string{"Package rule 1", "Package rule 2"},
			wantErr:   false,
		},
		"both global and package rules (global first)": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title: "Global Rules",
					Instructions: []rules.Instruction{
						{Rule: "Global rule 1"},
						{Rule: "Global rule 2"},
					},
				})

				packagePath := filepath.Join(dir, "package.json")
				createRuleFile(t, packagePath, &rules.RuleFile{
					Title: "Package Rules",
					Instructions: []rules.Instruction{
						{Rule: "Package rule 1"},
						{Rule: "Package rule 2"},
					},
				})

				svc := createTestRuleService(t, dir)
				return globalPath, packagePath, svc
			},
			wantLen:   4,
			wantOrder: []string{"Global rule 1", "Global rule 2", "Package rule 1", "Package rule 2"},
			wantErr:   false,
		},
		"global rules with imports": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				// Create base rule file
				basePath := filepath.Join(dir, "base.json")
				createRuleFile(t, basePath, &rules.RuleFile{
					Title: "Base Rules",
					Instructions: []rules.Instruction{
						{Rule: "Base rule 1"},
					},
				})

				// Create global rule that imports base
				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title:   "Global Rules",
					Imports: []string{"base.json"},
					Instructions: []rules.Instruction{
						{Rule: "Global rule 1"},
					},
				})

				svc := createTestRuleService(t, dir)
				return globalPath, filepath.Join(dir, "nonexistent.json"), svc
			},
			wantLen:   2,
			wantOrder: []string{"Base rule 1", "Global rule 1"},
			wantErr:   false,
		},
		"package rules with imports": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				// Create base rule file
				basePath := filepath.Join(dir, "base.json")
				createRuleFile(t, basePath, &rules.RuleFile{
					Title: "Base Rules",
					Instructions: []rules.Instruction{
						{Rule: "Base rule 1"},
					},
				})

				// Create package rule that imports base
				packagePath := filepath.Join(dir, "package.json")
				createRuleFile(t, packagePath, &rules.RuleFile{
					Title:   "Package Rules",
					Imports: []string{"base.json"},
					Instructions: []rules.Instruction{
						{Rule: "Package rule 1"},
					},
				})

				svc := createTestRuleService(t, dir)
				return filepath.Join(dir, "nonexistent.json"), packagePath, svc
			},
			wantLen:   2,
			wantOrder: []string{"Base rule 1", "Package rule 1"},
			wantErr:   false,
		},
		"both with imports (complex chain)": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				// Create base rule
				basePath := filepath.Join(dir, "base.json")
				createRuleFile(t, basePath, &rules.RuleFile{
					Title: "Base Rules",
					Instructions: []rules.Instruction{
						{Rule: "Base rule"},
					},
				})

				// Create global rule that imports base
				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title:   "Global Rules",
					Imports: []string{"base.json"},
					Instructions: []rules.Instruction{
						{Rule: "Global rule"},
					},
				})

				// Create package-specific base
				pkgBasePath := filepath.Join(dir, "pkg-base.json")
				createRuleFile(t, pkgBasePath, &rules.RuleFile{
					Title: "Package Base Rules",
					Instructions: []rules.Instruction{
						{Rule: "Package base rule"},
					},
				})

				// Create package rule that imports pkg-base
				packagePath := filepath.Join(dir, "package.json")
				createRuleFile(t, packagePath, &rules.RuleFile{
					Title:   "Package Rules",
					Imports: []string{"pkg-base.json"},
					Instructions: []rules.Instruction{
						{Rule: "Package rule"},
					},
				})

				svc := createTestRuleService(t, dir)
				return globalPath, packagePath, svc
			},
			wantLen:   4,
			wantOrder: []string{"Base rule", "Global rule", "Package base rule", "Package rule"},
			wantErr:   false,
		},
		"empty package config file (no instructions)": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title: "Global Rules",
					Instructions: []rules.Instruction{
						{Rule: "Global rule 1"},
					},
				})

				// Create empty package config
				packagePath := filepath.Join(dir, "package.json")
				if err := os.WriteFile(packagePath, []byte(`{"title":"Empty","instructions":[]}`), 0644); err != nil {
					t.Fatalf("write empty package config: %v", err)
				}

				svc := createTestRuleService(t, dir)
				return globalPath, packagePath, svc
			},
			wantLen:   1,
			wantOrder: []string{"Global rule 1"},
			wantErr:   false,
		},
		"neither global nor package config exists": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()
				svc := createTestRuleService(t, dir)
				return filepath.Join(dir, "nonexistent-global.json"),
					filepath.Join(dir, "nonexistent-package.json"),
					svc
			},
			wantLen: 0,
			wantErr: false,
		},
		"package imports global rules": {
			setup: func(t *testing.T) (string, string, rules.RuleService) {
				t.Helper()
				dir := t.TempDir()

				globalPath := filepath.Join(dir, "global.json")
				createRuleFile(t, globalPath, &rules.RuleFile{
					Title: "Global Rules",
					Instructions: []rules.Instruction{
						{Rule: "Global rule 1"},
					},
				})

				// Package imports global explicitly
				packagePath := filepath.Join(dir, "package.json")
				createRuleFile(t, packagePath, &rules.RuleFile{
					Title:   "Package Rules",
					Imports: []string{"global.json"},
					Instructions: []rules.Instruction{
						{Rule: "Package rule 1"},
					},
				})

				svc := createTestRuleService(t, dir)
				return globalPath, packagePath, svc
			},
			wantLen:   3,
			wantOrder: []string{"Global rule 1", "Global rule 1", "Package rule 1"},
			wantErr:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			globalPath, packagePath, svc := tc.setup(t)

			instructions, err := ComposeInstructions(globalPath, packagePath, svc)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.errContains != "" && err != nil {
				if !contains(err.Error(), tc.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.errContains)
				}
			}

			if len(instructions) != tc.wantLen {
				t.Errorf("got %d instructions, want %d", len(instructions), tc.wantLen)
			}

			// Check order
			if len(tc.wantOrder) > 0 {
				if len(instructions) != len(tc.wantOrder) {
					t.Errorf("instruction count mismatch: got %d, want %d", len(instructions), len(tc.wantOrder))
				}

				for i, want := range tc.wantOrder {
					if i >= len(instructions) {
						break
					}
					if instructions[i].Rule != want {
						t.Errorf("instruction %d: got %q, want %q", i, instructions[i].Rule, want)
					}
				}
			}
		})
	}
}

// createRuleFile creates a rule file for testing
func createRuleFile(t *testing.T, path string, rule *rules.RuleFile) {
	t.Helper()

	content := ruleToJSON(t, rule)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write rule file %s: %v", path, err)
	}
}

// ruleToJSON converts a RuleFile to JSON string for testing
func ruleToJSON(t *testing.T, rule *rules.RuleFile) string {
	t.Helper()

	// Build JSON manually for simplicity
	result := `{
  "title": "` + rule.Title + `",`

	if len(rule.Imports) > 0 {
		result += `
  "imports": [`
		for i, imp := range rule.Imports {
			if i > 0 {
				result += ","
			}
			result += `"` + imp + `"`
		}
		result += `],`
	}

	result += `
  "instructions": [`

	for i, instr := range rule.Instructions {
		if i > 0 {
			result += ","
		}
		result += `
    {`
		if instr.Heading != "" {
			result += `
      "heading": "` + instr.Heading + `",`
		}
		result += `
      "rule": "` + instr.Rule + `"`

		if len(instr.References) > 0 {
			result += `,
      "references": [`
			for j, ref := range instr.References {
				if j > 0 {
					result += ","
				}
				result += `
        {
          "title": "` + ref.Title + `",
          "path": "` + ref.Path + `"
        }`
			}
			result += `
      ]`
		}

		result += `
    }`
	}

	result += `
  ]
}`

	return result
}

// createTestRuleService creates a rule service for testing
func createTestRuleService(t *testing.T, baseDir string) rules.RuleService {
	t.Helper()

	configSvc := newTestConfigService(baseDir)
	return rules.NewRuleService(configSvc)
}

// testConfigService implements RuleConfigService for testing
type testConfigService struct {
	baseDir string
}

func newTestConfigService(baseDir string) *testConfigService {
	return &testConfigService{baseDir: baseDir}
}

func (s *testConfigService) LoadRuleFile(path string) (*rules.RuleFile, error) {
	// If path is not absolute, make it relative to baseDir
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.baseDir, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rule rules.RuleFile
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

func (s *testConfigService) SaveRuleFile(path string, rule *rules.RuleFile) error {
	return fmt.Errorf("not implemented")
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

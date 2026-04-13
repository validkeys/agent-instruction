package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockConfigService implements a simple mock for testing
type mockConfigService struct {
	files map[string]*RuleFile
}

func newMockConfigService() *mockConfigService {
	return &mockConfigService{
		files: make(map[string]*RuleFile),
	}
}

func (m *mockConfigService) LoadRuleFile(path string) (*RuleFile, error) {
	absPath, _ := filepath.Abs(path)
	if rule, ok := m.files[absPath]; ok {
		return rule, nil
	}
	// Try without absolute path
	if rule, ok := m.files[path]; ok {
		return rule, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockConfigService) addFile(path string, rule *RuleFile) {
	absPath, _ := filepath.Abs(path)
	m.files[absPath] = rule
	m.files[path] = rule // Store both for flexibility
}

func TestResolveImports(t *testing.T) {
	tests := map[string]struct {
		setup   func(*mockConfigService, string) string
		wantErr bool
		errMsg  string
		check   func(*testing.T, []Instruction)
	}{
		"simple import chain - A imports B": {
			setup: func(mock *mockConfigService, dir string) string {
				b := filepath.Join(dir, "b.json")
				mock.addFile(b, &RuleFile{
					Title: "B",
					Instructions: []Instruction{
						{Rule: "Rule from B"},
					},
				})

				a := filepath.Join(dir, "a.json")
				mock.addFile(a, &RuleFile{
					Title:   "A",
					Imports: []string{b},
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: false,
			check: func(t *testing.T, instructions []Instruction) {
				if len(instructions) != 2 {
					t.Fatalf("expected 2 instructions, got %d", len(instructions))
				}
				if instructions[0].Rule != "Rule from B" {
					t.Errorf("first instruction: got %q, want 'Rule from B'", instructions[0].Rule)
				}
				if instructions[1].Rule != "Rule from A" {
					t.Errorf("second instruction: got %q, want 'Rule from A'", instructions[1].Rule)
				}
			},
		},
		"diamond dependency - A imports B and C, both import D": {
			setup: func(mock *mockConfigService, dir string) string {
				d := filepath.Join(dir, "d.json")
				mock.addFile(d, &RuleFile{
					Title: "D",
					Instructions: []Instruction{
						{Rule: "Rule from D"},
					},
				})

				b := filepath.Join(dir, "b.json")
				mock.addFile(b, &RuleFile{
					Title:   "B",
					Imports: []string{d},
					Instructions: []Instruction{
						{Rule: "Rule from B"},
					},
				})

				c := filepath.Join(dir, "c.json")
				mock.addFile(c, &RuleFile{
					Title:   "C",
					Imports: []string{d},
					Instructions: []Instruction{
						{Rule: "Rule from C"},
					},
				})

				a := filepath.Join(dir, "a.json")
				mock.addFile(a, &RuleFile{
					Title:   "A",
					Imports: []string{b, c},
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: false,
			check: func(t *testing.T, instructions []Instruction) {
				// D should appear once (visited set prevents duplicate)
				// Order: D (via B), B, C (skips D as already visited), A
				if len(instructions) != 4 {
					t.Fatalf("expected 4 instructions, got %d", len(instructions))
				}
				expected := []string{"Rule from D", "Rule from B", "Rule from C", "Rule from A"}
				for i, want := range expected {
					if instructions[i].Rule != want {
						t.Errorf("instruction %d: got %q, want %q", i, instructions[i].Rule, want)
					}
				}
			},
		},
		"cycle detection - A imports B, B imports A": {
			setup: func(mock *mockConfigService, dir string) string {
				a := filepath.Join(dir, "a.json")
				b := filepath.Join(dir, "b.json")

				mock.addFile(b, &RuleFile{
					Title:   "B",
					Imports: []string{a},
					Instructions: []Instruction{
						{Rule: "Rule from B"},
					},
				})

				mock.addFile(a, &RuleFile{
					Title:   "A",
					Imports: []string{b},
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
		"self-import - A imports A": {
			setup: func(mock *mockConfigService, dir string) string {
				a := filepath.Join(dir, "a.json")
				mock.addFile(a, &RuleFile{
					Title:   "A",
					Imports: []string{a},
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
		"deep nesting - A→B→C→D→E": {
			setup: func(mock *mockConfigService, dir string) string {
				e := filepath.Join(dir, "e.json")
				mock.addFile(e, &RuleFile{
					Title: "E",
					Instructions: []Instruction{
						{Rule: "Rule from E"},
					},
				})

				d := filepath.Join(dir, "d.json")
				mock.addFile(d, &RuleFile{
					Title:   "D",
					Imports: []string{e},
					Instructions: []Instruction{
						{Rule: "Rule from D"},
					},
				})

				c := filepath.Join(dir, "c.json")
				mock.addFile(c, &RuleFile{
					Title:   "C",
					Imports: []string{d},
					Instructions: []Instruction{
						{Rule: "Rule from C"},
					},
				})

				b := filepath.Join(dir, "b.json")
				mock.addFile(b, &RuleFile{
					Title:   "B",
					Imports: []string{c},
					Instructions: []Instruction{
						{Rule: "Rule from B"},
					},
				})

				a := filepath.Join(dir, "a.json")
				mock.addFile(a, &RuleFile{
					Title:   "A",
					Imports: []string{b},
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: false,
			check: func(t *testing.T, instructions []Instruction) {
				if len(instructions) != 5 {
					t.Fatalf("expected 5 instructions, got %d", len(instructions))
				}
				expected := []string{"Rule from E", "Rule from D", "Rule from C", "Rule from B", "Rule from A"}
				for i, want := range expected {
					if instructions[i].Rule != want {
						t.Errorf("instruction %d: got %q, want %q", i, instructions[i].Rule, want)
					}
				}
			},
		},
		"no imports - single file": {
			setup: func(mock *mockConfigService, dir string) string {
				a := filepath.Join(dir, "a.json")
				mock.addFile(a, &RuleFile{
					Title: "A",
					Instructions: []Instruction{
						{Rule: "Rule from A"},
					},
				})

				return a
			},
			wantErr: false,
			check: func(t *testing.T, instructions []Instruction) {
				if len(instructions) != 1 {
					t.Fatalf("expected 1 instruction, got %d", len(instructions))
				}
				if instructions[0].Rule != "Rule from A" {
					t.Errorf("got %q, want 'Rule from A'", instructions[0].Rule)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock := newMockConfigService()
			dir := t.TempDir()

			rootFile := tc.setup(mock, dir)

			resolver := NewResolver(mock)
			instructions, err := resolver.ResolveImports(rootFile)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("error message:\ngot:  %q\nwant: contains %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.check != nil {
				tc.check(t, instructions)
			}
		})
	}
}

func TestResolverWithRealFiles(t *testing.T) {
	// Test with real files on disk
	dir := t.TempDir()

	// Create test files
	shared := &RuleFile{
		Title: "Shared Rules",
		Instructions: []Instruction{
			{Heading: "Shared", Rule: "Shared rule"},
		},
	}
	writeRuleFile(t, filepath.Join(dir, "shared.json"), shared)

	base := &RuleFile{
		Title:   "Base Rules",
		Imports: []string{"./shared.json"},
		Instructions: []Instruction{
			{Heading: "Base", Rule: "Base rule"},
		},
	}
	writeRuleFile(t, filepath.Join(dir, "base.json"), base)

	main := &RuleFile{
		Title:   "Main Rules",
		Imports: []string{"./base.json"},
		Instructions: []Instruction{
			{Heading: "Main", Rule: "Main rule"},
		},
	}
	mainPath := filepath.Join(dir, "main.json")
	writeRuleFile(t, mainPath, main)

	// Create config service wrapper that reads real files
	configSvc := &configServiceWrapper{
		loadFunc: func(path string) (*RuleFile, error) {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			var rule RuleFile
			if err := json.Unmarshal(data, &rule); err != nil {
				return nil, err
			}
			return &rule, nil
		},
	}

	resolver := NewResolver(configSvc)

	instructions, err := resolver.ResolveImports(mainPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(instructions))
	}

	expected := []string{"Shared rule", "Base rule", "Main rule"}
	for i, want := range expected {
		if instructions[i].Rule != want {
			t.Errorf("instruction %d: got %q, want %q", i, instructions[i].Rule, want)
		}
	}
}

// configServiceWrapper implements the minimum interface needed for testing
type configServiceWrapper struct {
	loadFunc func(string) (*RuleFile, error)
	saveFunc func(string, *RuleFile) error
}

func (c *configServiceWrapper) LoadRuleFile(path string) (*RuleFile, error) {
	return c.loadFunc(path)
}

func (c *configServiceWrapper) SaveRuleFile(path string, rule *RuleFile) error {
	if c.saveFunc != nil {
		return c.saveFunc(path, rule)
	}
	// Default implementation for tests that don't need save
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Helper to write rule files for real file tests
func writeRuleFile(t *testing.T, path string, rule *RuleFile) {
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		t.Fatalf("marshal rule: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write rule file: %v", err)
	}
}

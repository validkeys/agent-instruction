package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEdgeCases_CycleDetection(t *testing.T) {
	tests := map[string]struct {
		setup   func(string) string
		wantErr bool
		errMsg  string
	}{
		"detects cycle with real fixture files": {
			setup: func(dir string) string {
				// Use testdata fixtures
				return filepath.Join("..", "..", "testdata", "rules", "cycle-a.json")
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
		"self-import": {
			setup: func(dir string) string {
				selfPath := filepath.Join(dir, "self.json")
				rule := &RuleFile{
					Title:   "Self Import",
					Imports: []string{"./self.json"},
					Instructions: []Instruction{
						{Rule: "Self import rule"},
					},
				}
				writeEdgeCaseFile(t, selfPath, rule)
				return selfPath
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
		"indirect three-node cycle": {
			setup: func(dir string) string {
				c := filepath.Join(dir, "c.json")
				writeEdgeCaseFile(t, c, &RuleFile{
					Title:   "C",
					Imports: []string{"./a.json"},
					Instructions: []Instruction{
						{Rule: "C rule"},
					},
				})

				b := filepath.Join(dir, "b.json")
				writeEdgeCaseFile(t, b, &RuleFile{
					Title:   "B",
					Imports: []string{"./c.json"},
					Instructions: []Instruction{
						{Rule: "B rule"},
					},
				})

				a := filepath.Join(dir, "a.json")
				writeEdgeCaseFile(t, a, &RuleFile{
					Title:   "A",
					Imports: []string{"./b.json"},
					Instructions: []Instruction{
						{Rule: "A rule"},
					},
				})

				return a
			},
			wantErr: true,
			errMsg:  "import cycle",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			rootPath := tc.setup(dir)

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

			service := NewRuleService(configSvc)
			_, err := service.ResolveRules(rootPath)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("error message:\ngot:  %q\nwant: contains %q", err.Error(), tc.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEdgeCases_MissingFiles(t *testing.T) {
	tests := map[string]struct {
		setup   func(string) string
		wantErr bool
	}{
		"import file not found": {
			setup: func(dir string) string {
				return filepath.Join("..", "..", "testdata", "rules", "missing-import.json")
			},
			wantErr: true,
		},
		"root file not found": {
			setup: func(dir string) string {
				return filepath.Join(dir, "nonexistent.json")
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			rootPath := tc.setup(dir)

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

			service := NewRuleService(configSvc)
			_, err := service.ResolveRules(rootPath)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEdgeCases_MalformedJSON(t *testing.T) {
	tests := map[string]struct {
		content string
		wantErr bool
	}{
		"invalid JSON syntax": {
			content: `{ "title": "Test", "instructions": [`,
			wantErr: true,
		},
		"missing required title": {
			content: `{ "instructions": [{"rule": "Test"}] }`,
			wantErr: true,
		},
		"empty instructions array": {
			content: `{ "title": "Test", "instructions": [] }`,
			wantErr: true,
		},
		"null instructions": {
			content: `{ "title": "Test", "instructions": null }`,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.json")
			if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			configSvc := &configServiceWrapper{
				loadFunc: func(p string) (*RuleFile, error) {
					data, err := os.ReadFile(p)
					if err != nil {
						return nil, err
					}
					var rule RuleFile
					if err := json.Unmarshal(data, &rule); err != nil {
						return nil, err
					}
					// Validate
					if err := rule.Validate(); err != nil {
						return nil, err
					}
					return &rule, nil
				},
			}

			service := NewRuleService(configSvc)
			_, err := service.ResolveRules(path)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEdgeCases_DeepNesting(t *testing.T) {
	dir := t.TempDir()

	// Create 15 levels of nesting
	depth := 15
	var prevPath string

	for i := depth - 1; i >= 0; i-- {
		var path string
		if i == depth-1 {
			path = filepath.Join(dir, "level-15.json")
		} else {
			path = filepath.Join(dir, "level-"+string(rune('0'+i))+".json")
		}

		var imports []string
		if prevPath != "" {
			imports = []string{"./" + filepath.Base(prevPath)}
		}

		rule := &RuleFile{
			Title:   filepath.Base(path),
			Imports: imports,
			Instructions: []Instruction{
				{Rule: "Rule at level " + filepath.Base(path)},
			},
		}
		writeEdgeCaseFile(t, path, rule)
		prevPath = path
	}

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

	service := NewRuleService(configSvc)
	instructions, err := service.ResolveRules(prevPath)
	if err != nil {
		t.Fatalf("deep nesting failed: %v", err)
	}

	if len(instructions) != depth {
		t.Errorf("expected %d instructions, got %d", depth, len(instructions))
	}
}

func TestEdgeCases_EmptyImportsArray(t *testing.T) {
	dir := t.TempDir()

	rule := &RuleFile{
		Title:   "Empty Imports",
		Imports: []string{},
		Instructions: []Instruction{
			{Rule: "Single rule"},
		},
	}
	path := filepath.Join(dir, "empty-imports.json")
	writeEdgeCaseFile(t, path, rule)

	configSvc := &configServiceWrapper{
		loadFunc: func(p string) (*RuleFile, error) {
			data, err := os.ReadFile(p)
			if err != nil {
				return nil, err
			}
			var r RuleFile
			if err := json.Unmarshal(data, &r); err != nil {
				return nil, err
			}
			return &r, nil
		},
	}

	service := NewRuleService(configSvc)
	instructions, err := service.ResolveRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instructions) != 1 {
		t.Errorf("expected 1 instruction, got %d", len(instructions))
	}
}

func TestEdgeCases_ComplexDiamondWithMultiplePaths(t *testing.T) {
	dir := t.TempDir()

	// Create a more complex diamond: A imports B and C, B imports D and E, C imports D and E
	e := filepath.Join(dir, "e.json")
	writeEdgeCaseFile(t, e, &RuleFile{
		Title:        "E",
		Instructions: []Instruction{{Rule: "E"}},
	})

	d := filepath.Join(dir, "d.json")
	writeEdgeCaseFile(t, d, &RuleFile{
		Title:        "D",
		Instructions: []Instruction{{Rule: "D"}},
	})

	b := filepath.Join(dir, "b.json")
	writeEdgeCaseFile(t, b, &RuleFile{
		Title:   "B",
		Imports: []string{"./d.json", "./e.json"},
		Instructions: []Instruction{
			{Rule: "B"},
		},
	})

	c := filepath.Join(dir, "c.json")
	writeEdgeCaseFile(t, c, &RuleFile{
		Title:   "C",
		Imports: []string{"./d.json", "./e.json"},
		Instructions: []Instruction{
			{Rule: "C"},
		},
	})

	a := filepath.Join(dir, "a.json")
	writeEdgeCaseFile(t, a, &RuleFile{
		Title:   "A",
		Imports: []string{"./b.json", "./c.json"},
		Instructions: []Instruction{
			{Rule: "A"},
		},
	})

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

	service := NewRuleService(configSvc)
	instructions, err := service.ResolveRules(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have D, E, B, C, A (visited set prevents duplicates)
	if len(instructions) != 5 {
		t.Errorf("expected 5 unique instructions, got %d", len(instructions))
	}
}

// Benchmarks

func BenchmarkResolveImports_SimpleChain(b *testing.B) {
	dir := b.TempDir()

	// Create chain: A -> B -> C
	c := filepath.Join(dir, "c.json")
	writeEdgeCaseFile(b, c, &RuleFile{
		Title:        "C",
		Instructions: []Instruction{{Rule: "C"}},
	})

	bPath := filepath.Join(dir, "b.json")
	writeEdgeCaseFile(b, bPath, &RuleFile{
		Title:   "B",
		Imports: []string{"./c.json"},
		Instructions: []Instruction{
			{Rule: "B"},
		},
	})

	a := filepath.Join(dir, "a.json")
	writeEdgeCaseFile(b, a, &RuleFile{
		Title:   "A",
		Imports: []string{"./b.json"},
		Instructions: []Instruction{
			{Rule: "A"},
		},
	})

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

	service := NewRuleService(configSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ResolveRules(a)
		if err != nil {
			b.Fatalf("benchmark error: %v", err)
		}
	}
}

func BenchmarkResolveImports_DeepNesting(b *testing.B) {
	dir := b.TempDir()

	// Create 20 levels deep
	depth := 20
	var prevPath string

	for i := depth - 1; i >= 0; i-- {
		var path string
		if i == depth-1 {
			path = filepath.Join(dir, "level-20.json")
		} else if prevPath == "" {
			path = filepath.Join(dir, "level-0.json")
		} else {
			// Simpler naming
			path = filepath.Join(dir, filepath.Base(prevPath)[:len(filepath.Base(prevPath))-5]+"-child.json")
			if i == depth-2 {
				path = filepath.Join(dir, "level-19.json")
			}
		}

		var imports []string
		if prevPath != "" {
			imports = []string{"./" + filepath.Base(prevPath)}
		}

		rule := &RuleFile{
			Title:   filepath.Base(path),
			Imports: imports,
			Instructions: []Instruction{
				{Rule: "Rule " + filepath.Base(path)},
			},
		}
		writeEdgeCaseFile(b, path, rule)
		prevPath = path
	}

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

	service := NewRuleService(configSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ResolveRules(prevPath)
		if err != nil {
			b.Fatalf("benchmark error: %v", err)
		}
	}
}

// Helper function that works for both testing.T and testing.B
func writeEdgeCaseFile(tb testing.TB, path string, rule *RuleFile) {
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		tb.Fatalf("marshal rule: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		tb.Fatalf("write rule file: %v", err)
	}
}

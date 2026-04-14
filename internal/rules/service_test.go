package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRuleService_ResolveRules(t *testing.T) {
	tests := map[string]struct {
		setup       func(t *testing.T) (string, *configServiceWrapper)
		wantErr     bool
		wantCount   int
		wantRules   []string
	}{
		"resolves imports and returns merged instructions": {
			setup: func(t *testing.T) (string, *configServiceWrapper) {
				t.Helper()
				dir := t.TempDir()
				shared := &RuleFile{
					Title:        "Shared",
					Instructions: []Instruction{{Rule: "Shared rule"}},
				}
				writeTestRuleFile(t, filepath.Join(dir, "shared.json"), shared)
				main := &RuleFile{
					Title:        "Main",
					Imports:      []string{"./shared.json"},
					Instructions: []Instruction{{Rule: "Main rule"}},
				}
				mainPath := filepath.Join(dir, "main.json")
				writeTestRuleFile(t, mainPath, main)
				svc := &configServiceWrapper{
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
				return mainPath, svc
			},
			wantErr:   false,
			wantCount: 2,
			wantRules: []string{"Shared rule", "Main rule"},
		},
		"returns error when rule file does not exist": {
			setup: func(t *testing.T) (string, *configServiceWrapper) {
				t.Helper()
				svc := &configServiceWrapper{
					loadFunc: func(path string) (*RuleFile, error) {
						return nil, os.ErrNotExist
					},
				}
				return "/nonexistent/rule.json", svc
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rootPath, svc := tc.setup(t)
			service := NewRuleService(svc)

			instructions, err := service.ResolveRules(rootPath)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr {
				if len(instructions) != tc.wantCount {
					t.Fatalf("instruction count: got %d, want %d", len(instructions), tc.wantCount)
				}
				for i, want := range tc.wantRules {
					if instructions[i].Rule != want {
						t.Errorf("instruction[%d]: got %q, want %q", i, instructions[i].Rule, want)
					}
				}
			}
		})
	}
}

func TestRuleService_LoadRuleFile(t *testing.T) {
	tests := map[string]struct {
		setup     func(t *testing.T) (string, *configServiceWrapper)
		wantErr   bool
		wantTitle string
	}{
		"loads rule file successfully": {
			setup: func(t *testing.T) (string, *configServiceWrapper) {
				t.Helper()
				dir := t.TempDir()
				rule := &RuleFile{
					Title:        "Test Rule",
					Instructions: []Instruction{{Rule: "Test instruction"}},
				}
				path := filepath.Join(dir, "test.json")
				writeTestRuleFile(t, path, rule)
				svc := &configServiceWrapper{
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
				return path, svc
			},
			wantErr:   false,
			wantTitle: "Test Rule",
		},
		"returns error when file does not exist": {
			setup: func(t *testing.T) (string, *configServiceWrapper) {
				t.Helper()
				svc := &configServiceWrapper{
					loadFunc: func(p string) (*RuleFile, error) {
						return nil, os.ErrNotExist
					},
				}
				return "/nonexistent/rule.json", svc
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path, svc := tc.setup(t)
			service := NewRuleService(svc)

			loaded, err := service.LoadRuleFile(path)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr && loaded.Title != tc.wantTitle {
				t.Errorf("title: got %q, want %q", loaded.Title, tc.wantTitle)
			}
		})
	}
}

func TestRuleService_SaveRuleFile(t *testing.T) {
	tests := map[string]struct {
		rule    *RuleFile
		wantErr bool
	}{
		"delegates save to config service and persists rule": {
			rule: &RuleFile{
				Title:        "Test Rule",
				Instructions: []Instruction{{Rule: "Test instruction"}},
			},
			wantErr: false,
		},
		"returns error when config service save fails": {
			rule: &RuleFile{
				Title:        "Failing Rule",
				Instructions: []Instruction{{Rule: "Test instruction"}},
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.json")
			saved := false
			var savedRule *RuleFile

			configSvc := &mockConfigServiceForSave{
				saveFunc: func(p string, r *RuleFile) error {
					if tc.wantErr {
						return os.ErrPermission
					}
					saved = true
					savedRule = r
					data, err := json.MarshalIndent(r, "", "  ")
					if err != nil {
						return err
					}
					return os.WriteFile(p, data, 0644)
				},
			}
			service := &DefaultRuleService{configService: configSvc}

			err := service.SaveRuleFile(path, tc.rule)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr {
				if !saved {
					t.Error("save func was not called")
				}
				if savedRule == nil || savedRule.Title != tc.rule.Title {
					t.Errorf("saved title: got %q, want %q", savedRule.Title, tc.rule.Title)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read saved file: %v", err)
				}
				var loaded RuleFile
				if err := json.Unmarshal(data, &loaded); err != nil {
					t.Fatalf("parse saved file: %v", err)
				}
				if loaded.Title != tc.rule.Title {
					t.Errorf("loaded title: got %q, want %q", loaded.Title, tc.rule.Title)
				}
			}
		})
	}
}

func TestRuleService_AddInstruction(t *testing.T) {
	tests := map[string]struct {
		existing    *RuleFile
		instruction Instruction
		wantErr     bool
		wantRules   []string
	}{
		"appends instruction to existing rule file": {
			existing: &RuleFile{
				Title:        "Test Rule",
				Instructions: []Instruction{{Rule: "Existing instruction"}},
			},
			instruction: Instruction{Heading: "New Section", Rule: "New instruction"},
			wantErr:     false,
			wantRules:   []string{"Existing instruction", "New instruction"},
		},
		"returns error when rule file does not exist": {
			existing:    nil, // file will not be created
			instruction: Instruction{Rule: "New instruction"},
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.json")

			if tc.existing != nil {
				writeTestRuleFile(t, path, tc.existing)
			}

			configSvc := &fullConfigService{
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
				saveFunc: func(p string, r *RuleFile) error {
					data, err := json.MarshalIndent(r, "", "  ")
					if err != nil {
						return err
					}
					return os.WriteFile(p, data, 0644)
				},
			}
			service := &DefaultRuleService{configService: configSvc}

			err := service.AddInstruction(path, tc.instruction)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				var loaded RuleFile
				if err := json.Unmarshal(data, &loaded); err != nil {
					t.Fatalf("parse file: %v", err)
				}
				if len(loaded.Instructions) != len(tc.wantRules) {
					t.Fatalf("instruction count: got %d, want %d", len(loaded.Instructions), len(tc.wantRules))
				}
				for i, want := range tc.wantRules {
					if loaded.Instructions[i].Rule != want {
						t.Errorf("instruction[%d]: got %q, want %q", i, loaded.Instructions[i].Rule, want)
					}
				}
			}
		})
	}
}

// Test helpers and mocks

func writeTestRuleFile(t *testing.T, path string, rule *RuleFile) {
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		t.Fatalf("marshal rule: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write rule file: %v", err)
	}
}

type mockConfigServiceForSave struct {
	saveFunc func(string, *RuleFile) error
}

func (m *mockConfigServiceForSave) LoadRuleFile(path string) (*RuleFile, error) {
	return nil, os.ErrNotExist
}

func (m *mockConfigServiceForSave) SaveRuleFile(path string, rule *RuleFile) error {
	return m.saveFunc(path, rule)
}

type fullConfigService struct {
	loadFunc func(string) (*RuleFile, error)
	saveFunc func(string, *RuleFile) error
}

func (f *fullConfigService) LoadRuleFile(path string) (*RuleFile, error) {
	return f.loadFunc(path)
}

func (f *fullConfigService) SaveRuleFile(path string, rule *RuleFile) error {
	return f.saveFunc(path, rule)
}

// Tests for FileConfigService

func TestFileConfigService_LoadRuleFile(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) (string, string) // returns baseDir and filePath
		wantErr bool
		want    *RuleFile
	}{
		"loads rule with absolute path": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rule.json")
				rule := &RuleFile{
					Title: "Absolute Path Test",
					Instructions: []Instruction{
						{Rule: "Test rule"},
					},
				}
				writeTestRuleFile(t, path, rule)
				return dir, path
			},
			wantErr: false,
			want: &RuleFile{
				Title: "Absolute Path Test",
			},
		},
		"loads rule with relative path": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rule.json")
				rule := &RuleFile{
					Title: "Relative Path Test",
					Instructions: []Instruction{
						{Rule: "Test rule"},
					},
				}
				writeTestRuleFile(t, path, rule)
				return dir, "rule.json"
			},
			wantErr: false,
			want: &RuleFile{
				Title: "Relative Path Test",
			},
		},
		"returns error for nonexistent file": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				return dir, "nonexistent.json"
			},
			wantErr: true,
		},
		"returns error for invalid JSON": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "invalid.json")
				if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return dir, path
			},
			wantErr: true,
		},
		"returns error for invalid rule": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "invalid.json")
				// Missing required title
				rule := &RuleFile{}
				data, _ := json.MarshalIndent(rule, "", "  ")
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return dir, path
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			baseDir, filePath := tc.setup(t)

			svc := NewFileConfigService(baseDir)
			got, err := svc.LoadRuleFile(filePath)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr {
				if got.Title != tc.want.Title {
					t.Errorf("Title: got %q, want %q", got.Title, tc.want.Title)
				}
			}
		})
	}
}

func TestFileConfigService_SaveRuleFile(t *testing.T) {
	tests := map[string]struct {
		baseDir  func(t *testing.T) string
		filePath string
		rule     *RuleFile
		wantErr  bool
	}{
		"saves rule with absolute path": {
			baseDir: func(t *testing.T) string {
				return t.TempDir()
			},
			filePath: "", // Will be set in test
			rule: &RuleFile{
				Title: "Absolute Path Save",
				Instructions: []Instruction{
					{Rule: "Test rule"},
				},
			},
			wantErr: false,
		},
		"saves rule with relative path": {
			baseDir: func(t *testing.T) string {
				return t.TempDir()
			},
			filePath: "output.json",
			rule: &RuleFile{
				Title: "Relative Path Save",
				Instructions: []Instruction{
					{Rule: "Test rule"},
				},
			},
			wantErr: false,
		},
		"returns error for invalid rule": {
			baseDir: func(t *testing.T) string {
				return t.TempDir()
			},
			filePath: "invalid.json",
			rule:     &RuleFile{}, // Missing required title
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			baseDir := tc.baseDir(t)
			svc := NewFileConfigService(baseDir)

			filePath := tc.filePath
			if filePath == "" {
				filePath = filepath.Join(baseDir, "output.json")
			}

			err := svc.SaveRuleFile(filePath, tc.rule)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr {
				// Verify file was written
				fullPath := filePath
				if !filepath.IsAbs(filePath) {
					fullPath = filepath.Join(baseDir, filePath)
				}

				data, err := os.ReadFile(fullPath)
				if err != nil {
					t.Fatalf("failed to read saved file: %v", err)
				}

				var loaded RuleFile
				if err := json.Unmarshal(data, &loaded); err != nil {
					t.Fatalf("failed to parse saved file: %v", err)
				}

				if loaded.Title != tc.rule.Title {
					t.Errorf("Title: got %q, want %q", loaded.Title, tc.rule.Title)
				}
			}
		})
	}
}

func TestRuleService_AddInstruction_ErrorPath(t *testing.T) {
	// Test the error path when LoadRuleFile fails
	configSvc := &fullConfigService{
		loadFunc: func(p string) (*RuleFile, error) {
			return nil, os.ErrNotExist
		},
		saveFunc: func(p string, r *RuleFile) error {
			return nil
		},
	}

	service := NewRuleService(configSvc)

	err := service.AddInstruction("nonexistent.json", Instruction{Rule: "Test"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

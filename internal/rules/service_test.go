package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRuleService_ResolveRules(t *testing.T) {
	dir := t.TempDir()

	// Create test rule files
	shared := &RuleFile{
		Title: "Shared",
		Instructions: []Instruction{
			{Rule: "Shared rule"},
		},
	}
	sharedPath := filepath.Join(dir, "shared.json")
	writeTestRuleFile(t, sharedPath, shared)

	main := &RuleFile{
		Title:   "Main",
		Imports: []string{"./shared.json"},
		Instructions: []Instruction{
			{Rule: "Main rule"},
		},
	}
	mainPath := filepath.Join(dir, "main.json")
	writeTestRuleFile(t, mainPath, main)

	// Create service with real file config service
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

	instructions, err := service.ResolveRules(mainPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instructions) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(instructions))
	}

	if instructions[0].Rule != "Shared rule" {
		t.Errorf("first instruction: got %q, want 'Shared rule'", instructions[0].Rule)
	}

	if instructions[1].Rule != "Main rule" {
		t.Errorf("second instruction: got %q, want 'Main rule'", instructions[1].Rule)
	}
}

func TestRuleService_LoadRuleFile(t *testing.T) {
	dir := t.TempDir()

	rule := &RuleFile{
		Title: "Test Rule",
		Instructions: []Instruction{
			{Rule: "Test instruction"},
		},
	}
	path := filepath.Join(dir, "test.json")
	writeTestRuleFile(t, path, rule)

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

	loaded, err := service.LoadRuleFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loaded.Title != rule.Title {
		t.Errorf("title: got %q, want %q", loaded.Title, rule.Title)
	}

	if len(loaded.Instructions) != len(rule.Instructions) {
		t.Fatalf("instructions length: got %d, want %d", len(loaded.Instructions), len(rule.Instructions))
	}
}

func TestRuleService_SaveRuleFile(t *testing.T) {
	dir := t.TempDir()

	rule := &RuleFile{
		Title: "Test Rule",
		Instructions: []Instruction{
			{Rule: "Test instruction"},
		},
	}
	path := filepath.Join(dir, "test.json")

	// Track saves
	saved := false
	var savedRule *RuleFile

	configSvc := &mockConfigServiceForSave{
		saveFunc: func(p string, r *RuleFile) error {
			saved = true
			savedRule = r
			// Actually write to disk for verification
			data, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(p, data, 0644)
		},
	}

	service := &DefaultRuleService{
		configService: configSvc,
	}

	err := service.SaveRuleFile(path, rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !saved {
		t.Error("SaveRuleFile was not called")
	}

	if savedRule == nil {
		t.Fatal("savedRule is nil")
	}

	if savedRule.Title != rule.Title {
		t.Errorf("saved title: got %q, want %q", savedRule.Title, rule.Title)
	}

	// Verify file was written
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded RuleFile
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to parse saved file: %v", err)
	}

	if loaded.Title != rule.Title {
		t.Errorf("loaded title: got %q, want %q", loaded.Title, rule.Title)
	}
}

func TestRuleService_AddInstruction(t *testing.T) {
	dir := t.TempDir()

	// Create initial rule file
	rule := &RuleFile{
		Title: "Test Rule",
		Instructions: []Instruction{
			{Rule: "Existing instruction"},
		},
	}
	path := filepath.Join(dir, "test.json")
	writeTestRuleFile(t, path, rule)

	// Create service
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

	service := &DefaultRuleService{
		configService: configSvc,
	}

	newInstruction := Instruction{
		Heading: "New Section",
		Rule:    "New instruction",
	}

	err := service.AddInstruction(path, newInstruction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was updated
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var loaded RuleFile
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to parse file: %v", err)
	}

	if len(loaded.Instructions) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(loaded.Instructions))
	}

	if loaded.Instructions[0].Rule != "Existing instruction" {
		t.Errorf("first instruction: got %q, want 'Existing instruction'", loaded.Instructions[0].Rule)
	}

	if loaded.Instructions[1].Rule != "New instruction" {
		t.Errorf("second instruction: got %q, want 'New instruction'", loaded.Instructions[1].Rule)
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

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
	return nil, nil
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

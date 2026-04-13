package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/validkeys/agent-instruction/internal/rules"
)

func TestDefaultConfigService_LoadConfig(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantErr bool
		want    *Config
	}{
		"loads valid config": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.json")
				cfg := &Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"api"},
				}
				data, _ := json.MarshalIndent(cfg, "", "  ")
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
			want: &Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{"api"},
			},
		},
		"returns error for invalid JSON": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.json")
				if err := os.WriteFile(path, []byte("invalid json"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		"returns error for invalid config": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.json")
				// Missing required fields
				cfg := &Config{}
				data, _ := json.MarshalIndent(cfg, "", "  ")
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		"returns error for nonexistent file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "nonexistent.json")
			},
			wantErr: true,
		},
	}

	svc := &DefaultConfigService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			got, err := svc.LoadConfig(path)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr {
				if got.Version != tc.want.Version {
					t.Errorf("Version: got %q, want %q", got.Version, tc.want.Version)
				}
				if len(got.Frameworks) != len(tc.want.Frameworks) {
					t.Errorf("Frameworks length mismatch")
				}
			}
		})
	}
}

func TestDefaultConfigService_SaveConfig(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		config   *Config
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"saves valid config": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "config.json")
			},
			config: &Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{},
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}

				var cfg Config
				if err := json.Unmarshal(data, &cfg); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}

				if cfg.Version != "1.0" {
					t.Errorf("Version: got %q, want %q", cfg.Version, "1.0")
				}
			},
		},
		"returns error for invalid config": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "config.json")
			},
			config: &Config{
				// Missing required fields
			},
			wantErr: true,
		},
		"uses atomic write": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.json")
				// Create existing file
				if err := os.WriteFile(path, []byte("old"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			config: &Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				// Verify no temp files left
				dir := filepath.Dir(path)
				entries, _ := os.ReadDir(dir)
				for _, e := range entries {
					if filepath.Ext(e.Name()) == ".tmp" {
						t.Error("temp file not cleaned up")
					}
				}
			},
		},
	}

	svc := &DefaultConfigService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := svc.SaveConfig(path, tc.config)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.validate != nil {
				tc.validate(t, path)
			}
		})
	}
}

func TestDefaultConfigService_LoadRuleFile(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantErr bool
		want    *rules.RuleFile
	}{
		"loads valid rule file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rule.json")
				rule := &rules.RuleFile{
					Title: "Test Rule",
					Instructions: []rules.Instruction{
						{Rule: "Always test"},
					},
				}
				data, _ := json.MarshalIndent(rule, "", "  ")
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
			want: &rules.RuleFile{
				Title: "Test Rule",
				Instructions: []rules.Instruction{
					{Rule: "Always test"},
				},
			},
		},
		"returns error for invalid JSON": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rule.json")
				if err := os.WriteFile(path, []byte("invalid"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		"returns error for invalid rule": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "rule.json")
				// Missing required title
				rule := &rules.RuleFile{}
				data, _ := json.MarshalIndent(rule, "", "  ")
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
	}

	svc := &DefaultConfigService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			got, err := svc.LoadRuleFile(path)

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

func TestDefaultConfigService_SaveRuleFile(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		rule     *rules.RuleFile
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"saves valid rule file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "rule.json")
			},
			rule: &rules.RuleFile{
				Title: "Test Rule",
				Instructions: []rules.Instruction{
					{Rule: "Always test"},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}

				var rule rules.RuleFile
				if err := json.Unmarshal(data, &rule); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}

				if rule.Title != "Test Rule" {
					t.Errorf("Title: got %q, want %q", rule.Title, "Test Rule")
				}
			},
		},
		"returns error for invalid rule": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "rule.json")
			},
			rule:    &rules.RuleFile{},
			wantErr: true,
		},
	}

	svc := &DefaultConfigService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := svc.SaveRuleFile(path, tc.rule)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.validate != nil {
				tc.validate(t, path)
			}
		})
	}
}

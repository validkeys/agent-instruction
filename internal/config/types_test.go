package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := map[string]struct {
		config  Config
		wantErr bool
		errMsg  string
	}{
		"valid config with single framework": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{},
			},
			wantErr: false,
		},
		"valid config with multiple frameworks": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"claude", "agents"},
				Packages:   []string{"api", "web"},
			},
			wantErr: false,
		},
		"missing version": {
			config: Config{
				Frameworks: []string{"claude"},
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		"missing frameworks": {
			config: Config{
				Version: "1.0",
			},
			wantErr: true,
			errMsg:  "at least one framework is required",
		},
		"invalid framework": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"invalid"},
			},
			wantErr: true,
			errMsg:  "invalid framework",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.config.Validate()

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr && tc.errMsg != "" {
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("error message %q does not contain %q", err.Error(), tc.errMsg)
				}
			}
		})
	}
}

func TestConfigMarshalUnmarshal(t *testing.T) {
	tests := map[string]struct {
		config Config
	}{
		"basic config": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{"api"},
			},
		},
		"config with multiple packages": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"claude", "agents"},
				Packages:   []string{"api", "web", "mobile"},
			},
		},
		"config with empty packages": {
			config: Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tc.config)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			// Unmarshal back
			var decoded Config
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			// Compare
			if decoded.Version != tc.config.Version {
				t.Errorf("version: got %q, want %q", decoded.Version, tc.config.Version)
			}

			if len(decoded.Frameworks) != len(tc.config.Frameworks) {
				t.Errorf("frameworks count: got %d, want %d", len(decoded.Frameworks), len(tc.config.Frameworks))
			}

			if len(decoded.Packages) != len(tc.config.Packages) {
				t.Errorf("packages count: got %d, want %d", len(decoded.Packages), len(tc.config.Packages))
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantErr bool
		want    *Config
	}{
		"loads valid config file": {
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
		"returns error for nonexistent file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "nonexistent.json")
			},
			wantErr: true,
		},
		"returns error for invalid JSON": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.json")
				if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			got, err := LoadConfig(path)

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

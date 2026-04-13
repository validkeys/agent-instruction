package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/validkeys/agent-instruction/internal/config"
)

func TestDiscoverPackages(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) (string, *config.Config)
		wantLen  int
		wantPkgs []string // relative paths from repo root
		wantErr  bool
	}{
		"auto mode finds packages with agent-instruction.json": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				// Create package structure
				createPackage(t, dir, "packages/api")
				createPackage(t, dir, "packages/web")
				createPackage(t, dir, "services/auth")

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  3,
			wantPkgs: []string{"packages/api", "packages/web", "services/auth"},
			wantErr:  false,
		},
		"auto mode with empty packages array": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"manual mode validates explicit paths": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				createPackage(t, dir, "packages/web")

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"packages/api", "packages/web"},
				}
				return dir, cfg
			},
			wantLen:  2,
			wantPkgs: []string{"packages/api", "packages/web"},
			wantErr:  false,
		},
		"manual mode skips invalid paths": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"packages/api", "packages/missing", "packages/also-missing"},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"excludes .git directory": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				createPackage(t, dir, ".git/hooks") // Should be excluded

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"excludes node_modules directory": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				createPackage(t, dir, "node_modules/some-package") // Should be excluded

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"excludes .agent-instruction directory": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				// Create .agent-instruction/rules with agent-instruction.json (should be excluded)
				rulesDir := filepath.Join(dir, ".agent-instruction", "rules")
				if err := os.MkdirAll(rulesDir, 0755); err != nil {
					t.Fatalf("create rules dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(rulesDir, "agent-instruction.json"), []byte("{}"), 0644); err != nil {
					t.Fatalf("write config: %v", err)
				}

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"excludes .backup directory": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				createPackage(t, dir, ".backup/old-package") // Should be excluded

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  1,
			wantPkgs: []string{"packages/api"},
			wantErr:  false,
		},
		"supports nested packages": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				createPackage(t, dir, "packages/api")
				createPackage(t, dir, "packages/api/v2") // Nested package

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  2,
			wantPkgs: []string{"packages/api", "packages/api/v2"},
			wantErr:  false,
		},
		"returns empty array for no packages": {
			setup: func(t *testing.T) (string, *config.Config) {
				t.Helper()
				dir := t.TempDir()

				cfg := &config.Config{
					Version:    "1.0",
					Frameworks: []string{"claude"},
					Packages:   []string{"auto"},
				}
				return dir, cfg
			},
			wantLen:  0,
			wantPkgs: []string{},
			wantErr:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repoRoot, cfg := tc.setup(t)

			packages, err := DiscoverPackages(cfg, repoRoot)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(packages) != tc.wantLen {
				t.Errorf("got %d packages, want %d\nPackages: %v", len(packages), tc.wantLen, packages)
			}

			// Convert to relative paths for comparison
			relativePkgs := make([]string, len(packages))
			for i, pkg := range packages {
				rel, err := filepath.Rel(repoRoot, pkg)
				if err != nil {
					t.Fatalf("get relative path: %v", err)
				}
				relativePkgs[i] = rel
			}

			// Check that all expected packages are present
			for _, want := range tc.wantPkgs {
				found := false
				for _, got := range relativePkgs {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected package %q not found in: %v", want, relativePkgs)
				}
			}
		})
	}
}

func TestWalkForPackages(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantLen int
		wantErr bool
	}{
		"finds packages in flat structure": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				createPackage(t, dir, "api")
				createPackage(t, dir, "web")
				createPackage(t, dir, "auth")
				return dir
			},
			wantLen: 3,
			wantErr: false,
		},
		"finds packages in nested structure": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				createPackage(t, dir, "packages/frontend/app")
				createPackage(t, dir, "packages/backend/api")
				createPackage(t, dir, "services/auth")
				return dir
			},
			wantLen: 3,
			wantErr: false,
		},
		"handles empty directory": {
			setup: func(t *testing.T) string {
				t.Helper()
				return t.TempDir()
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repoRoot := tc.setup(t)

			packages, err := walkForPackages(repoRoot)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(packages) != tc.wantLen {
				t.Errorf("got %d packages, want %d", len(packages), tc.wantLen)
			}
		})
	}
}

func TestValidatePackagePaths(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) (string, []string)
		wantLen int
		wantErr bool
	}{
		"validates all existing paths": {
			setup: func(t *testing.T) (string, []string) {
				t.Helper()
				dir := t.TempDir()
				createPackage(t, dir, "api")
				createPackage(t, dir, "web")
				return dir, []string{"api", "web"}
			},
			wantLen: 2,
			wantErr: false,
		},
		"filters out missing paths": {
			setup: func(t *testing.T) (string, []string) {
				t.Helper()
				dir := t.TempDir()
				createPackage(t, dir, "api")
				return dir, []string{"api", "web", "missing"}
			},
			wantLen: 1,
			wantErr: false,
		},
		"filters out paths without agent-instruction.json": {
			setup: func(t *testing.T) (string, []string) {
				t.Helper()
				dir := t.TempDir()
				createPackage(t, dir, "api")

				// Create directory without agent-instruction.json
				webDir := filepath.Join(dir, "web")
				if err := os.MkdirAll(webDir, 0755); err != nil {
					t.Fatalf("create dir: %v", err)
				}

				return dir, []string{"api", "web"}
			},
			wantLen: 1,
			wantErr: false,
		},
		"returns empty for all invalid paths": {
			setup: func(t *testing.T) (string, []string) {
				t.Helper()
				dir := t.TempDir()
				return dir, []string{"missing1", "missing2"}
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repoRoot, paths := tc.setup(t)

			validPaths, err := validatePackagePaths(paths, repoRoot)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(validPaths) != tc.wantLen {
				t.Errorf("got %d valid paths, want %d", len(validPaths), tc.wantLen)
			}
		})
	}
}

func TestSymlinkCycleDetection(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantErr bool
	}{
		"detects simple cycle (symlink to parent)": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()

				// Create packages/api directory
				apiDir := filepath.Join(dir, "packages", "api")
				if err := os.MkdirAll(apiDir, 0755); err != nil {
					t.Fatalf("create api dir: %v", err)
				}

				// Create symlink from api/link back to packages
				linkPath := filepath.Join(apiDir, "link")
				targetPath := filepath.Join(dir, "packages")
				if err := os.Symlink(targetPath, linkPath); err != nil {
					t.Skip("symlink not supported on this system")
				}

				createPackage(t, dir, "packages/api")

				return dir
			},
			wantErr: false, // Should handle gracefully, not error
		},
		"handles valid symlink": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()

				// Create actual package
				createPackage(t, dir, "packages/api")

				// Create symlink to package
				linkPath := filepath.Join(dir, "packages", "api-link")
				targetPath := filepath.Join(dir, "packages", "api")
				if err := os.Symlink(targetPath, linkPath); err != nil {
					t.Skip("symlink not supported on this system")
				}

				return dir
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repoRoot := tc.setup(t)

			cfg := &config.Config{
				Version:    "1.0",
				Frameworks: []string{"claude"},
				Packages:   []string{"auto"},
			}

			_, err := DiscoverPackages(cfg, repoRoot)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// createPackage creates a package directory with agent-instruction.json
func createPackage(t *testing.T, baseDir, pkgPath string) {
	t.Helper()

	fullPath := filepath.Join(baseDir, pkgPath)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		t.Fatalf("create package dir %s: %v", pkgPath, err)
	}

	configPath := filepath.Join(fullPath, "agent-instruction.json")
	content := `{
  "version": "1.0",
  "instructions": []
}`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("write agent-instruction.json for %s: %v", pkgPath, err)
	}
}

// containsPath checks if a slice contains a path (with normalized separators)
func containsPath(paths []string, target string) bool {
	normalizedTarget := filepath.ToSlash(target)
	for _, p := range paths {
		if filepath.ToSlash(p) == normalizedTarget || strings.HasSuffix(filepath.ToSlash(p), normalizedTarget) {
			return true
		}
	}
	return false
}

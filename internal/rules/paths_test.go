package rules

import (
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tests := map[string]struct {
		importPath string
		baseDir    string
		want       string
		wantErr    bool
	}{
		"relative path from base": {
			importPath: "./test.json",
			baseDir:    "/home/user/rules",
			want:       "/home/user/rules/test.json",
			wantErr:    false,
		},
		"relative path nested": {
			importPath: "./nested/test.json",
			baseDir:    "/home/user/rules",
			want:       "/home/user/rules/nested/test.json",
			wantErr:    false,
		},
		"relative path parent directory": {
			importPath: "../shared/test.json",
			baseDir:    "/home/user/rules",
			want:       "/home/user/shared/test.json",
			wantErr:    false,
		},
		"absolute path": {
			importPath: "/etc/rules/test.json",
			baseDir:    "/home/user/rules",
			want:       "/etc/rules/test.json",
			wantErr:    false,
		},
		"no leading dot relative path": {
			importPath: "test.json",
			baseDir:    "/home/user/rules",
			want:       "/home/user/rules/test.json",
			wantErr:    false,
		},
		"empty import path": {
			importPath: "",
			baseDir:    "/home/user/rules",
			want:       "",
			wantErr:    true,
		},
		"empty base dir": {
			importPath: "./test.json",
			baseDir:    "",
			want:       "",
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ResolvePath(tc.importPath, tc.baseDir)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Normalize paths for comparison (handles OS differences)
			wantAbs, _ := filepath.Abs(tc.want)
			gotAbs, _ := filepath.Abs(got)

			if gotAbs != wantAbs {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestResolveImportPath(t *testing.T) {
	tests := map[string]struct {
		importPath      string
		currentFilePath string
		want            string
		wantErr         bool
	}{
		"relative import from current dir": {
			importPath:      "./shared.json",
			currentFilePath: "/home/user/rules/main.json",
			want:            "/home/user/rules/shared.json",
			wantErr:         false,
		},
		"relative import from nested dir": {
			importPath:      "../base.json",
			currentFilePath: "/home/user/rules/nested/main.json",
			want:            "/home/user/rules/base.json",
			wantErr:         false,
		},
		"relative import to nested dir": {
			importPath:      "./sub/shared.json",
			currentFilePath: "/home/user/rules/main.json",
			want:            "/home/user/rules/sub/shared.json",
			wantErr:         false,
		},
		"absolute import path": {
			importPath:      "/etc/rules/shared.json",
			currentFilePath: "/home/user/rules/main.json",
			want:            "/etc/rules/shared.json",
			wantErr:         false,
		},
		"import without leading dot": {
			importPath:      "shared.json",
			currentFilePath: "/home/user/rules/main.json",
			want:            "/home/user/rules/shared.json",
			wantErr:         false,
		},
		"empty import path": {
			importPath:      "",
			currentFilePath: "/home/user/rules/main.json",
			want:            "",
			wantErr:         true,
		},
		"empty current file path": {
			importPath:      "./shared.json",
			currentFilePath: "",
			want:            "",
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ResolveImportPath(tc.importPath, tc.currentFilePath)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Normalize paths for comparison (handles OS differences)
			wantAbs, _ := filepath.Abs(tc.want)
			gotAbs, _ := filepath.Abs(got)

			if gotAbs != wantAbs {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

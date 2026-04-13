package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) (path string, baseDir string)
		wantErr bool
	}{
		"accepts path within base directory": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := filepath.Join(base, "file.txt")
				return path, base
			},
			wantErr: false,
		},
		"accepts nested path within base directory": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				subdir := filepath.Join(base, "subdir")
				os.Mkdir(subdir, 0755)
				path := filepath.Join(subdir, "file.txt")
				return path, base
			},
			wantErr: false,
		},
		"rejects path with .. trying to escape": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := filepath.Join(base, "..", "outside.txt")
				return path, base
			},
			wantErr: true,
		},
		"rejects absolute path outside base": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := "/tmp/outside.txt"
				return path, base
			},
			wantErr: true,
		},
		"rejects path starting with ../": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := "../outside.txt"
				return path, base
			},
			wantErr: true,
		},
		"handles symlink within base": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()

				// Create target file
				target := filepath.Join(base, "target.txt")
				os.WriteFile(target, []byte("content"), 0644)

				// Create symlink
				link := filepath.Join(base, "link.txt")
				if err := os.Symlink(target, link); err != nil {
					t.Skip("symlinks not supported")
				}

				return link, base
			},
			wantErr: false,
		},
		"rejects symlink pointing outside base": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()

				// Create target outside base
				outside := t.TempDir()
				target := filepath.Join(outside, "target.txt")
				os.WriteFile(target, []byte("content"), 0644)

				// Create symlink inside base pointing outside
				link := filepath.Join(base, "link.txt")
				if err := os.Symlink(target, link); err != nil {
					t.Skip("symlinks not supported")
				}

				return link, base
			},
			wantErr: true,
		},
		"accepts nonexistent file in valid directory": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := filepath.Join(base, "newfile.txt")
				return path, base
			},
			wantErr: false,
		},
		"accepts file in nonexistent subdirectory": {
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				base := t.TempDir()
				path := filepath.Join(base, "newdir", "file.txt")
				return path, base
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path, baseDir := tc.setup(t)

			err := ValidatePath(path, baseDir)

			if tc.wantErr && err == nil {
				t.Fatalf("expected error for path %s with base %s, got nil", path, baseDir)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for path %s with base %s: %v", path, baseDir, err)
			}
		})
	}
}

func TestIsSymlink(t *testing.T) {
	tests := map[string]struct {
		setup func(t *testing.T) string
		want  bool
	}{
		"returns true for symlink": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				target := filepath.Join(dir, "target.txt")
				os.WriteFile(target, []byte("content"), 0644)

				link := filepath.Join(dir, "link.txt")
				if err := os.Symlink(target, link); err != nil {
					t.Skip("symlinks not supported")
				}
				return link
			},
			want: true,
		},
		"returns false for regular file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "file.txt")
				os.WriteFile(path, []byte("content"), 0644)
				return path
			},
			want: false,
		},
		"returns false for nonexistent file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "nonexistent.txt")
			},
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			got, err := IsSymlink(path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

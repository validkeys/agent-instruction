package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAtomic(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string // returns filepath
		content  []byte
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"creates new file with content": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "test.txt")
			},
			content: []byte("test content"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("got %q, want %q", string(content), "test content")
				}
			},
		},
		"overwrites existing file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(path, []byte("old content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			content: []byte("new content"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				if string(content) != "new content" {
					t.Errorf("got %q, want %q", string(content), "new content")
				}
			},
		},
		"preserves existing file permissions": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(path, []byte("content"), 0600); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			content: []byte("new content"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("stat file: %v", err)
				}
				// Check that permissions are 0600
				perm := info.Mode().Perm()
				if perm != 0600 {
					t.Errorf("permissions not preserved: got %v, want 0600", perm)
				}
			},
		},
		"uses default permissions for new file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "test.txt")
			},
			content: []byte("content"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("stat file: %v", err)
				}
				// Should use default 0644
				perm := info.Mode().Perm()
				if perm != 0644 {
					t.Errorf("wrong default permissions: got %v, want 0644", perm)
				}
			},
		},
		"cleans up temp file on error": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				// Create a path to a read-only directory
				roDir := filepath.Join(dir, "readonly")
				if err := os.Mkdir(roDir, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				path := filepath.Join(roDir, "test.txt")

				// Create initial file
				if err := os.WriteFile(path, []byte("initial"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				// Make directory read-only
				if err := os.Chmod(roDir, 0555); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				// Restore on cleanup
				t.Cleanup(func() {
					os.Chmod(roDir, 0755)
				})

				return path
			},
			content: []byte("new content"),
			wantErr: true,
			validate: func(t *testing.T, path string) {
				t.Helper()
				dir := filepath.Dir(path)

				// Check no temp files left behind
				entries, err := os.ReadDir(dir)
				if err != nil {
					return // may not be readable
				}

				for _, entry := range entries {
					if filepath.Ext(entry.Name()) == ".tmp" {
						t.Errorf("temp file not cleaned up: %s", entry.Name())
					}
				}
			},
		},
		"handles empty content": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "test.txt")
			},
			content: []byte{},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				if len(content) != 0 {
					t.Errorf("expected empty file, got %d bytes", len(content))
				}
			},
		},
		"handles large content": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "test.txt")
			},
			content: make([]byte, 1024*1024), // 1MB
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				if len(content) != 1024*1024 {
					t.Errorf("size mismatch: got %d, want %d", len(content), 1024*1024)
				}
			},
		},
		"creates parent directory if needed": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				// Create directory structure
				subdir := filepath.Join(dir, "subdir")
				if err := os.Mkdir(subdir, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filepath.Join(subdir, "test.txt")
			},
			content: []byte("content"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				if _, err := os.Stat(path); err != nil {
					t.Errorf("file not created: %v", err)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := WriteAtomic(path, tc.content)

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

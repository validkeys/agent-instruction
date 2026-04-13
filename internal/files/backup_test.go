package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBackup(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string // returns filepath
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"creates backup with correct permissions": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				content := []byte("test content")
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				backupPath := path + ".backup"

				// Check backup exists
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Fatal("backup file not created")
				}

				// Check content matches
				original, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read original: %v", err)
				}
				backup, err := os.ReadFile(backupPath)
				if err != nil {
					t.Fatalf("read backup: %v", err)
				}
				if string(original) != string(backup) {
					t.Errorf("backup content mismatch:\noriginal: %s\nbackup: %s", original, backup)
				}

				// Check permissions match
				origInfo, _ := os.Stat(path)
				backupInfo, _ := os.Stat(backupPath)
				if origInfo.Mode() != backupInfo.Mode() {
					t.Errorf("permissions mismatch: original %v, backup %v", origInfo.Mode(), backupInfo.Mode())
				}
			},
		},
		"handles file with different permissions": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				content := []byte("test content")
				if err := os.WriteFile(path, content, 0600); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				backupPath := path + ".backup"

				origInfo, _ := os.Stat(path)
				backupInfo, _ := os.Stat(backupPath)
				if origInfo.Mode() != backupInfo.Mode() {
					t.Errorf("permissions not preserved: original %v, backup %v", origInfo.Mode(), backupInfo.Mode())
				}
			},
		},
		"returns nil when file doesn't exist": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "nonexistent.txt")
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				backupPath := path + ".backup"
				if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
					t.Error("backup should not exist for nonexistent file")
				}
			},
		},
		"does not overwrite existing backup": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				backupPath := path + ".backup"

				// Create original file
				if err := os.WriteFile(path, []byte("original content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				// Create existing backup
				if err := os.WriteFile(backupPath, []byte("existing backup"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				return path
			},
			wantErr: true,
			validate: func(t *testing.T, path string) {
				t.Helper()
				backupPath := path + ".backup"

				// Existing backup should be unchanged
				content, err := os.ReadFile(backupPath)
				if err != nil {
					t.Fatalf("read backup: %v", err)
				}
				if string(content) != "existing backup" {
					t.Error("existing backup was overwritten")
				}
			},
		},
		"handles permission denied": {
			setup: func(t *testing.T) string {
				t.Helper()
				if os.Geteuid() == 0 {
					t.Skip("test not valid when running as root")
				}

				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")

				// Create file
				if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				// Make directory read-only
				if err := os.Chmod(dir, 0555); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				// Restore permissions on cleanup
				t.Cleanup(func() {
					os.Chmod(dir, 0755)
				})

				return path
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := CreateBackup(path)

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

func TestBackupExists(t *testing.T) {
	tests := map[string]struct {
		setup func(t *testing.T) string
		want  bool
	}{
		"returns true when backup exists": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				backupPath := path + ".backup"

				if err := os.WriteFile(backupPath, []byte("backup"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}

				return path
			},
			want: true,
		},
		"returns false when backup doesn't exist": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "test.txt")
			},
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			got := BackupExists(path)

			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultFileService_ReadFile(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T) string
		wantErr bool
	}{
		"reads existing file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
		},
		"returns error for nonexistent file": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "nonexistent.txt")
			},
			wantErr: true,
		},
	}

	svc := &DefaultFileService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			content, err := svc.ReadFile(path)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr && string(content) != "content" {
				t.Errorf("got %q, want %q", string(content), "content")
			}
		})
	}
}

func TestDefaultFileService_WriteFile(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		content  []byte
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"writes new file": {
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
		"overwrites existing file atomically": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(path, []byte("old"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			content: []byte("new"),
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				if string(content) != "new" {
					t.Errorf("got %q, want %q", string(content), "new")
				}
			},
		},
	}

	svc := &DefaultFileService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := svc.WriteFile(path, tc.content)

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

func TestDefaultFileService_BackupFile(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		"creates backup": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			wantErr: false,
			validate: func(t *testing.T, path string) {
				t.Helper()
				backupPath := path + ".backup"
				if _, err := os.Stat(backupPath); err != nil {
					t.Errorf("backup not created: %v", err)
				}
			},
		},
	}

	svc := &DefaultFileService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := svc.BackupFile(path)

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

func TestDefaultFileService_ParseManaged(t *testing.T) {
	tests := map[string]struct {
		content []byte
		wantErr bool
		want    *ManagedContent
	}{
		"parses valid managed section": {
			content: []byte(`User content

<!-- BEGIN AGENT-INSTRUCTION -->
Managed content
<!-- END AGENT-INSTRUCTION -->

More content`),
			wantErr: false,
			want: &ManagedContent{
				Before:  "User content\n\n",
				Managed: "\nManaged content\n",
				After:   "\n\nMore content",
			},
		},
	}

	svc := &DefaultFileService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := svc.ParseManaged(tc.content)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr {
				if got.Before != tc.want.Before {
					t.Errorf("Before mismatch: got %q, want %q", got.Before, tc.want.Before)
				}
			}
		})
	}
}

func TestDefaultFileService_UpdateManaged(t *testing.T) {
	tests := map[string]struct {
		setup       func(t *testing.T) string
		newContent  string
		wantErr     bool
		validate    func(t *testing.T, path string)
		wantBackup  bool
	}{
		"updates existing managed section with backup": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				content := `User content

<!-- BEGIN AGENT-INSTRUCTION -->
Old content
<!-- END AGENT-INSTRUCTION -->

More user content`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			newContent: "New managed content",
			wantErr:    false,
			wantBackup: true,
			validate: func(t *testing.T, path string) {
				t.Helper()

				// Check backup exists
				backupPath := path + ".backup"
				backupContent, err := os.ReadFile(backupPath)
				if err != nil {
					t.Errorf("backup not created: %v", err)
				}
				if !contains(string(backupContent), "Old content") {
					t.Error("backup doesn't contain old content")
				}

				// Check updated file
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				contentStr := string(content)

				if !contains(contentStr, "New managed content") {
					t.Error("file doesn't contain new content")
				}
				if !contains(contentStr, "User content") {
					t.Error("user content was lost")
				}
				if !contains(contentStr, "More user content") {
					t.Error("trailing user content was lost")
				}
			},
		},
		"creates managed section in file without one": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "test.txt")
				content := "Just user content"
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return path
			},
			newContent: "Generated content",
			wantErr:    false,
			wantBackup: true,
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				contentStr := string(content)

				if !contains(contentStr, "Generated content") {
					t.Error("managed content not added")
				}
				if !contains(contentStr, "Just user content") {
					t.Error("original content was lost")
				}
				if !contains(contentStr, BeginMarker) || !contains(contentStr, EndMarker) {
					t.Error("markers not added")
				}
			},
		},
		"creates new file with managed section": {
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				return filepath.Join(dir, "new.txt")
			},
			newContent: "Generated content",
			wantErr:    false,
			wantBackup: false, // No backup for new file
			validate: func(t *testing.T, path string) {
				t.Helper()
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read file: %v", err)
				}
				contentStr := string(content)

				if !contains(contentStr, "Generated content") {
					t.Error("managed content not present")
				}
				if !contains(contentStr, BeginMarker) || !contains(contentStr, EndMarker) {
					t.Error("markers not present")
				}

				// Verify no backup was created
				backupPath := path + ".backup"
				if _, err := os.Stat(backupPath); err == nil {
					t.Error("backup should not exist for new file")
				}
			},
		},
	}

	svc := &DefaultFileService{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.setup(t)

			err := svc.UpdateManaged(path, tc.newContent)

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

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteAtomic writes content to a file atomically using a temp file and rename
// This prevents partial writes and corruption if interrupted
// Preserves existing file permissions or uses 0644 for new files
func WriteAtomic(path string, content []byte) error {
	// Get existing file permissions, or use default for new files
	perm := os.FileMode(0644)
	if info, err := os.Stat(path); err == nil {
		perm = info.Mode()
	}

	// Create temp file in same directory as target
	// This ensures temp and target are on same filesystem (required for atomic rename)
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Ensure cleanup on any error
	defer func() {
		_ = tempFile.Close()    // Ignore close error in cleanup
		_ = os.Remove(tempPath) // Clean up temp file (ignore error if already renamed)
	}()

	// Write content to temp file
	if _, err := tempFile.Write(content); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	// Sync to disk before rename
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("sync temp file: %w", err)
	}

	// Close before rename (required on Windows)
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	// Set permissions on temp file before rename
	if err := os.Chmod(tempPath, perm); err != nil {
		return fmt.Errorf("set permissions: %w", err)
	}

	// Atomic rename - this is the key operation
	// If this fails, original file is unchanged
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}

package files

import (
	"fmt"
	"os"
)

// CreateBackup creates a backup of the file with .backup extension
// If the file doesn't exist, returns nil (no backup needed)
// If backup already exists, returns an error (won't overwrite)
// Preserves file permissions exactly
func CreateBackup(path string) error {
	// Check if original file exists
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no file to backup
		}
		return fmt.Errorf("read file for backup: %w", err)
	}

	// Get original file permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file for backup: %w", err)
	}

	backupPath := path + ".backup"

	// Check if backup already exists
	if _, err := os.Stat(backupPath); err == nil {
		return fmt.Errorf("backup already exists at %s: will not overwrite", backupPath)
	}

	// Write backup with same permissions as original
	if err := os.WriteFile(backupPath, content, info.Mode()); err != nil {
		return fmt.Errorf("write backup: %w", err)
	}

	return nil
}

// BackupExists checks if a backup file exists for the given path
func BackupExists(path string) bool {
	backupPath := path + ".backup"
	_, err := os.Stat(backupPath)
	return err == nil
}

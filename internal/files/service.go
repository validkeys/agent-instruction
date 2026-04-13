package files

import (
	"fmt"
	"os"
)

// FileService provides unified interface for all file operations
type FileService interface {
	// ReadFile reads file content
	ReadFile(path string) ([]byte, error)

	// WriteFile writes content atomically
	WriteFile(path string, content []byte) error

	// BackupFile creates a backup with .backup extension
	BackupFile(path string) error

	// ParseManaged parses content for managed sections
	ParseManaged(content []byte) (*ManagedContent, error)

	// UpdateManaged updates managed section in a file
	// Creates backup if file exists, writes atomically
	UpdateManaged(path string, newContent string) error
}

// DefaultFileService implements FileService using standard file operations
type DefaultFileService struct{}

// ReadFile reads file content
func (s *DefaultFileService) ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}
	return content, nil
}

// WriteFile writes content to file atomically
func (s *DefaultFileService) WriteFile(path string, content []byte) error {
	if err := WriteAtomic(path, content); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

// BackupFile creates a backup of the file
func (s *DefaultFileService) BackupFile(path string) error {
	if err := CreateBackup(path); err != nil {
		return fmt.Errorf("backup file %s: %w", path, err)
	}
	return nil
}

// ParseManaged parses content for managed sections
func (s *DefaultFileService) ParseManaged(content []byte) (*ManagedContent, error) {
	managed, err := ParseManagedContent(string(content))
	if err != nil {
		return nil, fmt.Errorf("parse managed content: %w", err)
	}
	return managed, nil
}

// UpdateManaged updates the managed section in a file
// Workflow:
//  1. Read existing file (if exists)
//  2. Parse managed sections
//  3. Replace managed content
//  4. Create backup if file exists
//  5. Write atomically
func (s *DefaultFileService) UpdateManaged(path string, newContent string) error {
	// Read existing file content
	existingContent, err := os.ReadFile(path)
	fileExists := err == nil

	var updatedContent string

	if fileExists {
		// Parse and replace managed section
		updatedContent, err = ReplaceManagedSection(string(existingContent), newContent)
		if err != nil {
			return fmt.Errorf("replace managed section in %s: %w", path, err)
		}

		// Create backup before modification
		if err := CreateBackup(path); err != nil {
			return fmt.Errorf("create backup of %s: %w", path, err)
		}
	} else {
		// New file - create with managed section
		updatedContent = BeginMarker + "\n" + newContent + "\n" + EndMarker
	}

	// Write atomically
	if err := WriteAtomic(path, []byte(updatedContent)); err != nil {
		return fmt.Errorf("write updated file %s: %w", path, err)
	}

	return nil
}

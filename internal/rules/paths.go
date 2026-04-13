package rules

import (
	"errors"
	"fmt"
	"path/filepath"
)

var (
	ErrEmptyImportPath = errors.New("import path cannot be empty")
	ErrEmptyBaseDir    = errors.New("base directory cannot be empty")
	ErrEmptyFilePath   = errors.New("file path cannot be empty")
)

// ResolvePath resolves an import path relative to a base directory.
// It handles both relative and absolute paths, returning an absolute normalized path.
func ResolvePath(importPath, baseDir string) (string, error) {
	if importPath == "" {
		return "", fmt.Errorf("%w", ErrEmptyImportPath)
	}

	if baseDir == "" {
		return "", fmt.Errorf("%w", ErrEmptyBaseDir)
	}

	// If path is already absolute, return it normalized
	if filepath.IsAbs(importPath) {
		return filepath.Clean(importPath), nil
	}

	// Join relative path with base directory
	joined := filepath.Join(baseDir, importPath)

	// Convert to absolute path and normalize
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	return absPath, nil
}

// ResolveImportPath resolves an import path relative to the importing file's directory.
// It extracts the directory from currentFilePath and resolves importPath relative to it.
func ResolveImportPath(importPath, currentFilePath string) (string, error) {
	if importPath == "" {
		return "", fmt.Errorf("%w", ErrEmptyImportPath)
	}

	if currentFilePath == "" {
		return "", fmt.Errorf("%w", ErrEmptyFilePath)
	}

	// Get the directory containing the current file
	currentDir := filepath.Dir(currentFilePath)

	// Resolve the import relative to that directory
	return ResolvePath(importPath, currentDir)
}

package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath ensures path is safe and within expected directory
// Checks for:
//   - Path traversal attacks (../)
//   - Symlink attacks
//   - Paths outside base directory
func ValidatePath(path string, baseDir string) error {
	// Convert to absolute paths for comparison
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve absolute path: %w", err)
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("resolve absolute base directory: %w", err)
	}

	// Check for suspicious patterns in original path
	if strings.Contains(path, "..") {
		return fmt.Errorf("path contains suspicious pattern '..': %s", path)
	}

	// Evaluate any symlinks in the base directory
	evalBase, err := filepath.EvalSymlinks(absBase)
	if err != nil {
		// Base doesn't exist yet, use absolute path
		evalBase = absBase
	}

	// Evaluate any symlinks in the path
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If file doesn't exist yet, walk up until we find an existing parent
		// and reconstruct the path from there
		parts := []string{}
		current := absPath
		for {
			// Try to evaluate current path
			if evalCurrent, err := filepath.EvalSymlinks(current); err == nil {
				// Found an existing ancestor
				evalPath = evalCurrent
				// Append the parts we collected
				for i := len(parts) - 1; i >= 0; i-- {
					evalPath = filepath.Join(evalPath, parts[i])
				}
				break
			}

			// Move up one level
			parent := filepath.Dir(current)
			if parent == current {
				// Reached root without finding existing path
				// This shouldn't happen for temp dirs, but use absolute path
				evalPath = absPath
				break
			}

			parts = append(parts, filepath.Base(current))
			current = parent
		}
	}

	// Ensure the evaluated path is within base directory
	relPath, err := filepath.Rel(evalBase, evalPath)
	if err != nil {
		return fmt.Errorf("compute relative path: %w", err)
	}

	// Check if path escapes base directory
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("path %s is outside base directory %s", path, baseDir)
	}

	return nil
}

// IsSymlink checks if a path is a symbolic link
func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat %s: %w", path, err)
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

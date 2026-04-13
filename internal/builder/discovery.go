package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/validkeys/agent-instruction/internal/config"
)

// excludedDirs are directories to skip during package discovery
var excludedDirs = map[string]bool{
	".git":               true,
	"node_modules":       true,
	".agent-instruction": true,
	".backup":            true,
}

// DiscoverPackages finds all packages in the monorepo based on config
func DiscoverPackages(cfg *config.Config, repoRoot string) ([]string, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Check if auto mode (empty array or contains "auto")
	isAutoMode := len(cfg.Packages) == 0
	for _, pkg := range cfg.Packages {
		if pkg == "auto" {
			isAutoMode = true
			break
		}
	}

	var packages []string
	var err error

	if isAutoMode {
		// Auto mode: walk directory tree
		packages, err = walkForPackages(repoRoot)
		if err != nil {
			return nil, fmt.Errorf("walk for packages: %w", err)
		}
	} else {
		// Manual mode: validate explicit paths
		packages, err = validatePackagePaths(cfg.Packages, repoRoot)
		if err != nil {
			return nil, fmt.Errorf("validate package paths: %w", err)
		}
	}

	// Sort packages for consistent ordering
	sort.Strings(packages)

	return packages, nil
}

// walkForPackages walks the repository tree and finds all directories with agent-instruction.json
func walkForPackages(repoRoot string) ([]string, error) {
	packages := make([]string, 0)
	visited := make(map[string]bool)

	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this directory should be excluded
		dirName := info.Name()
		if excludedDirs[dirName] {
			return filepath.SkipDir
		}

		// Handle symlinks with cycle detection
		if info.Mode()&os.ModeSymlink != 0 {
			// Resolve symlink
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				// Skip broken symlinks
				return nil
			}

			// Check for cycles (symlink pointing to ancestor)
			if visited[realPath] {
				return filepath.SkipDir
			}

			// Check if symlink points outside repo (could cause infinite loops)
			if !strings.HasPrefix(realPath, repoRoot) {
				return filepath.SkipDir
			}

			visited[realPath] = true
		}

		// Check for agent-instruction.json
		configPath := filepath.Join(path, "agent-instruction.json")
		if _, err := os.Stat(configPath); err == nil {
			packages = append(packages, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory tree: %w", err)
	}

	return packages, nil
}

// validatePackagePaths validates explicit package paths and returns valid ones
func validatePackagePaths(paths []string, repoRoot string) ([]string, error) {
	validPaths := make([]string, 0, len(paths))

	for _, pkgPath := range paths {
		// Skip empty paths
		if pkgPath == "" {
			continue
		}

		// Make absolute path
		absPath := pkgPath
		if !filepath.IsAbs(pkgPath) {
			absPath = filepath.Join(repoRoot, pkgPath)
		}

		// Check if directory exists
		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Path doesn't exist - skip with warning
				continue
			}
			return nil, fmt.Errorf("stat %s: %w", pkgPath, err)
		}

		// Check if it's a directory
		if !info.IsDir() {
			continue
		}

		// Check for agent-instruction.json
		configPath := filepath.Join(absPath, "agent-instruction.json")
		if _, err := os.Stat(configPath); err != nil {
			if os.IsNotExist(err) {
				// No agent-instruction.json - skip with warning
				continue
			}
			return nil, fmt.Errorf("stat %s: %w", configPath, err)
		}

		validPaths = append(validPaths, absPath)
	}

	return validPaths, nil
}

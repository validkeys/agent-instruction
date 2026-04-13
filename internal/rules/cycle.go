package rules

import (
	"fmt"
	"strings"
)

// ImportContext tracks state during import resolution
type ImportContext struct {
	visited   map[string]bool // All paths seen during resolution
	pathStack []string        // Current import chain for cycle detection
}

// NewImportContext creates a new import context
func NewImportContext() *ImportContext {
	return &ImportContext{
		visited:   make(map[string]bool),
		pathStack: make([]string, 0),
	}
}

// detectCycle checks if path is in the current import chain (cycle detection)
func detectCycle(path string, ctx *ImportContext) error {
	// Check if path is currently in the import chain
	for _, p := range ctx.pathStack {
		if p == path {
			// Cycle detected - build error message showing full cycle
			return buildCycleError(ctx.pathStack, path)
		}
	}

	return nil
}

// buildCycleError creates a descriptive error message showing the import cycle
func buildCycleError(stack []string, newPath string) error {
	// Build cycle path: stack + newPath
	cyclePath := append(append([]string{}, stack...), newPath)
	cycleStr := strings.Join(cyclePath, " → ")

	return fmt.Errorf("import cycle: %s", cycleStr)
}

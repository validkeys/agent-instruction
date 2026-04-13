package rules

import (
	"fmt"
)

// ConfigServiceInterface defines the minimal interface needed for import resolution
type ConfigServiceInterface interface {
	LoadRuleFile(path string) (*RuleFile, error)
}

// Resolver handles import resolution with cycle detection
type Resolver struct {
	configService ConfigServiceInterface
}

// NewResolver creates a new resolver instance
func NewResolver(configService ConfigServiceInterface) *Resolver {
	return &Resolver{
		configService: configService,
	}
}

// ResolveImports resolves all imports starting from rootFile using depth-first traversal
func (r *Resolver) ResolveImports(rootFile string) ([]Instruction, error) {
	ctx := NewImportContext()
	return r.resolveRecursive(rootFile, ctx)
}

// resolveRecursive recursively resolves imports using depth-first traversal
func (r *Resolver) resolveRecursive(filePath string, ctx *ImportContext) ([]Instruction, error) {
	// Normalize path to absolute
	absPath, err := ResolvePath(filePath, ".")
	if err != nil {
		return nil, fmt.Errorf("resolve path %s: %w", filePath, err)
	}

	// Check for cycle - must happen before visited check
	if err := detectCycle(absPath, ctx); err != nil {
		return nil, err
	}

	// If already visited (but not in current path), skip to avoid duplicate processing
	if ctx.visited[absPath] {
		return []Instruction{}, nil
	}

	// Mark as visited and add to path stack
	ctx.visited[absPath] = true
	ctx.pathStack = append(ctx.pathStack, absPath)

	// Remove from path stack when we exit this function (backtrack)
	defer func() {
		ctx.pathStack = ctx.pathStack[:len(ctx.pathStack)-1]
	}()

	// Load rule file
	rule, err := r.configService.LoadRuleFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("load rule file %s: %w", absPath, err)
	}

	// Collect all instructions from imports first (depth-first)
	var allInstructions []Instruction

	for _, importPath := range rule.Imports {
		// Resolve import path relative to current file
		resolvedPath, err := ResolveImportPath(importPath, absPath)
		if err != nil {
			return nil, fmt.Errorf("resolve import %s in %s: %w", importPath, absPath, err)
		}

		// Recursively resolve imports
		importedInstructions, err := r.resolveRecursive(resolvedPath, ctx)
		if err != nil {
			return nil, err
		}

		// Collect imported instructions
		allInstructions = append(allInstructions, importedInstructions...)
	}

	// Add local instructions after all imports
	allInstructions = append(allInstructions, rule.Instructions...)

	return allInstructions, nil
}

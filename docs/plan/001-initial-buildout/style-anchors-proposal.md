# Style Anchors Proposal for agent-instruction

**Generated:** 2026-04-13
**Purpose:** Define concrete code examples that guide AI-assisted development
**Based on:** Cobra repository analysis and project requirements

## Overview

Style anchors are reference documents containing real code examples that show AI agents exactly how to write code for this project. They prevent common mistakes and ensure consistency across all generated code.

## Proposed Style Anchors

### 1. cobra-command-structure.md

**Purpose:** Command definition patterns for CLI commands

**Content to include:**
- Complete example of a Cobra command with Use, Short, Long, Example
- Args validation patterns (NoArgs, ExactArgs, MinimumNArgs)
- RunE function signature and error handling
- Flag definition and binding
- Command registration with AddCommand

**Source reference:**
- `/Users/kydavis/Sites/ai-use-repos/cobra/command.go` (lines 54-146)
- `/Users/kydavis/Sites/ai-use-repos/cobra/args.go`

**Use cases:**
- M3: Init command implementation
- M4: Build command implementation
- M5: Add and List commands

**Example structure:**
```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize agent-instruction in repository",
    Long:  `Creates .agent-instruction directory structure...`,
    Example: "  agent-instruction init\n  agent-instruction init --non-interactive",
    Args:  cobra.NoArgs,
    RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
    // Command implementation
    // Return errors, don't call os.Exit
    // Use cmd.OutOrStdout() for output
}
```

---

### 2. table-driven-testing.md

**Purpose:** Testing patterns for all test files

**Content to include:**
- Table-driven test structure with map[string]struct
- Helper functions for test execution
- Test naming conventions
- Assertion patterns
- Test setup and teardown

**Source reference:**
- `/Users/kydavis/Sites/ai-use-repos/cobra/args_test.go` (lines 440-486)
- `/Users/kydavis/Sites/ai-use-repos/cobra/command_test.go` (lines 30-57, 88-111)

**Use cases:**
- All milestones with _test.go files
- M7: Comprehensive test coverage

**Example structure:**
```go
func TestBuildCommand(t *testing.T) {
    tests := map[string]struct {
        setup    func() string // returns temp dir
        args     []string
        wantErr  bool
        validate func(t *testing.T, dir string)
    }{
        "happy path with global rules": {
            setup: func() string {
                dir := setupTestRepo(t)
                createConfig(t, dir, Config{...})
                return dir
            },
            args:    []string{},
            wantErr: false,
            validate: func(t *testing.T, dir string) {
                assertFileExists(t, filepath.Join(dir, "CLAUDE.md"))
            },
        },
        "error when not initialized": {
            setup:   func() string { return t.TempDir() },
            args:    []string{},
            wantErr: true,
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            dir := tc.setup()
            defer os.RemoveAll(dir)

            output, err := executeCommand(buildCmd, tc.args...)

            if tc.wantErr && err == nil {
                t.Fatal("expected error, got nil")
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if tc.validate != nil {
                tc.validate(t, dir)
            }
        })
    }
}

// Helper functions
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
    buf := new(bytes.Buffer)
    cmd.SetOut(buf)
    cmd.SetErr(buf)
    cmd.SetArgs(args)
    err := cmd.Execute()
    return buf.String(), err
}
```

---

### 3. file-operations.md

**Purpose:** Safe file handling with atomic writes and backups

**Content to include:**
- Atomic write pattern (write to temp, rename)
- Backup creation before modification
- Managed section parsing and replacement
- File permission preservation
- Error handling for file operations

**Source reference:**
- Go standard library patterns (os, io/ioutil, filepath)
- Project requirements (technical-requirements.yaml lines 161-169)

**Use cases:**
- M1: File service implementation
- M3: Init command backup logic
- M4: Build command file generation

**Example structure:**
```go
// AtomicWrite writes content to a file atomically
func AtomicWrite(path string, content []byte, perm os.FileMode) error {
    // Write to temp file first
    tempPath := path + ".tmp"
    if err := os.WriteFile(tempPath, content, perm); err != nil {
        return fmt.Errorf("write temp file: %w", err)
    }

    // Atomic rename
    if err := os.Rename(tempPath, path); err != nil {
        os.Remove(tempPath) // cleanup on failure
        return fmt.Errorf("rename temp file: %w", err)
    }

    return nil
}

// CreateBackup creates a backup of the file with .backup extension
func CreateBackup(path string) error {
    content, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // no file to backup
        }
        return fmt.Errorf("read file for backup: %w", err)
    }

    backupPath := path + ".backup"
    if err := os.WriteFile(backupPath, content, 0644); err != nil {
        return fmt.Errorf("write backup: %w", err)
    }

    return nil
}

// ReplaceManagedSection finds and replaces content between markers
func ReplaceManagedSection(content string, newSection string) (string, error) {
    const beginMarker = "<!-- BEGIN AGENT-INSTRUCTION -->"
    const endMarker = "<!-- END AGENT-INSTRUCTION -->"

    beginIdx := strings.Index(content, beginMarker)
    endIdx := strings.Index(content, endMarker)

    if beginIdx == -1 && endIdx == -1 {
        // No managed section, append to end
        return content + "\n" + beginMarker + "\n" + newSection + "\n" + endMarker + "\n", nil
    }

    if beginIdx == -1 || endIdx == -1 {
        return "", fmt.Errorf("malformed managed section: only one marker found")
    }

    if beginIdx > endIdx {
        return "", fmt.Errorf("malformed managed section: end marker before begin marker")
    }

    // Replace content between markers
    before := content[:beginIdx+len(beginMarker)]
    after := content[endIdx:]
    return before + "\n" + newSection + "\n" + after, nil
}
```

---

### 4. json-config-handling.md

**Purpose:** JSON configuration parsing and validation

**Content to include:**
- Struct definitions with json tags
- Marshal/Unmarshal patterns
- Validation functions
- Nested struct composition
- Optional fields with pointers

**Source reference:**
- Go encoding/json standard library
- Project data model (technical-requirements.yaml lines 101-143)

**Use cases:**
- M0: Core data structures
- M1: Config service
- M2: Rule service

**Example structure:**
```go
// Config represents .agent-instruction/config.json
type Config struct {
    Version    string   `json:"version"`
    Packages   []string `json:"packages"`
    Frameworks []string `json:"frameworks"`
}

// RuleFile represents a rule file (.agent-instruction/rules/*.json)
type RuleFile struct {
    Title        string        `json:"title"`
    Instructions []Instruction `json:"instructions"`
    Imports      []string      `json:"imports,omitempty"`
}

// Instruction represents a single instruction rule
type Instruction struct {
    Heading    string      `json:"heading,omitempty"`
    Rule       string      `json:"rule"`
    References []Reference `json:"references,omitempty"`
}

// Reference represents a reference to another file or section
type Reference struct {
    Title string `json:"title"`
    Path  string `json:"path"`
}

// LoadConfig reads and validates config.json
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("validate config: %w", err)
    }

    return &cfg, nil
}

// validateConfig checks config for required fields and valid values
func validateConfig(cfg *Config) error {
    if cfg.Version == "" {
        return fmt.Errorf("version is required")
    }

    if len(cfg.Frameworks) == 0 {
        return fmt.Errorf("at least one framework is required")
    }

    for _, fw := range cfg.Frameworks {
        if fw != "claude" && fw != "agents" {
            return fmt.Errorf("invalid framework: %s (must be 'claude' or 'agents')", fw)
        }
    }

    return nil
}

// SaveConfig writes config to disk with proper formatting
func SaveConfig(path string, cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal config: %w", err)
    }

    data = append(data, '\n') // add trailing newline

    if err := AtomicWrite(path, data, 0644); err != nil {
        return fmt.Errorf("write config: %w", err)
    }

    return nil
}
```

---

### 5. error-handling.md

**Purpose:** Consistent error handling and messaging

**Content to include:**
- fmt.Errorf with %w for error wrapping
- Clear, actionable error messages
- Error context accumulation
- CLI-appropriate error formatting
- No panic in normal code paths

**Source reference:**
- `/Users/kydavis/Sites/ai-use-repos/cobra/args.go` (error message patterns)
- Go error handling best practices

**Use cases:**
- All milestones
- NFR: Actionable error messages

**Example structure:**
```go
// Good: Clear, actionable error messages with context
func buildInstructionFiles(cfg *Config) error {
    if err := validatePackages(cfg.Packages); err != nil {
        return fmt.Errorf("invalid package configuration: %w", err)
    }

    for _, pkg := range cfg.Packages {
        if err := buildPackageFiles(pkg, cfg); err != nil {
            return fmt.Errorf("build files for package %s: %w", pkg, err)
        }
    }

    return nil
}

// Good: Validation errors suggest solutions
func validateRuleFile(path string, rule *RuleFile) error {
    if rule.Title == "" {
        return fmt.Errorf("rule file %s: title is required", path)
    }

    if len(rule.Instructions) == 0 {
        return fmt.Errorf("rule file %s: must contain at least one instruction", path)
    }

    for i, instr := range rule.Instructions {
        if instr.Rule == "" {
            return fmt.Errorf("rule file %s: instruction %d: rule text is required", path, i)
        }
    }

    return nil
}

// Good: Error wrapping preserves context
func loadRuleFile(path string) (*RuleFile, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("rule file not found: %s (check file path and spelling)", path)
        }
        return nil, fmt.Errorf("read rule file %s: %w", path, err)
    }

    var rule RuleFile
    if err := json.Unmarshal(data, &rule); err != nil {
        return nil, fmt.Errorf("parse rule file %s: %w (check JSON syntax)", path, err)
    }

    return &rule, nil
}

// Bad: Don't use panic for normal errors
func badExample(path string) {
    if _, err := os.Stat(path); err != nil {
        panic(err) // DON'T DO THIS
    }
}

// Good: Return errors for caller to handle
func goodExample(path string) error {
    if _, err := os.Stat(path); err != nil {
        return fmt.Errorf("check path: %w", err)
    }
    return nil
}
```

---

### 6. graph-traversal-cycle-detection.md

**Purpose:** Import resolution algorithm with cycle detection

**Content to include:**
- Depth-first traversal pattern
- Cycle detection with path stack
- Visited set tracking
- Result collection during traversal
- Clear error messages for cycles

**Source reference:**
- Standard CS algorithms
- Project requirements (technical-requirements.yaml lines 442-474)

**Use cases:**
- M2: Rule management and import resolution

**Example structure:**
```go
// ImportResolver handles rule import resolution with cycle detection
type ImportResolver struct {
    visited  map[string]bool
    pathStack []string
    baseDir  string
}

// NewImportResolver creates a new resolver
func NewImportResolver(baseDir string) *ImportResolver {
    return &ImportResolver{
        visited:   make(map[string]bool),
        pathStack: make([]string, 0),
        baseDir:   baseDir,
    }
}

// Resolve recursively resolves imports and returns merged instructions
func (r *ImportResolver) Resolve(rulePath string) ([]Instruction, error) {
    // Make path absolute
    absPath := r.resolveRelativePath(rulePath)

    // Check for cycle
    if r.inPathStack(absPath) {
        return nil, fmt.Errorf("import cycle detected: %s", r.formatCyclePath(absPath))
    }

    // Skip if already visited (but not a cycle)
    if r.visited[absPath] {
        return nil, nil
    }

    // Mark as visited and add to path stack
    r.visited[absPath] = true
    r.pathStack = append(r.pathStack, absPath)
    defer func() {
        // Remove from path stack when done
        r.pathStack = r.pathStack[:len(r.pathStack)-1]
    }()

    // Load rule file
    rule, err := loadRuleFile(absPath)
    if err != nil {
        return nil, err
    }

    // Collect all instructions (imports first, then local)
    var allInstructions []Instruction

    // Recursively resolve imports (depth-first)
    for _, importPath := range rule.Imports {
        importedInstructions, err := r.Resolve(importPath)
        if err != nil {
            return nil, err
        }
        allInstructions = append(allInstructions, importedInstructions...)
    }

    // Add local instructions after imports
    allInstructions = append(allInstructions, rule.Instructions...)

    return allInstructions, nil
}

// inPathStack checks if path is currently in the traversal path
func (r *ImportResolver) inPathStack(path string) bool {
    for _, p := range r.pathStack {
        if p == path {
            return true
        }
    }
    return false
}

// formatCyclePath creates a readable cycle error message
func (r *ImportResolver) formatCyclePath(cyclePath string) string {
    path := append(r.pathStack, cyclePath)
    return strings.Join(path, " → ")
}

// resolveRelativePath converts relative path to absolute based on baseDir
func (r *ImportResolver) resolveRelativePath(path string) string {
    if filepath.IsAbs(path) {
        return path
    }
    return filepath.Join(r.baseDir, path)
}
```

---

## Implementation Guidelines

### For AI Agents

When generating code for agent-instruction:

1. **Always reference the appropriate style anchor** before writing code
2. **Match the patterns exactly** - don't improvise variations
3. **Copy error message style** from error-handling.md
4. **Use table-driven tests** from table-driven-testing.md for all tests
5. **Follow file operation patterns** from file-operations.md for any file I/O

### For Developers

When reviewing AI-generated code:

1. Check that it matches style anchor patterns
2. Verify test structure follows table-driven-testing.md
3. Confirm error messages are actionable (error-handling.md)
4. Ensure file operations are atomic (file-operations.md)

### Creating Style Anchor Files

Each style anchor should be created as a markdown file in:
```
/Users/kydavis/Sites/agent-instruction/docs/style-anchors/
```

Format:
- Clear heading explaining the pattern
- Complete, working code examples
- Annotations explaining key decisions
- "Good" vs "Bad" examples where appropriate
- References to source material

## Next Steps

1. **Create style anchor files** in `docs/style-anchors/` directory
2. **Reference in technical-requirements.yaml** (already listed at line 42)
3. **Include in CLAUDE.md** so AI agents automatically see them
4. **Update as patterns evolve** during development
5. **Add to milestone task descriptions** so agents know which anchors to use

## Benefits

✅ **Consistency** - All code follows the same patterns
✅ **Quality** - Proven patterns from production code (Cobra)
✅ **Speed** - AI doesn't need to guess, just follow examples
✅ **Maintainability** - Single source of truth for coding standards
✅ **Onboarding** - New developers/AI agents learn by example

## Validation

Style anchors are working if:
- AI-generated code matches examples without correction
- Test files all follow table-driven pattern
- Error messages are consistently formatted
- File operations are atomic and safe
- No repeated mistakes across milestones

# Error Handling

**Purpose:** Consistent error handling and messaging
**Source:** Go error handling best practices and Cobra patterns
**Use cases:** All functions that can fail, especially CLI commands

---

## Error Wrapping with %w

✅ **Good: Use fmt.Errorf with %w to wrap errors**

```go
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
```

❌ **Bad: Don't lose error context**

```go
func buildInstructionFiles(cfg *Config) error {
    if err := validatePackages(cfg.Packages); err != nil {
        return err // Lost context: which operation failed?
    }
    // ...
}
```

❌ **Bad: Don't use %v for wrapping**

```go
func buildInstructionFiles(cfg *Config) error {
    if err := validatePackages(cfg.Packages); err != nil {
        return fmt.Errorf("validation failed: %v", err) // %v doesn't preserve error chain
    }
    // ...
}
```

---

## Actionable Error Messages

✅ **Good: Tell user what went wrong AND what to do**

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("config not found at %s: run 'agent-instruction init' to create it", path)
        }
        if os.IsPermission(err) {
            return nil, fmt.Errorf("permission denied reading %s: check file permissions with 'ls -l'", path)
        }
        return nil, fmt.Errorf("read config: %w", err)
    }

    // ... rest of function
}
```

❌ **Bad: Generic error with no guidance**

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err) // What should user do?
    }
    // ...
}
```

---

## Validation Error Messages

✅ **Good: Specific field-level errors with suggestions**

```go
func validateRuleFile(path string, rule *RuleFile) error {
    if rule.Title == "" {
        return fmt.Errorf("rule file %s: title is required (add \"title\": \"Rule Name\" to JSON)", path)
    }

    if len(rule.Instructions) == 0 {
        return fmt.Errorf("rule file %s: must contain at least one instruction (add \"instructions\" array)", path)
    }

    for i, instr := range rule.Instructions {
        if instr.Rule == "" {
            return fmt.Errorf("rule file %s: instruction %d: rule text is required (\"rule\" field cannot be empty)", path, i)
        }
    }

    return nil
}
```

❌ **Bad: Vague validation errors**

```go
func validateRuleFile(path string, rule *RuleFile) error {
    if rule.Title == "" {
        return fmt.Errorf("invalid rule file") // Which file? What's wrong?
    }
    // ...
}
```

---

## Error Context Accumulation

```go
// Each layer adds context as error bubbles up
func ProcessPackage(pkgName string) error {
    if err := loadPackageConfig(pkgName); err != nil {
        return fmt.Errorf("load config for package %s: %w", pkgName, err)
    }

    if err := buildPackageFiles(pkgName); err != nil {
        return fmt.Errorf("build files for package %s: %w", pkgName, err)
    }

    return nil
}

func loadPackageConfig(pkgName string) error {
    path := filepath.Join(pkgName, "config.json")
    if _, err := os.ReadFile(path); err != nil {
        return fmt.Errorf("read config file %s: %w", path, err)
    }
    return nil
}

// Final error might look like:
// "load config for package api: read config file api/config.json: no such file or directory"
```

---

## No Panic in Normal Code

✅ **Good: Return errors for caller to handle**

```go
func initialize(dir string) error {
    if _, err := os.Stat(dir); err != nil {
        return fmt.Errorf("check directory: %w", err)
    }

    if err := createStructure(dir); err != nil {
        return fmt.Errorf("create structure: %w", err)
    }

    return nil
}
```

❌ **Bad: Using panic for normal error conditions**

```go
func initialize(dir string) {
    if _, err := os.Stat(dir); err != nil {
        panic(fmt.Sprintf("directory check failed: %v", err)) // DON'T DO THIS
    }
    // ...
}
```

**Note:** Panic is acceptable for:
- Programmer errors (e.g., nil pointer bugs)
- Truly unrecoverable situations
- Init-time configuration errors (before main() runs)

---

## CLI-Appropriate Error Formatting

```go
// Command error handler in main.go
func main() {
    rootCmd := cmd.NewRootCommand()

    if err := rootCmd.Execute(); err != nil {
        // Cobra already prints error, but we can customize
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

// In command RunE functions
func runBuild(cmd *cobra.Command, args []string) error {
    cfg, err := loadConfig()
    if err != nil {
        // Return error directly - Cobra will format it
        return err
    }

    if err := buildFiles(cfg); err != nil {
        // Add helpful context
        return fmt.Errorf("build failed: %w\n\nTip: Check your config.json and rule files for syntax errors", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "✓ Build complete\n")
    return nil
}
```

---

## Multi-Error Collection

```go
// When you need to collect multiple errors
type ErrorCollector struct {
    errors []error
}

func (ec *ErrorCollector) Add(err error) {
    if err != nil {
        ec.errors = append(ec.errors, err)
    }
}

func (ec *ErrorCollector) Error() error {
    if len(ec.errors) == 0 {
        return nil
    }

    if len(ec.errors) == 1 {
        return ec.errors[0]
    }

    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("multiple errors (%d):\n", len(ec.errors)))
    for i, err := range ec.errors {
        sb.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
    }
    return fmt.Errorf("%s", sb.String())
}

// Usage
func validateAllRules(rules []*RuleFile) error {
    ec := &ErrorCollector{}

    for _, rule := range rules {
        ec.Add(validateRuleFile(rule.Path, rule))
    }

    return ec.Error()
}
```

---

## Sentinel Errors for Special Cases

```go
// Define sentinel errors for specific conditions
var (
    ErrNotInitialized = errors.New("agent-instruction not initialized: run 'agent-instruction init'")
    ErrImportCycle    = errors.New("import cycle detected")
    ErrInvalidConfig  = errors.New("invalid configuration")
)

// Check for specific errors
func processConfig() error {
    cfg, err := loadConfig()
    if err != nil {
        if errors.Is(err, ErrNotInitialized) {
            // Handle not initialized case
            return fmt.Errorf("%w in current directory", ErrNotInitialized)
        }
        return err
    }
    // ...
}
```

---

## Error Checking Patterns

```go
// Pattern: Check error type
func handleError(err error) {
    if os.IsNotExist(err) {
        fmt.Println("File not found")
    } else if os.IsPermission(err) {
        fmt.Println("Permission denied")
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
}

// Pattern: Unwrap to check root cause
func isTimeoutError(err error) bool {
    for err != nil {
        if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
            return true
        }
        err = errors.Unwrap(err)
    }
    return false
}
```

---

## Testing Error Messages

```go
func TestLoadConfigErrors(t *testing.T) {
    tests := map[string]struct {
        setup       func() string
        wantErrMsg  string
        wantErrType error
    }{
        "returns clear error when file not found": {
            setup: func() string {
                return filepath.Join(t.TempDir(), "nonexistent.json")
            },
            wantErrMsg: "run 'agent-instruction init'",
        },
        "returns clear error for invalid JSON": {
            setup: func() string {
                dir := t.TempDir()
                path := filepath.Join(dir, "config.json")
                os.WriteFile(path, []byte("{ invalid json"), 0644)
                return path
            },
            wantErrMsg: "parse JSON",
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            path := tc.setup()

            _, err := LoadConfig(path)

            if err == nil {
                t.Fatal("expected error, got nil")
            }

            if !strings.Contains(err.Error(), tc.wantErrMsg) {
                t.Errorf("error message %q does not contain %q", err.Error(), tc.wantErrMsg)
            }
        })
    }
}
```

---

## Complete Example: Error Handling Flow

```go
// Top-level command function
func runInit(cmd *cobra.Command, args []string) error {
    dir, _ := cmd.Flags().GetString("dir")

    // Validate input
    if err := validateDirectory(dir); err != nil {
        return fmt.Errorf("invalid directory: %w", err)
    }

    // Check existing state
    if err := checkNotAlreadyInitialized(dir); err != nil {
        return err // Already has good context
    }

    // Perform operations with rollback on error
    if err := initializeWithRollback(dir); err != nil {
        return fmt.Errorf("initialization failed: %w", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "✓ Successfully initialized in %s\n", dir)
    return nil
}

func validateDirectory(dir string) error {
    info, err := os.Stat(dir)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("directory does not exist: %s", dir)
        }
        return fmt.Errorf("stat directory: %w", err)
    }

    if !info.IsDir() {
        return fmt.Errorf("path is not a directory: %s", dir)
    }

    return nil
}

func checkNotAlreadyInitialized(dir string) error {
    configPath := filepath.Join(dir, ".agent-instruction", "config.json")
    if _, err := os.Stat(configPath); err == nil {
        return fmt.Errorf("already initialized at %s: use 'agent-instruction build' to regenerate files", dir)
    }
    return nil
}

func initializeWithRollback(dir string) error {
    agentDir := filepath.Join(dir, ".agent-instruction")

    // Create directory
    if err := os.MkdirAll(agentDir, 0755); err != nil {
        return fmt.Errorf("create .agent-instruction directory: %w", err)
    }

    // If any step fails, clean up
    defer func() {
        if r := recover(); r != nil {
            os.RemoveAll(agentDir)
            panic(r)
        }
    }()

    // Create config
    if err := createDefaultConfig(agentDir); err != nil {
        os.RemoveAll(agentDir) // Cleanup on error
        return fmt.Errorf("create config: %w", err)
    }

    return nil
}
```

---

## Key Principles

1. **Always use %w for error wrapping** - Preserves error chain for errors.Is/As
2. **Add context at each layer** - Build up helpful error messages
3. **Be specific about what failed** - Include file paths, field names, indices
4. **Suggest solutions** - Tell user how to fix the problem
5. **Never panic in normal code** - Return errors instead
6. **Test error messages** - Verify errors are helpful

---

## References

- Go error handling: https://go.dev/blog/error-handling-and-go
- Working with errors: https://go.dev/blog/go1.13-errors
- Error wrapping: https://pkg.go.dev/errors
- Source: `/Users/kydavis/Sites/ai-use-repos/cobra/args.go` (error message patterns)

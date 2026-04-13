# Table-Driven Testing

**Purpose:** Testing patterns for all test files
**Source:** Cobra test suite and Go best practices
**Use cases:** All `_test.go` files in the project

---

## Complete Test Example

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
                createConfig(t, dir, Config{
                    Version:    "1.0",
                    Frameworks: []string{"claude"},
                    Packages:   []string{},
                })
                createRuleFile(t, dir, "global.json", &RuleFile{
                    Title: "Global Rules",
                    Instructions: []Instruction{
                        {Rule: "Always use error wrapping"},
                    },
                })
                return dir
            },
            args:    []string{},
            wantErr: false,
            validate: func(t *testing.T, dir string) {
                assertFileExists(t, filepath.Join(dir, "CLAUDE.md"))
                content := readFile(t, filepath.Join(dir, "CLAUDE.md"))
                if !strings.Contains(content, "Always use error wrapping") {
                    t.Error("expected rule in output")
                }
            },
        },
        "error when not initialized": {
            setup: func() string {
                return t.TempDir() // Empty directory
            },
            args:    []string{},
            wantErr: true,
        },
        "handles missing config file": {
            setup: func() string {
                dir := t.TempDir()
                // Create .agent-instruction dir but no config
                os.MkdirAll(filepath.Join(dir, ".agent-instruction"), 0755)
                return dir
            },
            args:    []string{},
            wantErr: true,
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            dir := tc.setup()
            defer os.RemoveAll(dir)

            // Change to test directory
            oldWd, _ := os.Getwd()
            os.Chdir(dir)
            defer os.Chdir(oldWd)

            output, err := executeCommand(buildCmd, tc.args...)

            if tc.wantErr && err == nil {
                t.Fatal("expected error, got nil")
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("unexpected error: %v\nOutput: %s", err, output)
            }

            if tc.validate != nil {
                tc.validate(t, dir)
            }
        })
    }
}
```

---

## Test Structure Pattern

✅ **Good: Use map[string]struct for test cases**

```go
func TestFunction(t *testing.T) {
    tests := map[string]struct {
        input    string
        expected string
        wantErr  bool
    }{
        "descriptive test name": {
            input:    "test input",
            expected: "expected output",
            wantErr:  false,
        },
        "handles error case": {
            input:   "bad input",
            wantErr: true,
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            result, err := Function(tc.input)

            if tc.wantErr && err == nil {
                t.Fatal("expected error, got nil")
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if !tc.wantErr && result != tc.expected {
                t.Errorf("got %q, want %q", result, tc.expected)
            }
        })
    }
}
```

❌ **Bad: Avoid slice-based table tests (less clear)**

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string // Have to add name field manually
        input    string
        expected string
    }{
        {"test 1", "input", "output"}, // Name is separate from structure
    }
    // ... rest of test
}
```

---

## Helper Functions

```go
// executeCommand runs a Cobra command and captures output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
    buf := new(bytes.Buffer)
    cmd.SetOut(buf)
    cmd.SetErr(buf)
    cmd.SetArgs(args)

    err := cmd.Execute()
    return buf.String(), err
}

// setupTestRepo creates a temporary repository structure
func setupTestRepo(t *testing.T) string {
    t.Helper()

    dir := t.TempDir()
    agentDir := filepath.Join(dir, ".agent-instruction")

    if err := os.MkdirAll(filepath.Join(agentDir, "rules"), 0755); err != nil {
        t.Fatalf("create test repo: %v", err)
    }

    return dir
}

// createConfig writes a config.json file for testing
func createConfig(t *testing.T, baseDir string, cfg Config) {
    t.Helper()

    path := filepath.Join(baseDir, ".agent-instruction", "config.json")
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        t.Fatalf("marshal config: %v", err)
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        t.Fatalf("write config: %v", err)
    }
}

// createRuleFile writes a rule file for testing
func createRuleFile(t *testing.T, baseDir, filename string, rule *RuleFile) {
    t.Helper()

    path := filepath.Join(baseDir, ".agent-instruction", "rules", filename)
    data, err := json.MarshalIndent(rule, "", "  ")
    if err != nil {
        t.Fatalf("marshal rule: %v", err)
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        t.Fatalf("write rule: %v", err)
    }
}

// assertFileExists fails if file doesn't exist
func assertFileExists(t *testing.T, path string) {
    t.Helper()

    if _, err := os.Stat(path); os.IsNotExist(err) {
        t.Fatalf("expected file to exist: %s", path)
    }
}

// readFile reads file content and fails on error
func readFile(t *testing.T, path string) string {
    t.Helper()

    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("read file %s: %v", path, err)
    }
    return string(data)
}
```

---

## Test Naming Conventions

✅ **Good: Descriptive, behavior-focused names**

```go
tests := map[string]struct{
    // ...
}{
    "creates config with default frameworks":     {...},
    "returns error when directory already exists": {...},
    "handles import cycle detection":             {...},
    "merges instructions from multiple rules":    {...},
}
```

❌ **Bad: Generic or implementation-focused names**

```go
tests := map[string]struct{
    // ...
}{
    "test 1":        {...}, // Not descriptive
    "edge case":     {...}, // What edge case?
    "calls loadRuleFile": {...}, // Tests behavior, not implementation
}
```

---

## Assertion Patterns

```go
// Error checking
if tc.wantErr && err == nil {
    t.Fatal("expected error, got nil")
}
if !tc.wantErr && err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Value comparison
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// String contains
if !strings.Contains(output, expected) {
    t.Errorf("output missing expected text: %q\nGot: %s", expected, output)
}

// Slice comparison
if len(got) != len(want) {
    t.Errorf("got %d items, want %d", len(got), len(want))
}

// Custom validation
if tc.validate != nil {
    tc.validate(t, result)
}
```

---

## Setup and Teardown

```go
func TestWithSetupTeardown(t *testing.T) {
    tests := map[string]struct {
        setup    func() string
        teardown func(dir string)
        // ... other fields
    }{
        "test case": {
            setup: func() string {
                dir := t.TempDir()
                // Setup logic
                return dir
            },
            teardown: func(dir string) {
                // Custom cleanup (if needed beyond TempDir)
                os.RemoveAll(filepath.Join(dir, "backup"))
            },
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            dir := tc.setup()
            if tc.teardown != nil {
                defer tc.teardown(dir)
            }

            // Test execution
        })
    }
}
```

---

## Testing Commands with Flags

```go
func TestCommandFlags(t *testing.T) {
    tests := map[string]struct {
        args     []string
        validate func(t *testing.T, output string)
    }{
        "with package flag": {
            args: []string{"--package", "api"},
            validate: func(t *testing.T, output string) {
                if !strings.Contains(output, "Building for package: api") {
                    t.Error("expected package name in output")
                }
            },
        },
        "with dry-run flag": {
            args: []string{"--dry-run"},
            validate: func(t *testing.T, output string) {
                if !strings.Contains(output, "DRY RUN") {
                    t.Error("expected dry-run indicator")
                }
            },
        },
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            output, err := executeCommand(buildCmd, tc.args...)
            if err != nil {
                t.Fatalf("command failed: %v", err)
            }
            tc.validate(t, output)
        })
    }
}
```

---

## Key Principles

1. **Use `map[string]struct` for test cases** - Name is part of the key
2. **Always call `t.Helper()` in helper functions** - Better error line numbers
3. **Use `t.TempDir()` for temporary files** - Automatic cleanup
4. **Descriptive test names describe behavior** - "creates X when Y" not "test 1"
5. **Validate both success and failure cases** - Test error paths thoroughly
6. **Keep test cases independent** - Each test should run in isolation

---

## References

- Go testing best practices: https://go.dev/doc/tutorial/add-a-test
- Table-driven tests: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
- Source: `/Users/kydavis/Sites/ai-use-repos/cobra/args_test.go` (lines 440-486)
- Source: `/Users/kydavis/Sites/ai-use-repos/cobra/command_test.go` (lines 30-57, 88-111)

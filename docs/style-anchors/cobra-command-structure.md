# Cobra Command Structure

**Purpose:** Command definition patterns for CLI commands
**Source:** Cobra CLI framework (`github.com/spf13/cobra`)
**Use cases:** All command implementations (init, build, add, list)

---

## Complete Command Example

```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize agent-instruction in repository",
    Long:  `Creates .agent-instruction directory structure and configuration files.

This command sets up the initial directory structure for managing AI agent
instructions in your repository. It creates config.json, directories for rules,
and optional backup of existing files.`,
    Example: `  # Initialize with interactive prompts
  agent-instruction init

  # Initialize with defaults (non-interactive)
  agent-instruction init --non-interactive

  # Initialize in a specific directory
  agent-instruction init --dir /path/to/repo`,
    Args:  cobra.NoArgs,
    RunE:  runInit,
}

func init() {
    // Register command with root
    rootCmd.AddCommand(initCmd)

    // Define flags
    initCmd.Flags().Bool("non-interactive", false, "skip interactive prompts")
    initCmd.Flags().StringP("dir", "d", ".", "target directory")
}

func runInit(cmd *cobra.Command, args []string) error {
    // Get flag values
    nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
    targetDir, _ := cmd.Flags().GetString("dir")

    // Business logic implementation
    if err := validateDirectory(targetDir); err != nil {
        return fmt.Errorf("invalid target directory: %w", err)
    }

    if err := createDirectoryStructure(targetDir); err != nil {
        return fmt.Errorf("create directory structure: %w", err)
    }

    // Use cmd.OutOrStdout() for output (testable)
    fmt.Fprintf(cmd.OutOrStdout(), "✓ Initialized agent-instruction in %s\n", targetDir)

    return nil // Return error, don't call os.Exit()
}
```

---

## Args Validation Patterns

```go
// NoArgs - command accepts no arguments
var initCmd = &cobra.Command{
    Use:  "init",
    Args: cobra.NoArgs,
    RunE: runInit,
}

// ExactArgs - require exact number of arguments
var addCmd = &cobra.Command{
    Use:  "add <rule-name>",
    Args: cobra.ExactArgs(1),
    RunE: runAdd,
}

// MinimumNArgs - require at least N arguments
var listCmd = &cobra.Command{
    Use:  "list [packages...]",
    Args: cobra.MinimumNArgs(0),
    RunE: runList,
}

// Custom validation function
var buildCmd = &cobra.Command{
    Use:  "build",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) > 0 {
            return fmt.Errorf("build does not accept arguments, use --package flag instead")
        }
        return nil
    },
    RunE: runBuild,
}
```

---

## RunE Function Signature

✅ **Good: Use RunE for error handling**

```go
func runInit(cmd *cobra.Command, args []string) error {
    // Implementation that can return errors
    if err := doWork(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
```

❌ **Bad: Don't use Run with os.Exit()**

```go
func runInit(cmd *cobra.Command, args []string) {
    if err := doWork(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1) // DON'T DO THIS - not testable
    }
}
```

---

## Flag Definition Patterns

```go
func init() {
    rootCmd.AddCommand(buildCmd)

    // String flag with shorthand
    buildCmd.Flags().StringP("package", "p", "", "target package name")

    // Boolean flag
    buildCmd.Flags().Bool("dry-run", false, "preview changes without writing files")

    // String slice flag (repeatable)
    buildCmd.Flags().StringSlice("frameworks", []string{"claude"}, "target frameworks")

    // Mark flag as required
    buildCmd.MarkFlagRequired("package")
}
```

---

## Output Handling

✅ **Good: Use cmd.OutOrStdout() for testable output**

```go
func runBuild(cmd *cobra.Command, args []string) error {
    result := generateOutput()
    fmt.Fprintf(cmd.OutOrStdout(), "Generated: %s\n", result)
    return nil
}
```

❌ **Bad: Direct stdout/stderr (hard to test)**

```go
func runBuild(cmd *cobra.Command, args []string) error {
    result := generateOutput()
    fmt.Println("Generated:", result) // Hard to capture in tests
    return nil
}
```

---

## Command Registration

```go
// In cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "agent-instruction",
    Short: "Manage AI agent instructions in monorepos",
    Long:  `A CLI tool for managing CLAUDE.md and AGENTS.md files...`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Register subcommands
    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(buildCmd)
    rootCmd.AddCommand(addCmd)
    rootCmd.AddCommand(listCmd)
}
```

---

## Key Principles

1. **Always use `RunE` not `Run`** - Enables proper error handling
2. **Return errors, never call `os.Exit()`** - Makes commands testable
3. **Use `cmd.OutOrStdout()` for output** - Allows output capture in tests
4. **Validate arguments with `Args` field** - Clear error messages for wrong usage
5. **Define flags in `init()` function** - Consistent flag registration
6. **Use clear, actionable help text** - Users should understand command purpose immediately

---

## References

- Cobra documentation: https://github.com/spf13/cobra
- Source: `/Users/kydavis/Sites/ai-use-repos/cobra/command.go` (lines 54-146)
- Args patterns: `/Users/kydavis/Sites/ai-use-repos/cobra/args.go`

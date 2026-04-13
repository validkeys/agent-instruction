# ADR-004: Sentinel Error Pattern

**Status**: Accepted

**Date**: 2026-04-13

**Deciders**: Development Team

---

## Context

The agent-instruction tool performs validation on configuration files, rule files, and user input. When validation fails, we need to:

1. **Communicate specific errors** to users (what went wrong, why)
2. **Enable programmatic error checking** (distinguish error types in code)
3. **Support error wrapping** (add context while preserving original error)
4. **Follow Go idioms** (use standard library patterns)

Go provides several error handling patterns:
- String comparison: `if err.Error() == "some message"`
- Error types: `if _, ok := err.(*ValidationError); ok`
- Sentinel errors: `if errors.Is(err, ErrNotFound)`
- Error wrapping: `fmt.Errorf("context: %w", err)`

The question: Which pattern should we use for validation and domain errors?

## Decision

We will use **sentinel errors with error wrapping** for all validation and domain errors:

```go
// Package-level sentinel errors
var (
    ErrVersionRequired    = errors.New("config version is required")
    ErrFrameworkRequired  = errors.New("at least one framework is required")
    ErrInvalidFramework   = errors.New("invalid framework")
)

// Validation wraps sentinel with context
func (c *Config) Validate() error {
    if c.Version == "" {
        return fmt.Errorf("%w", ErrVersionRequired)
    }

    if len(c.Frameworks) == 0 {
        return fmt.Errorf("%w", ErrFrameworkRequired)
    }

    for _, fw := range c.Frameworks {
        if !validFrameworks[fw] {
            return fmt.Errorf("%w: %s (must be 'claude' or 'agents')",
                ErrInvalidFramework, fw)
        }
    }

    return nil
}

// Callers can check error type
if err := cfg.Validate(); err != nil {
    if errors.Is(err, config.ErrInvalidFramework) {
        // Handle invalid framework specifically
    }
    return err
}
```

This pattern applies to:
- Config validation (`internal/config/types.go`)
- Rule validation (`internal/rules/types.go`)
- Future validation logic

## Alternatives Considered

### Alternative 1: String Error Messages Only
```go
func (c *Config) Validate() error {
    if c.Version == "" {
        return errors.New("config version is required")
    }
    // ...
}
```

**Pros:**
- Simplest approach
- No package-level variables
- Direct error creation

**Cons:**
- No programmatic error checking
- String comparison brittle (`err.Error() == "..."`)
- Cannot distinguish error types
- Wrapping loses error identity

**Rejected because:** Cannot distinguish error types programmatically. Callers must parse error strings.

### Alternative 2: Custom Error Types
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (c *Config) Validate() error {
    if c.Version == "" {
        return &ValidationError{Field: "version", Message: "required"}
    }
    // ...
}
```

**Pros:**
- Rich error context (fields, codes)
- Type assertion possible
- Structured error data

**Cons:**
- More boilerplate per error
- Harder to wrap with fmt.Errorf
- Doesn't work well with errors.Is()
- Over-engineered for simple validation

**Rejected because:** Too complex for validation errors. Sentinel errors are simpler and sufficient.

### Alternative 3: Error Codes
```go
const (
    ErrCodeVersionRequired = 1001
    ErrCodeFrameworkRequired = 1002
)

type CodedError struct {
    Code    int
    Message string
}
```

**Pros:**
- Numeric codes for error types
- Easy to serialize/log
- Can map to HTTP status codes

**Cons:**
- Requires code registry
- Extra maintenance burden
- Doesn't work with errors.Is()
- Not idiomatic Go

**Rejected because:** Adds complexity without benefit. Error codes more suitable for APIs than CLI tools.

### Alternative 4: Wrapped String Errors
```go
func (c *Config) Validate() error {
    if c.Version == "" {
        return fmt.Errorf("validation: config version is required")
    }
    // ...
}
```

**Pros:**
- Context in error message
- Uses fmt.Errorf
- Simple

**Cons:**
- No programmatic checking
- Cannot use errors.Is()
- String comparison only

**Rejected because:** Same issue as Alternative 1 - no type checking.

## Consequences

### Positive

1. **Type-Safe Error Checking**
   ```go
   if errors.Is(err, config.ErrVersionRequired) {
       // Handle specifically
   }
   ```
   Callers can check error types without string comparison.

2. **Error Wrapping Support**
   ```go
   return fmt.Errorf("%w: %s", ErrInvalidFramework, framework)
   ```
   Add context while preserving error identity for `errors.Is()`.

3. **Discoverable**
   - Package-level variables are easy to find
   - Exported errors appear in godoc
   - IDE autocomplete shows available errors

4. **Idiomatic Go**
   - Standard library pattern (e.g., `io.EOF`, `sql.ErrNoRows`)
   - Works with `errors.Is()` and `errors.As()`
   - Familiar to Go developers

5. **Consistent Pattern**
   - Same approach across all packages
   - Easy to extend with new errors
   - Clear naming convention (`Err*`)

### Negative

1. **Package-Level Variables**
   - Errors defined at package scope
   - Could clutter package namespace
   - Mitigation: Group with `var (...)` block

2. **Limited Context**
   - Sentinel error is just a string
   - Must wrap to add context (field values, etc.)
   - Mitigation: Wrapping with `%w` is easy and expected

3. **Not Structured**
   - No error codes, fields, metadata
   - Mitigation: Not needed for CLI tool, string messages sufficient

### Neutral

1. **Naming Convention**
   - All sentinels prefixed with `Err`
   - Clear and consistent
   - Example: `ErrVersionRequired`, `ErrInvalidFramework`

2. **Wrapping is Optional**
   - Can return sentinel directly: `return ErrVersionRequired`
   - Can wrap with context: `return fmt.Errorf("%w: %s", ErrInvalidFramework, fw)`
   - Choice depends on whether additional context is useful

## Implementation Notes

### Defining Sentinel Errors
```go
// Group related errors in var block
var (
    ErrVersionRequired   = errors.New("config version is required")
    ErrFrameworkRequired = errors.New("at least one framework is required")
    ErrInvalidFramework  = errors.New("invalid framework")
)
```

Place at top of type definition file (e.g., `config/types.go`, `rules/types.go`).

### Using Sentinel Errors
```go
// Return directly if no additional context needed
if c.Version == "" {
    return fmt.Errorf("%w", ErrVersionRequired)
}

// Wrap with additional context
for _, fw := range c.Frameworks {
    if !validFrameworks[fw] {
        return fmt.Errorf("%w: %s (must be 'claude' or 'agents')",
            ErrInvalidFramework, fw)
    }
}
```

Always use `%w` (not `%v`) to preserve error for `errors.Is()`.

### Checking Sentinel Errors
```go
// In caller code
if err := cfg.Validate(); err != nil {
    if errors.Is(err, config.ErrInvalidFramework) {
        fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
        fmt.Fprintf(os.Stderr, "Valid frameworks: claude, agents\n")
        return err
    }
    // Generic error handling
    return fmt.Errorf("config validation failed: %w", err)
}
```

Use `errors.Is()` for checking, not type assertion or string comparison.

### Error Messages
Sentinel error messages should be:
- **Actionable**: User knows what to fix
- **Specific**: Clear which field/constraint failed
- **Lowercase**: Start with lowercase (will be wrapped)
- **No punctuation**: No trailing period (will be added when wrapped)

Good examples:
- `"config version is required"`
- `"invalid framework"`
- `"rule text is required"`

Bad examples:
- `"Error: Config version is required."` (capitalized, punctuation)
- `"Something went wrong"` (vague)
- `"Invalid input"` (not specific)

## Testing Strategy

Tests verify:
1. Validation returns correct sentinel error
2. Error wrapping preserves sentinel identity
3. `errors.Is()` works with wrapped errors
4. Error messages are actionable

Example test:
```go
func TestConfigValidation(t *testing.T) {
    cfg := Config{} // Missing version

    err := cfg.Validate()

    if !errors.Is(err, ErrVersionRequired) {
        t.Errorf("expected ErrVersionRequired, got %v", err)
    }

    if !strings.Contains(err.Error(), "version is required") {
        t.Errorf("error message missing context: %v", err)
    }
}
```

## Error Categories

Current sentinel errors by package:

**config/types.go:**
- `ErrVersionRequired` - config.json missing version
- `ErrFrameworkRequired` - config.json missing frameworks
- `ErrInvalidFramework` - unknown framework name

**rules/types.go:**
- `ErrTitleRequired` - rule file missing title
- `ErrInstructionsRequired` - rule file has no instructions
- `ErrRuleTextRequired` - instruction missing rule text

Future errors will follow same pattern.

## References

- [Go blog: Working with Errors](https://go.dev/blog/go1.13-errors)
- [errors.Is() documentation](https://pkg.go.dev/errors#Is)
- [Sentinel errors in Go](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- golang-design-patterns skill - Error handling section

## Related Decisions

- ADR-001: Interface-Based Service Design (services return sentinel errors)
- Error wrapping strategy (using %w throughout codebase)
- Future: User-facing error formatting (CLI output)

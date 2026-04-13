# Go Standards Compliance Tracking

**Status**: ✅ Compliant (High & Medium priority items completed)
**Last Updated**: 2026-04-13
**Skill**: golang-design-patterns

## Overview

This document tracks alignment of the codebase with idiomatic Go design patterns. Items are prioritized by impact and ordered by file.

## Priority Legend

- 🔴 **High**: Affects error handling, API stability, or production readiness
- 🟡 **Medium**: Improves maintainability or testing
- 🟢 **Low**: Nice-to-have improvements

---

## Action Items

### 🔴 High Priority

#### 1. Add error wrapping with `%w` for better error traces

**Files**:
- `internal/config/types.go`
- `internal/rules/types.go`

**Current Pattern**:
```go
return fmt.Errorf("version is required")
return fmt.Errorf("invalid framework: %s", fw)
```

**Required Pattern**:
```go
// Define sentinel errors at package level
var (
    ErrMissingVersion    = errors.New("version is required")
    ErrMissingFrameworks = errors.New("at least one framework is required")
    ErrInvalidFramework  = errors.New("invalid framework")
)

// Use %w in validation
return fmt.Errorf("config validation failed: %w", ErrMissingVersion)
return fmt.Errorf("%w: %s (must be 'claude' or 'agents')", ErrInvalidFramework, fw)
```

**Files to Update**:
- [x] `internal/config/types.go:14` - version validation
- [x] `internal/config/types.go:19` - frameworks validation
- [x] `internal/config/types.go:29` - invalid framework validation
- [x] `internal/rules/types.go:28` - title validation
- [x] `internal/rules/types.go:32` - instructions validation
- [x] `internal/rules/types.go:37` - instruction rule validation

**Benefit**: Enables `errors.Is()` and `errors.As()` for proper error handling up the stack.

---

### 🟡 Medium Priority

#### 2. Remove or create missing test fixtures

**File**: `internal/commands/testhelpers_test.go`

**Issue**: Tests reference fixtures that don't exist:
- Line 216: `../../testdata/config/valid-config.json`
- Line 224: `../../testdata/rules/simple-rule.json`

**Options**:
1. Create the fixtures and testdata directory structure
2. Remove these tests until fixtures are needed
3. Use inline test data instead of fixtures

**Recommendation**: Remove these tests for now (simpler). Add back when you need to test fixture loading.

**Files to Update**:
- [x] `internal/commands/testhelpers_test.go:210-242` - Remove `TestLoadFixtures`

---

#### 3. Call `Validate()` in constructors when added

**Status**: Future consideration

**Pattern**: When you add `New*` constructors, call validation immediately:

```go
func NewConfig(version string, frameworks []string, packages []string) (*Config, error) {
    c := &Config{
        Version:    version,
        Frameworks: frameworks,
        Packages:   packages,
    }
    if err := c.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    return c, nil
}
```

**Benefit**: Fail fast at construction time, not at use time.

**Action**: Keep in mind for future development. No immediate action needed.

---

### 🟢 Low Priority

#### 4. Consider functional options if configuration grows

**File**: `internal/commands/root.go`

**Current**: Simple constructor (appropriate for current needs)

**Future Pattern** (if you add 3+ configuration options):
```go
type RootOption func(*cobra.Command)

func WithDebug(enabled bool) RootOption {
    return func(cmd *cobra.Command) {
        cmd.PersistentFlags().Bool("debug", enabled, "enable debug output")
    }
}

func WithTimeout(d time.Duration) RootOption {
    return func(cmd *cobra.Command) {
        // configure timeout
    }
}

func NewRootCmd(opts ...RootOption) *cobra.Command {
    cmd := &cobra.Command{...}
    for _, opt := range opts {
        opt(cmd)
    }
    return cmd
}
```

**Action**: Monitor. Apply when you have 3+ flags/options.

---

#### 5. Add compile-time interface checks

**Status**: Future consideration

**Pattern**: When you define interfaces (Validator, Loader, etc.):

```go
// Validator validates data structures
type Validator interface {
    Validate() error
}

// Compile-time checks (add at package level)
var (
    _ Validator = (*Config)(nil)
    _ Validator = (*RuleFile)(nil)
)
```

**Action**: Apply when you add interfaces to the codebase.

---

## Compliance Checklist

### Current Strengths ✅

- [x] No `init()` functions
- [x] Proper error handling flow (early returns, flat happy path)
- [x] Table-driven tests throughout
- [x] Proper use of `t.Helper()` in test utilities
- [x] No global mutable state
- [x] Simple, clear constructors

### Standards to Apply

- [x] Error wrapping with `%w` and sentinel errors (🔴 High)
- [x] Clean test fixtures or remove references (🟡 Medium)
- [ ] Constructor validation pattern (🟡 Medium - future)
- [ ] Functional options pattern (🟢 Low - when needed)
- [ ] Compile-time interface checks (🟢 Low - when interfaces added)

---

## Next Steps

1. **Address High Priority Items First**
   - Add sentinel errors and `%w` wrapping to validation methods
   - Update all error returns in `config` and `rules` packages

2. **Clean Up Tests**
   - Remove or implement fixture loading tests

3. **Continue Development**
   - Apply patterns as needed when adding new features
   - Reference this doc during code reviews

---

## References

- Go Design Patterns Skill: `golang-design-patterns`
- Error Handling: See `samber/cc-skills-golang@golang-error-handling`
- Project Layout: See `samber/cc-skills-golang@golang-project-layout`

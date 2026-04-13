# ADR-001: Interface-Based Service Design

**Status**: Accepted

**Date**: 2026-04-13

**Deciders**: Development Team

---

## Context

As we build the agent-instruction CLI tool for managing CLAUDE.md and AGENTS.md files, we need to decide how to structure our service layer. The tool needs to perform file operations, configuration management, and rule processing. We need an architecture that is:

1. **Testable** - Easy to write unit tests without requiring actual file system operations
2. **Mockable** - Ability to substitute implementations for testing
3. **Extensible** - New implementations can be added without changing callers
4. **Maintainable** - Clear contracts between components

The primary question is: Should we use concrete types directly, or abstract them behind interfaces?

## Decision

We will use **interface-based service design** for all major service layers:

```go
// Define interface
type FileService interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, content []byte) error
    BackupFile(path string) error
    ParseManaged(content []byte) (*ManagedContent, error)
    UpdateManaged(path string, newContent string) error
}

// Provide default implementation
type DefaultFileService struct{}

// Constructor establishes pattern for future extensibility
func NewFileService() *DefaultFileService {
    return &DefaultFileService{}
}
```

This pattern applies to:
- `FileService` - File operations, backups, managed sections
- `ConfigService` - Configuration and rule file management
- Future services as they are added

## Alternatives Considered

### Alternative 1: Concrete Types Only
```go
type FileOperations struct {
    // fields
}

func (f *FileOperations) ReadFile(...) {...}
```

**Pros:**
- Simpler, less boilerplate
- Fewer abstractions to maintain
- More direct code flow

**Cons:**
- Hard to test without actual file system
- Difficult to mock for unit tests
- Changes to implementation affect all callers
- Cannot easily swap implementations

**Rejected because:** Testing becomes difficult, especially for CLI commands that need to verify behavior without touching the real filesystem.

### Alternative 2: Functional Approach
```go
type FileOp func(path string) ([]byte, error)

func MakeReadFile() FileOp {
    return func(path string) ([]byte, error) {
        return os.ReadFile(path)
    }
}
```

**Pros:**
- Very flexible
- Composition over inheritance
- Easy to create test doubles

**Cons:**
- Unfamiliar pattern in Go
- Harder to discover related operations
- No clear grouping of related functionality
- More complex for newcomers

**Rejected because:** Go idiomatic code prefers interfaces over function composition for service abstractions.

### Alternative 3: Dependency Injection Framework
Use a DI framework like Wire or Dig.

**Pros:**
- Automatic wiring
- Centralized configuration
- Lifecycle management

**Cons:**
- Additional dependency
- Complexity for small project
- Steeper learning curve
- Compile-time magic can be confusing

**Rejected because:** Overkill for current project size. Simple constructor injection is sufficient.

## Consequences

### Positive

1. **Testability Improved**
   - Commands can be tested with mock services
   - File operations don't require temp directories for basic logic tests
   - Fast, deterministic unit tests

2. **Clear Contracts**
   - Interface defines exact behavior expected
   - Documentation lives with the interface
   - Breaking changes are explicit

3. **Future-Proof**
   - Can add caching layer without changing callers
   - Can add observability (logging, metrics) by wrapping
   - Can provide alternative implementations (e.g., in-memory for testing)

4. **Dependency Injection Friendly**
   - Easy to inject mock implementations
   - Clear dependency graph
   - Supports constructor-based injection

### Negative

1. **Slight Increase in Boilerplate**
   - Each service needs interface + implementation
   - Constructor functions even for stateless services
   - Mitigation: Pattern is consistent and easy to copy

2. **Indirection**
   - One extra step to find actual implementation
   - Mitigation: Clear naming convention (`Default*Service`)

3. **Interface Maintenance**
   - Adding methods requires updating interface and all implementations
   - Mitigation: Keep interfaces focused and stable

### Neutral

1. **Constructor Pattern**
   - All services use `New*Service()` constructors
   - Even stateless services get constructors for consistency
   - Establishes pattern for future services that may need initialization

## Implementation Notes

### Service Pattern Template
```go
// 1. Define interface
type XService interface {
    Method1(...) error
    Method2(...) (Type, error)
}

// 2. Default implementation
type DefaultXService struct {
    // fields if needed
}

// 3. Constructor (even if stateless)
func NewXService() *DefaultXService {
    return &DefaultXService{}
}

// 4. Interface methods
func (s *DefaultXService) Method1(...) error {
    // implementation
}
```

### Testing Pattern
```go
// Mock implementation for tests
type MockXService struct {
    Method1Func func(...) error
}

func (m *MockXService) Method1(...) error {
    if m.Method1Func != nil {
        return m.Method1Func(...)
    }
    return nil
}
```

## References

- [Go interfaces best practices](https://go.dev/doc/effective_go#interfaces)
- [Accept interfaces, return structs](https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b)
- golang-design-patterns skill - Interface design section

## Related Decisions

- ADR-002: Atomic File Writes (uses FileService interface)
- Future: ADR on dependency injection strategy

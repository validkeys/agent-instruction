# Style Anchors Index

Quick reference guide for finding the right style anchor for your task.

---

## Quick Lookup by Task

| What are you doing? | Use this style anchor |
|---------------------|----------------------|
| Creating a CLI command | [cobra-command-structure.md](cobra-command-structure.md) |
| Writing tests | [table-driven-testing.md](table-driven-testing.md) |
| Reading/writing files | [file-operations.md](file-operations.md) |
| Parsing JSON config | [json-config-handling.md](json-config-handling.md) |
| Handling errors | [error-handling.md](error-handling.md) |
| Resolving imports | [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) |

---

## Quick Lookup by Code Element

### Commands
- **Command structure** → [cobra-command-structure.md](cobra-command-structure.md)
- **Command flags** → [cobra-command-structure.md](cobra-command-structure.md) (Flag Definition Patterns)
- **Command output** → [cobra-command-structure.md](cobra-command-structure.md) (Output Handling)
- **Command errors** → [error-handling.md](error-handling.md) (CLI-Appropriate Error Formatting)

### Tests
- **Test structure** → [table-driven-testing.md](table-driven-testing.md)
- **Test helpers** → [table-driven-testing.md](table-driven-testing.md) (Helper Functions)
- **Test assertions** → [table-driven-testing.md](table-driven-testing.md) (Assertion Patterns)
- **Testing commands** → [table-driven-testing.md](table-driven-testing.md) (Testing Commands with Flags)

### Files
- **Write files safely** → [file-operations.md](file-operations.md) (Atomic Write Pattern)
- **Read files** → [file-operations.md](file-operations.md) (Reading Files Safely)
- **Create backups** → [file-operations.md](file-operations.md) (Backup Creation)
- **Managed sections** → [file-operations.md](file-operations.md) (Managed Section Replacement)
- **Directory creation** → [file-operations.md](file-operations.md) (Safe Directory Creation)

### JSON
- **Define structs** → [json-config-handling.md](json-config-handling.md) (Struct Definitions)
- **Load JSON** → [json-config-handling.md](json-config-handling.md) (Loading JSON Files)
- **Save JSON** → [json-config-handling.md](json-config-handling.md) (Saving JSON Files)
- **Validate JSON** → [json-config-handling.md](json-config-handling.md) (Validation Functions)

### Errors
- **Wrap errors** → [error-handling.md](error-handling.md) (Error Wrapping with %w)
- **Error messages** → [error-handling.md](error-handling.md) (Actionable Error Messages)
- **Validation errors** → [error-handling.md](error-handling.md) (Validation Error Messages)
- **Multiple errors** → [error-handling.md](error-handling.md) (Multi-Error Collection)

### Imports/Dependencies
- **Resolve imports** → [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md)
- **Detect cycles** → [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) (Cycle Detection with Path Stack)
- **Dependency order** → [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) (Topological Sort Pattern)

---

## Lookup by Milestone

### M0: Core Data Structures
- [json-config-handling.md](json-config-handling.md) - Config, RuleFile, Instruction structs
- [error-handling.md](error-handling.md) - Error types and validation

### M1: Services Layer
- [file-operations.md](file-operations.md) - File service implementation
- [json-config-handling.md](json-config-handling.md) - Config service
- [error-handling.md](error-handling.md) - Service error handling

### M2: Rule Management
- [json-config-handling.md](json-config-handling.md) - Rule file loading
- [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) - Import resolution
- [error-handling.md](error-handling.md) - Validation errors

### M3: Init Command
- [cobra-command-structure.md](cobra-command-structure.md) - Command definition
- [file-operations.md](file-operations.md) - Directory creation and backups
- [json-config-handling.md](json-config-handling.md) - Creating default config

### M4: Build Command
- [cobra-command-structure.md](cobra-command-structure.md) - Command with flags
- [file-operations.md](file-operations.md) - Managed section updates
- [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) - Import resolution

### M5: Add and List Commands
- [cobra-command-structure.md](cobra-command-structure.md) - Multiple commands
- [json-config-handling.md](json-config-handling.md) - Config updates

### M6: Validation Layer
- [error-handling.md](error-handling.md) - Comprehensive error messages
- [json-config-handling.md](json-config-handling.md) - Schema validation

### M7: Test Coverage
- [table-driven-testing.md](table-driven-testing.md) - All test patterns
- All other anchors for testing their respective patterns

---

## Decision Tree

```
START: What are you implementing?

├─ CLI Command?
│  └─ → cobra-command-structure.md
│
├─ Test file?
│  └─ → table-driven-testing.md
│
├─ File I/O operation?
│  ├─ Reading/writing files? → file-operations.md
│  ├─ JSON config? → json-config-handling.md
│  └─ Creating directories? → file-operations.md
│
├─ Error handling?
│  └─ → error-handling.md
│
├─ Import resolution?
│  └─ → graph-traversal-cycle-detection.md
│
└─ Not sure?
   └─ Check "Quick Lookup by Code Element" above
```

---

## Pattern Combinations

Some tasks require multiple style anchors. Here are common combinations:

### Building a new command
1. [cobra-command-structure.md](cobra-command-structure.md) - Command setup
2. [error-handling.md](error-handling.md) - Error returns
3. [file-operations.md](file-operations.md) OR [json-config-handling.md](json-config-handling.md) - Business logic
4. [table-driven-testing.md](table-driven-testing.md) - Tests

### Adding file operations
1. [file-operations.md](file-operations.md) - Atomic writes and backups
2. [error-handling.md](error-handling.md) - File error handling
3. [table-driven-testing.md](table-driven-testing.md) - File operation tests

### Working with JSON configs
1. [json-config-handling.md](json-config-handling.md) - Struct definitions and parsing
2. [error-handling.md](error-handling.md) - Validation errors
3. [file-operations.md](file-operations.md) - Safe file writes
4. [table-driven-testing.md](table-driven-testing.md) - JSON tests

### Implementing import resolution
1. [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) - Resolver algorithm
2. [json-config-handling.md](json-config-handling.md) - Loading rule files
3. [error-handling.md](error-handling.md) - Cycle error messages
4. [table-driven-testing.md](table-driven-testing.md) - Cycle detection tests

---

## Search by Keyword

| Keyword | Found in |
|---------|----------|
| `cobra.Command` | cobra-command-structure.md |
| `RunE` | cobra-command-structure.md |
| `Args` | cobra-command-structure.md |
| `Flags` | cobra-command-structure.md |
| `map[string]struct` | table-driven-testing.md |
| `t.Run` | table-driven-testing.md |
| `t.Helper()` | table-driven-testing.md |
| `t.TempDir()` | table-driven-testing.md |
| `AtomicWrite` | file-operations.md |
| `CreateBackup` | file-operations.md |
| `os.WriteFile` | file-operations.md |
| `os.ReadFile` | file-operations.md |
| `json.Marshal` | json-config-handling.md |
| `json.Unmarshal` | json-config-handling.md |
| `json:` tag | json-config-handling.md |
| `omitempty` | json-config-handling.md |
| `fmt.Errorf` | error-handling.md |
| `%w` | error-handling.md |
| `errors.Is` | error-handling.md |
| `DFS` | graph-traversal-cycle-detection.md |
| `visited` | graph-traversal-cycle-detection.md |
| `cycle detection` | graph-traversal-cycle-detection.md |
| `pathStack` | graph-traversal-cycle-detection.md |

---

## Anti-Patterns to Avoid

Each style anchor includes "Bad" examples. Here's where to find them:

| Anti-Pattern | See |
|--------------|-----|
| Using `Run` instead of `RunE` | [cobra-command-structure.md](cobra-command-structure.md) |
| Direct stdout instead of `cmd.OutOrStdout()` | [cobra-command-structure.md](cobra-command-structure.md) |
| Slice-based table tests | [table-driven-testing.md](table-driven-testing.md) |
| Direct file writes without temp file | [file-operations.md](file-operations.md) |
| Missing omitempty on optional fields | [json-config-handling.md](json-config-handling.md) |
| Generic error messages | [error-handling.md](error-handling.md) |
| Using panic for normal errors | [error-handling.md](error-handling.md) |
| Only using visited set (missing cycles) | [graph-traversal-cycle-detection.md](graph-traversal-cycle-detection.md) |

---

## When to Consult Multiple Anchors

**Scenario: Implementing init command**
```
1. cobra-command-structure.md → Command skeleton
2. file-operations.md → Create .agent-instruction directory
3. json-config-handling.md → Create default config.json
4. error-handling.md → Handle all error cases
5. table-driven-testing.md → Write tests
```

**Scenario: Adding import resolution**
```
1. graph-traversal-cycle-detection.md → Algorithm implementation
2. json-config-handling.md → Load rule files with imports
3. error-handling.md → Format cycle error messages
4. file-operations.md → Read rule files safely
5. table-driven-testing.md → Test cycle detection
```

**Scenario: Creating a new service**
```
1. json-config-handling.md → Data structures if JSON-related
2. file-operations.md → File operations if needed
3. error-handling.md → Error handling throughout
4. table-driven-testing.md → Service tests
```

---

## Quick Reference Card

**Before writing code, ask:**

- [ ] Am I creating a command? → cobra-command-structure.md
- [ ] Am I writing tests? → table-driven-testing.md
- [ ] Will this read/write files? → file-operations.md
- [ ] Will this parse JSON? → json-config-handling.md
- [ ] Will this return errors? → error-handling.md (always yes!)
- [ ] Will this traverse a graph? → graph-traversal-cycle-detection.md

**Default rule:** When in doubt, always check **error-handling.md** - every function needs proper error handling.

---

## Style Anchor Overview

| File | Lines | Use When |
|------|-------|----------|
| cobra-command-structure.md | 180 | Building CLI commands |
| table-driven-testing.md | 210 | Writing any tests |
| file-operations.md | 200 | File I/O operations |
| json-config-handling.md | 220 | JSON config/data |
| error-handling.md | 250 | Error handling (always) |
| graph-traversal-cycle-detection.md | 280 | Import resolution |

**Total patterns:** 6 anchors covering all core project patterns

---

## For AI Agents

**Rule:** Before generating code, identify which style anchor(s) apply to your task using this index, then follow the patterns exactly as shown.

**Order of consultation:**
1. Identify task type (command, test, file op, etc.)
2. Find relevant anchor(s) using tables above
3. Read the complete anchor before writing code
4. Follow patterns exactly - don't improvise variations
5. Always check error-handling.md for error patterns

**Multi-file tasks:** Consult all relevant anchors and combine patterns appropriately.

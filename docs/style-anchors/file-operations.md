# File Operations

**Purpose:** Safe file handling with atomic writes and backups
**Source:** Go standard library patterns and project requirements
**Use cases:** Init, build commands, and all file I/O operations

---

## Atomic Write Pattern

✅ **Good: Write to temp file, then rename (atomic)**

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
```

❌ **Bad: Direct write (not atomic, can corrupt file)**

```go
func badWrite(path string, content []byte) error {
    // If this fails mid-write, file is corrupted
    return os.WriteFile(path, content, 0644)
}
```

---

## Backup Creation

```go
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

// RestoreBackup restores a file from its backup
func RestoreBackup(path string) error {
    backupPath := path + ".backup"

    content, err := os.ReadFile(backupPath)
    if err != nil {
        return fmt.Errorf("read backup: %w", err)
    }

    if err := AtomicWrite(path, content, 0644); err != nil {
        return fmt.Errorf("restore from backup: %w", err)
    }

    return nil
}
```

---

## Managed Section Replacement

```go
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

## File Permission Preservation

```go
// UpdateFilePreservingMode updates a file while preserving its permissions
func UpdateFilePreservingMode(path string, content []byte) error {
    // Get current file info
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            // New file, use default permissions
            return AtomicWrite(path, content, 0644)
        }
        return fmt.Errorf("stat file: %w", err)
    }

    // Preserve original permissions
    return AtomicWrite(path, content, info.Mode())
}
```

---

## Safe Directory Creation

```go
// EnsureDir creates directory and all parent directories if needed
func EnsureDir(path string) error {
    if err := os.MkdirAll(path, 0755); err != nil {
        return fmt.Errorf("create directory %s: %w", path, err)
    }
    return nil
}

// EnsureDirForFile creates parent directory for a file path
func EnsureDirForFile(filePath string) error {
    dir := filepath.Dir(filePath)
    return EnsureDir(dir)
}
```

---

## Reading Files Safely

```go
// ReadFileOrEmpty reads file content, returns empty string if not found
func ReadFileOrEmpty(path string) (string, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return "", nil
        }
        return "", fmt.Errorf("read file %s: %w", path, err)
    }
    return string(content), nil
}

// MustReadFile reads file or fails with descriptive error
func MustReadFile(path string) ([]byte, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("file not found: %s", path)
        }
        return nil, fmt.Errorf("read file %s: %w", path, err)
    }
    return content, nil
}
```

---

## Complete File Update Flow

```go
// UpdateConfigFile demonstrates complete safe update flow
func UpdateConfigFile(path string, updater func(*Config) error) error {
    // 1. Create backup
    if err := CreateBackup(path); err != nil {
        return fmt.Errorf("create backup: %w", err)
    }

    // 2. Read existing file
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("read config: %w", err)
    }

    // 3. Parse config
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return fmt.Errorf("parse config: %w", err)
    }

    // 4. Apply updates
    if err := updater(&cfg); err != nil {
        return fmt.Errorf("update config: %w", err)
    }

    // 5. Marshal updated config
    newData, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal config: %w", err)
    }

    // 6. Write atomically
    if err := AtomicWrite(path, newData, 0644); err != nil {
        // Attempt to restore backup on failure
        RestoreBackup(path)
        return fmt.Errorf("write config: %w", err)
    }

    return nil
}
```

---

## Directory Traversal

```go
// WalkRuleFiles walks all .json files in rules directory
func WalkRuleFiles(rulesDir string, fn func(path string) error) error {
    entries, err := os.ReadDir(rulesDir)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // no rules directory, not an error
        }
        return fmt.Errorf("read rules directory: %w", err)
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        if !strings.HasSuffix(entry.Name(), ".json") {
            continue
        }

        path := filepath.Join(rulesDir, entry.Name())
        if err := fn(path); err != nil {
            return fmt.Errorf("process %s: %w", entry.Name(), err)
        }
    }

    return nil
}
```

---

## File Existence Checks

```go
// FileExists checks if file exists (not a directory)
func FileExists(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !info.IsDir()
}

// DirExists checks if directory exists
func DirExists(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return info.IsDir()
}

// IsEmptyDir checks if directory is empty
func IsEmptyDir(path string) (bool, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return false, fmt.Errorf("read directory: %w", err)
    }
    return len(entries) == 0, nil
}
```

---

## Error Handling for File Operations

✅ **Good: Distinguish between error types**

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("config not found at %s: run 'agent-instruction init' first", path)
        }
        if os.IsPermission(err) {
            return nil, fmt.Errorf("permission denied reading %s: check file permissions", path)
        }
        return nil, fmt.Errorf("read config: %w", err)
    }

    // ... parse and return
}
```

❌ **Bad: Generic error handling**

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err // User doesn't know why it failed
    }
    // ...
}
```

---

## Key Principles

1. **Always use atomic writes** - Write to temp file, then rename
2. **Create backups before modifying** - Allow recovery on failure
3. **Preserve file permissions** - Don't change modes unintentionally
4. **Handle missing files gracefully** - Distinguish not-found from other errors
5. **Clean up temp files on failure** - Don't leave garbage
6. **Use absolute paths when possible** - Avoid ambiguity

---

## Common Patterns

```go
// Pattern: Ensure directory exists before creating file
func createConfigFile(baseDir string) error {
    configPath := filepath.Join(baseDir, ".agent-instruction", "config.json")

    // Ensure parent directory exists
    if err := EnsureDirForFile(configPath); err != nil {
        return err
    }

    // Create config
    cfg := &Config{Version: "1.0", Frameworks: []string{"claude"}}
    data, _ := json.MarshalIndent(cfg, "", "  ")

    // Write atomically
    return AtomicWrite(configPath, data, 0644)
}

// Pattern: Update existing file safely
func addPackageToConfig(configPath, packageName string) error {
    return UpdateConfigFile(configPath, func(cfg *Config) error {
        cfg.Packages = append(cfg.Packages, packageName)
        return nil
    })
}
```

---

## References

- Go os package: https://pkg.go.dev/os
- Go filepath package: https://pkg.go.dev/path/filepath
- Atomic file writes: https://www.joeshaw.org/dont-defer-close-on-writable-files/
- Project requirements: `/Users/kydavis/Sites/agent-instruction/docs/plan/001-initial-buildout/technical-requirements.yaml` (lines 161-169)

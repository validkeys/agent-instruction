# ADR-002: Atomic File Writes

**Status**: Accepted

**Date**: 2026-04-13

**Deciders**: Development Team

---

## Context

The agent-instruction tool modifies critical files (CLAUDE.md, AGENTS.md, config.json, rule files) that developers rely on for AI agent behavior. File corruption or partial writes would be severe issues because:

1. **Corruption Risk**: If the program is interrupted mid-write (crash, SIGKILL, power loss), files could be left in an inconsistent state
2. **Critical Files**: CLAUDE.md and AGENTS.md control agent behavior - corruption affects development workflow
3. **No User Rollback**: Users can't easily recover from corrupted configuration files
4. **Concurrent Access**: Multiple processes might read files while we're writing

Standard approaches like direct writes with `os.WriteFile()` or `io.Writer` don't guarantee atomicity - interruptions can leave partial content.

## Decision

We will use **atomic file writes via temp file + rename** for all file write operations:

```go
func WriteAtomic(path string, content []byte) error {
    // 1. Create temp file in same directory (ensures same filesystem)
    dir := filepath.Dir(path)
    tempFile, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil {
        return fmt.Errorf("create temp file: %w", err)
    }
    tempPath := tempFile.Name()

    // 2. Ensure cleanup on any error
    defer func() {
        tempFile.Close()
        os.Remove(tempPath)
    }()

    // 3. Write to temp file
    if _, err := tempFile.Write(content); err != nil {
        return fmt.Errorf("write temp file: %w", err)
    }

    // 4. Sync to disk
    if err := tempFile.Sync(); err != nil {
        return fmt.Errorf("sync temp file: %w", err)
    }

    // 5. Close (required on Windows before rename)
    if err := tempFile.Close(); err != nil {
        return fmt.Errorf("close temp file: %w", err)
    }

    // 6. Atomic rename
    if err := os.Rename(tempPath, path); err != nil {
        return fmt.Errorf("rename temp file: %w", err)
    }

    return nil
}
```

This pattern is used everywhere files are written:
- Configuration files (config.json)
- Rule files (rules/*.json)
- Generated instruction files (CLAUDE.md, AGENTS.md)

## Alternatives Considered

### Alternative 1: Direct Write with os.WriteFile
```go
func WriteFile(path string, content []byte) error {
    return os.WriteFile(path, content, 0644)
}
```

**Pros:**
- Simple, one-line implementation
- Less code to maintain
- No temp files to clean up

**Cons:**
- Not atomic - interruption leaves partial content
- Can corrupt existing file
- No protection against concurrent readers seeing partial state
- Users lose data on interruption

**Rejected because:** Risk of corruption is unacceptable for critical configuration files.

### Alternative 2: Write-Ahead Log (WAL)
```go
// Write intent to log first
WriteLog(operation)
// Execute operation
WriteFile(path, content)
// Mark complete in log
MarkComplete(operation)
```

**Pros:**
- Can replay failed operations
- Provides audit trail
- Supports transactions

**Cons:**
- Much more complex
- Requires log management
- Performance overhead
- Overkill for single-file operations

**Rejected because:** Excessive complexity for our use case. Atomic rename is sufficient and simpler.

### Alternative 3: Copy-on-Write with Backup
```go
// Create backup
CopyFile(path, path+".backup")
// Write new content
os.WriteFile(path, content, 0644)
```

**Pros:**
- Keeps old version
- User can recover manually
- Simple to implement

**Cons:**
- Still not atomic
- Backup creation itself can fail
- Two files to manage
- Doesn't prevent partial writes

**Rejected because:** Doesn't solve atomicity problem. We implement backups separately (see backup.go).

### Alternative 4: Lock File
```go
// Create lock
lock := CreateLock(path + ".lock")
defer lock.Release()
// Write file
os.WriteFile(path, content, 0644)
```

**Pros:**
- Prevents concurrent writes
- Simple locking mechanism

**Cons:**
- Doesn't prevent partial writes from crashes
- Lock cleanup issues if process dies
- Can deadlock
- Still not atomic

**Rejected because:** Doesn't address corruption from interruption, only concurrent access.

## Consequences

### Positive

1. **Data Integrity**
   - Files are never partially written
   - Interruptions don't corrupt files
   - Readers always see complete content (old or new, never partial)

2. **Reliability**
   - Tool is safe to use in CI/CD pipelines
   - Safe during system crashes or forced termination
   - No manual recovery needed

3. **Cross-Platform**
   - Works on Unix, Linux, macOS, Windows
   - Handles filesystem differences (close before rename on Windows)
   - Temp file in same directory ensures same filesystem

4. **Permission Preservation**
   - Existing file permissions are preserved
   - New files get default 0644 permissions
   - No permission escalation issues

### Negative

1. **Disk Space**
   - Requires temporary space equal to file size
   - Brief period where both files exist
   - Mitigation: Temp files are immediately cleaned up

2. **Performance Overhead**
   - Extra Sync() call adds latency
   - Rename operation cost
   - Mitigation: Still fast for typical file sizes (< 1MB)

3. **Cleanup Responsibility**
   - Must ensure temp files are cleaned up
   - Defer ensures cleanup even on error
   - Mitigation: Pattern is consistent and tested

### Neutral

1. **Temp File Naming**
   - Uses `.tmp-*` prefix (hidden on Unix systems)
   - Random suffix prevents conflicts
   - Cleaned up automatically

## Implementation Notes

### Sync is Critical
```go
if err := tempFile.Sync(); err != nil {
    return fmt.Errorf("sync temp file: %w", err)
}
```

Without `Sync()`, data may only be in OS cache. A crash before the OS flushes cache would lose data even after rename.

### Same Directory is Required
```go
dir := filepath.Dir(path)
tempFile, err := os.CreateTemp(dir, ".tmp-*")
```

Temp file MUST be in same directory as target to ensure they're on the same filesystem. `os.Rename()` only works atomically within a single filesystem.

### Close Before Rename (Windows)
```go
if err := tempFile.Close(); err != nil {
    return fmt.Errorf("close temp file: %w", err)
}
```

Windows requires files to be closed before rename. Defer handles cleanup, but explicit close ensures rename succeeds.

### Defer Cleanup Pattern
```go
defer func() {
    tempFile.Close()
    os.Remove(tempPath)
}()
```

Ensures temp file is cleaned up even if write or sync fails. If rename succeeds, `os.Remove()` on non-existent file is harmless.

## Testing Strategy

Atomic writes are tested for:
- Normal write (new file)
- Overwrite existing file
- Permission preservation
- Temp file cleanup on error
- Large file handling (1MB+)
- Concurrent write safety (implicit via atomicity)

See `internal/files/atomic_test.go` for complete test suite.

## References

- [POSIX rename() atomicity guarantees](https://pubs.opengroup.org/onlinepubs/9699919799/functions/rename.html)
- [Atomic file writes in Go](https://stackoverflow.com/questions/37329998/atomic-file-write-in-go)
- golang-design-patterns skill - Resource management section

## Related Decisions

- ADR-001: Interface-Based Service Design (FileService uses atomic writes)
- Backup strategy (backup.go) - separate concern from atomicity

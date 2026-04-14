# build Command

Generate instruction files from rule configurations.

## Synopsis

```bash
agent-instruction build [flags]
```

## Description

The `build` command generates or updates `CLAUDE.md` and/or `AGENTS.md` files throughout your repository based on the centralized rule configurations.

It discovers packages, composes instructions from global and package-level rules, and writes managed sections while preserving custom content.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dry-run` | boolean | `false` | Preview changes without writing files |
| `--verbose` | boolean | `false` | Show detailed progress output |
| `--no-parallel` | boolean | `false` | Disable parallel package processing |

## Behavior

### Package Discovery

The command discovers packages based on `config.json`:

- **`"packages": ["auto"]`** - Automatically finds all directories with `agent-instruction.json`
- **`"packages": ["path1", "path2"]`** - Uses specified package paths
- **`"packages": []`** - Processes only the root directory

### Instruction Composition

For each package:

1. Loads global rules from `.agent-instruction/rules/global.json`
2. Loads package rules from `<package>/agent-instruction.json` (if exists)
3. Resolves imports recursively
4. Detects circular dependencies
5. Composes final instruction set

### File Generation

For each configured framework:

1. Generates markdown content from instructions
2. Wraps content in managed markers: `<!-- BEGIN AGENT-INSTRUCTION -->` and `<!-- END AGENT-INSTRUCTION -->`
3. Preserves content outside markers
4. Writes to `CLAUDE.md` or `AGENTS.md`

### Parallel Processing

By default, processes multiple packages in parallel for speed. Use `--no-parallel` to process sequentially (useful for debugging).

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success - all packages processed |
| `1` | Not initialized (run `init` first) |
| `1` | Configuration error (invalid `config.json`) |
| `1` | Import resolution error (circular dependency, missing file) |
| `1` | File write error (permissions, disk space) |
| `1` | One or more packages failed to build |

## Examples

### Basic Build

```bash
agent-instruction build
```

Generates instruction files for all packages.

### Preview Changes

```bash
agent-instruction build --dry-run
```

Shows what would be generated without writing files:

```
DRY RUN: Preview mode - no files will be written

Discovering packages...
Found 3 package(s):
  - .
  - packages/api
  - packages/lib

  Would generate: CLAUDE.md
  Would generate: AGENTS.md
  Would generate: packages/api/CLAUDE.md
  Would generate: packages/api/AGENTS.md
  Would generate: packages/lib/CLAUDE.md
  Would generate: packages/lib/AGENTS.md

✓ DRY RUN: Checked 3 package(s) in 45ms
```

### Verbose Output

```bash
agent-instruction build --verbose
```

Shows detailed progress:

```
Discovering packages...
Found 3 package(s):
  - .
  - packages/api
  - packages/lib

Building packages in parallel...

Processing .
  ✓ Generated CLAUDE.md
  ✓ Generated AGENTS.md
Processing packages/api...
  ✓ Generated CLAUDE.md
  ✓ Generated AGENTS.md
Processing packages/lib...
  ✓ Generated CLAUDE.md
  ✓ Generated AGENTS.md

✓ Successfully processed 3 package(s) in 123ms
```

### Sequential Processing

```bash
agent-instruction build --no-parallel
```

Processes packages one at a time (useful for debugging build issues).

### Combine Flags

```bash
agent-instruction build --dry-run --verbose
```

Preview with detailed output to understand what will change.

## Common Patterns

### After Adding Rules

```bash
# Add new rule
agent-instruction add "New rule content" --rule global

# Rebuild all files
agent-instruction build
```

### After Editing Rule Files

```bash
# Edit rule file manually
vim .agent-instruction/rules/global.json

# Rebuild
agent-instruction build
```

### Verify Before Committing

```bash
# Build with verbose output
agent-instruction build --verbose

# Check differences
git diff CLAUDE.md AGENTS.md

# Commit if satisfied
git add .agent-instruction/ **/CLAUDE.md **/AGENTS.md
git commit -m "Update agent instructions"
```

### Monorepo Incremental Build

```bash
# Build everything
agent-instruction build

# Later, rebuild after changes
agent-instruction build --verbose
```

### Debug Build Issues

```bash
# Use verbose and sequential processing
agent-instruction build --verbose --no-parallel

# Or preview without writing
agent-instruction build --dry-run --verbose
```

## Managed Sections

### Marker Format

Generated files contain HTML-style comment markers:

```markdown
# Custom Header

Your custom content here is preserved.

<!-- BEGIN AGENT-INSTRUCTION -->
# Generated Content

This section is replaced on each build.
<!-- END AGENT-INSTRUCTION -->

# Custom Footer

More custom content preserved.
```

### Preserving Custom Content

- Content **outside markers** is preserved
- Content **inside markers** is regenerated
- Manual edits inside markers are **lost** on rebuild

**Best Practice:** Keep custom content outside markers, or edit rule files instead.

### First Build

On first build, if no markers exist:
- Wraps all content in markers
- Preserves existing content
- Adds markers for future builds

## Troubleshooting

### Not Initialized

**Error:**
```
Error: not initialized: run 'agent-instruction init' first
```

**Solution:**
Run `init` command first:

```bash
agent-instruction init
```

### Invalid Configuration

**Error:**
```
Error: invalid config: at least one framework is required
```

**Solution:**
Edit `.agent-instruction/config.json` to add frameworks:

```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

### Circular Import

**Error:**
```
Error: compose instructions: circular import detected: global -> testing -> global
```

**Solution:**
Remove circular imports from rule files:

```bash
# Find imports
grep -r "imports" .agent-instruction/rules/

# Edit to remove circular reference
vim .agent-instruction/rules/testing.json
```

### File Not Found

**Error:**
```
Error: compose instructions for packages/api: load rule file: open security.json: no such file or directory
```

**Solution:**
Create the missing rule file or remove the import:

```bash
# Create missing file
cat > .agent-instruction/rules/security.json <<EOF
{
  "title": "Security Rules",
  "instructions": []
}
EOF

# Or remove import from package config
vim packages/api/agent-instruction.json
```

### Permission Denied

**Error:**
```
Error: write CLAUDE.md: permission denied
```

**Solution:**
Check file permissions:

```bash
# Check permissions
ls -la CLAUDE.md

# Fix if needed
chmod u+w CLAUDE.md
```

### Build Failed for Some Packages

**Output:**
```
⚠ Processed 2 package(s) with 1 error(s) in 234ms
```

**Solution:**
Use verbose mode to see which package failed:

```bash
agent-instruction build --verbose --no-parallel
```

## Performance

### Timing Expectations

For typical repositories:

| Packages | Time (parallel) | Time (sequential) |
|----------|----------------|-------------------|
| 1-5 | < 100ms | < 200ms |
| 10-20 | < 500ms | < 1s |
| 50+ | < 2s | < 5s |

### Optimization

For large monorepos:
- Keep rule files under 1MB
- Minimize import chains (max depth 5-10)
- Use parallel processing (default)
- Cache Go dependencies in CI/CD

## Related Commands

- [`init`](init.md) - Initialize before building
- [`add`](add.md) - Add rules, then rebuild
- [`list`](list.md) - View rules before building

## See Also

- [Configuration Documentation](../configuration.md)
- [Import System](../configuration.md#import-system)
- [Examples](../examples.md)
- [Main README](../../README.md)

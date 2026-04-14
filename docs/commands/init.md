# init Command

Initialize agent-instruction in your repository.

## Synopsis

```bash
agent-instruction init [flags]
```

## Description

The `init` command sets up the `.agent-instruction` directory structure in your repository, creating the necessary configuration files and directory layout for managing AI instruction files.

This is typically the first command you run when adopting agent-instruction in a project.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--non-interactive` | boolean | `false` | Skip interactive prompts and use default values |
| `--frameworks` | string | - | Comma-separated list of frameworks: `claude`, `agents` |
| `--packages` | string | - | Comma-separated package paths or `auto` for discovery |

## Behavior

### Interactive Mode (default)

When run without `--non-interactive`, the command:

1. Detects existing `CLAUDE.md` or `AGENTS.md` files
2. Prompts whether to create backups (`.backup` extension)
3. Asks which frameworks to support (Claude, Agents, or both)
4. Asks about package discovery (auto, manual list, or root only)
5. Creates directory structure and initial configuration

### Non-Interactive Mode

When run with `--non-interactive`:

1. Automatically backs up existing instruction files
2. Uses both frameworks (`claude` and `agents`) by default
3. Uses automatic package discovery by default
4. Creates structure without prompts

### Created Structure

```
.agent-instruction/
├── config.json         # Main configuration
└── rules/
    └── global.json     # Global instruction rules
```

### Existing Files

If `CLAUDE.md` or `AGENTS.md` exist in the root directory:
- Backups are created: `CLAUDE.md.backup`, `AGENTS.md.backup`
- Original files are left untouched
- You can manually merge content into new managed files

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Already initialized (`.agent-instruction` exists) |
| `1` | File system error (permissions, disk space) |
| `1` | Invalid flag values |

## Examples

### Basic Interactive Setup

```bash
agent-instruction init
```

Prompts for all configuration options.

### Non-Interactive with Defaults

```bash
agent-instruction init --non-interactive
```

Creates structure with:
- Both Claude and Agents frameworks
- Automatic package discovery
- Backup of existing files

### Claude Framework Only

```bash
agent-instruction init --non-interactive --frameworks claude
```

Generates only `CLAUDE.md` files.

### Specific Packages

```bash
agent-instruction init --non-interactive --packages api,lib,worker
```

Configures to manage instruction files in:
- `api/`
- `lib/`
- `worker/`

### Root Only (No Packages)

```bash
agent-instruction init --non-interactive --packages ""
```

Manages only root-level instruction files, no package discovery.

### Multiple Frameworks with Specific Packages

```bash
agent-instruction init \
  --non-interactive \
  --frameworks claude,agents \
  --packages app,services/api,services/worker
```

## Common Patterns

### Convert Existing Repository

If you already have `CLAUDE.md` files:

```bash
# 1. Initialize with backup
agent-instruction init --non-interactive

# 2. Manually merge content from backups
cat CLAUDE.md.backup

# 3. Add content to rule files
agent-instruction add "Your existing rules..." --rule global

# 4. Build to regenerate files
agent-instruction build
```

### Monorepo Setup

```bash
# 1. Initialize at repository root
cd /path/to/monorepo
agent-instruction init --non-interactive

# 2. Edit config to specify packages
cat > .agent-instruction/config.json <<EOF
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["packages/app", "packages/lib", "services/api"]
}
EOF

# 3. Build to create files in all packages
agent-instruction build
```

### Single Package Project

```bash
# Initialize without package discovery
agent-instruction init --non-interactive --packages ""
```

## Troubleshooting

### Already Initialized

**Error:**
```
Error: already initialized: .agent-instruction directory exists
Use 'agent-instruction build' to regenerate files
```

**Solution:**
The directory already exists. Use `build` command to regenerate files, or remove `.agent-instruction/` to start fresh.

### Permission Denied

**Error:**
```
Error: create directory structure: permission denied
```

**Solution:**
Check directory permissions. You need write access to the current directory.

```bash
# Check permissions
ls -la

# Fix if needed (example)
chmod u+w .
```

### Invalid Framework

**Error:**
```
Error: invalid framework: xyz (must be 'claude' or 'agents')
```

**Solution:**
Use only valid framework names:

```bash
agent-instruction init --non-interactive --frameworks claude,agents
```

### Backup Creation Failed

**Error:**
```
Error: create backup of CLAUDE.md: permission denied
```

**Solution:**
Check that existing instruction files are readable:

```bash
# Check file permissions
ls -la CLAUDE.md

# Fix if needed
chmod u+r CLAUDE.md
```

## Related Commands

- [`build`](build.md) - Generate instruction files after initialization
- [`add`](add.md) - Add rules to rule files
- [`list`](list.md) - View configured rules

## Configuration Files

After initialization, edit these files:

### config.json

See [Configuration Reference](../configuration.md#configjson) for full documentation.

### rules/global.json

See [Configuration Reference](../configuration.md#rule-files) for rule file format.

## Next Steps

After initialization:

1. **Edit rules**: Add your project-specific instructions to `.agent-instruction/rules/global.json`
2. **Build files**: Run `agent-instruction build` to generate instruction files
3. **Add rules**: Use `agent-instruction add` to quickly add new rules during development
4. **Commit**: Add `.agent-instruction/` to version control

```bash
git add .agent-instruction/
git commit -m "Initialize agent-instruction"
```

## See Also

- [Configuration Documentation](../configuration.md)
- [Examples](../examples.md)
- [Main README](../../README.md)

# Installation Guide

Complete instructions for installing the agent-instruction skill for Claude Code.

## Prerequisites

Before installing the skill, ensure you have:

1. **Claude Code** installed and working
   - Download from: https://claude.ai/code
   - Verify: `claude --version`

2. **agent-instruction CLI** installed
   - Install: `go install github.com/validkeys/agent-instruction@latest`
   - Verify: `agent-instruction --version`

3. **Git** (for cloning the repository)
   - Verify: `git --version`

## Quick Installation

```bash
# Clone the repository
git clone https://github.com/validkeys/agent-instruction.git
cd agent-instruction

# Run installation script
bash skill/install.sh
```

The script will:
- ✓ Detect or create Claude Code skills directory
- ✓ Copy the skill to `~/.local/share/claude/skills/agent-instruction/`
- ✓ Validate the installation
- ✓ Check for agent-instruction CLI

## Custom Installation Location

Override the default skills directory:

```bash
# Use custom location
export CLAUDE_SKILLS_DIR=~/my-custom-path/skills
bash skill/install.sh
```

## Dry Run

Test the installation without making changes:

```bash
bash skill/install.sh --dry-run
```

This shows what would happen without actually installing.

## Manual Installation

If you prefer manual installation:

```bash
# Create skills directory
mkdir -p ~/.local/share/claude/skills

# Copy skill
cp -r skill/agent-instruction ~/.local/share/claude/skills/

# Verify
ls ~/.local/share/claude/skills/agent-instruction/skill.yaml
```

## Verification

### 1. Check Skill Files

```bash
ls ~/.local/share/claude/skills/agent-instruction/
```

**Expected output:**
```
CHANGELOG.md
README.md
examples/
install.md
skill.yaml
```

### 2. Verify skill.yaml

```bash
cat ~/.local/share/claude/skills/agent-instruction/skill.yaml
```

Should show valid YAML with name, version, commands, etc.

### 3. Check Claude Code Recognition

```bash
claude skills list
```

Should include `agent-instruction` in the list.

### 4. Verify CLI Access

```bash
agent-instruction --version
```

Should display version information.

## Troubleshooting

### Installation Script Fails

**Error:** `Permission denied`

```bash
# Make script executable
chmod +x skill/install.sh

# Run again
bash skill/install.sh
```

**Error:** `Skills directory does not exist`

```bash
# Create manually
mkdir -p ~/.local/share/claude/skills

# Run again
bash skill/install.sh
```

### Skill Not Recognized

**Problem:** Claude Code doesn't see the skill

**Solutions:**

1. Check installation location:
   ```bash
   echo ~/.local/share/claude/skills/agent-instruction
   ls -la ~/.local/share/claude/skills/agent-instruction/
   ```

2. Verify skill.yaml exists and is valid:
   ```bash
   cat ~/.local/share/claude/skills/agent-instruction/skill.yaml
   ```

3. Restart Claude Code if necessary

### CLI Not Found

**Problem:** `agent-instruction: command not found`

**Solutions:**

1. Install the CLI:
   ```bash
   go install github.com/validkeys/agent-instruction@latest
   ```

2. Ensure Go bin is in PATH:
   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"

   # Add to ~/.bashrc or ~/.zshrc for persistence
   echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
   ```

3. Verify installation:
   ```bash
   which agent-instruction
   agent-instruction --version
   ```

### Skill Already Installed

**Problem:** Installation script finds existing skill

**Solutions:**

1. Overwrite during installation (script will prompt)

2. Remove manually first:
   ```bash
   rm -rf ~/.local/share/claude/skills/agent-instruction
   bash skill/install.sh
   ```

3. Install to different location:
   ```bash
   export CLAUDE_SKILLS_DIR=~/skills-backup
   bash skill/install.sh
   ```

## Updating

To update the skill to a newer version:

```bash
# Pull latest changes
cd agent-instruction
git pull origin main

# Reinstall (overwrites existing)
bash skill/install.sh
```

The installation script will prompt before overwriting.

## Uninstallation

To remove the skill:

```bash
# Remove skill directory
rm -rf ~/.local/share/claude/skills/agent-instruction

# Verify removal
ls ~/.local/share/claude/skills/
```

The agent-instruction CLI remains installed. To remove it:

```bash
# Remove CLI binary
rm $(which agent-instruction)
```

## Platform-Specific Notes

### macOS

Default skills directory: `~/.local/share/claude/skills/`

No special configuration needed.

### Linux

Default skills directory: `~/.local/share/claude/skills/`

Ensure `~/.local/share` directory exists:
```bash
mkdir -p ~/.local/share
```

### Windows (WSL)

Use the Linux installation method within WSL.

Default skills directory: `~/.local/share/claude/skills/`

## Next Steps

After successful installation:

1. **Initialize a project**
   ```bash
   cd your-project
   agent-instruction init
   ```

2. **Test in Claude Code**
   - Start a Claude Code session
   - Try: "add a rule that we always use TDD"
   - Verify skill responds

3. **Read the examples**
   ```bash
   cat ~/.local/share/claude/skills/agent-instruction/examples/*.md
   ```

## Support

If you encounter issues:

1. Check the [Troubleshooting section](#troubleshooting) above
2. Review the [README](README.md) for usage guidance
3. Open an issue: https://github.com/validkeys/agent-instruction/issues

## Contributing

Found a bug in the installation script? Contributions welcome!

See: https://github.com/validkeys/agent-instruction/blob/main/CONTRIBUTING.md

# agent-instruction Skill for Claude Code

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/validkeys/agent-instruction)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A Claude Code skill that enables natural language management of AI instruction files in monorepos during development sessions.

## Overview

The `agent-instruction` skill allows Claude Code to:

- Add instruction rules during development conversations
- Build and update CLAUDE.md and AGENTS.md files
- List existing rules to avoid duplicates
- Manage rule files organized by category (testing, security, golang, etc.)

This skill bridges the gap between discovering best practices during development and codifying them as durable instructions.

## Installation

### Prerequisites

1. [Claude Code](https://claude.ai/code) installed
2. [agent-instruction CLI](https://github.com/validkeys/agent-instruction) installed and in PATH

### Install the Skill

```bash
# Clone or download the agent-instruction repository
cd agent-instruction

# Run the installation script
bash skill/install.sh
```

The script copies the skill to `~/.local/share/claude/skills/agent-instruction/`.

### Verify Installation

```bash
# Check Claude Code recognizes the skill
claude skills list | grep agent-instruction

# Verify CLI is available
agent-instruction --version
```

## Commands

### add_rule

Add a new instruction rule to a rule file.

```bash
agent-instruction add "[rule content]" --title "[title]" --rule="[rule-file]"
```

**Parameters:**
- `content` (required): The instruction rule text
- `--title` (optional): Heading for the rule section
- `--rule` (optional): Target rule file name (prompts if omitted)

**Examples:**

```bash
# Add a testing standard
agent-instruction add "Always use table-driven tests in Go" \
  --title="Testing Standards" \
  --rule="testing"

# Add with interactive rule file selection
agent-instruction add "Never use global variables" \
  --title="Code Standards"

# Add a security rule
agent-instruction add "Validate all user input before database queries" \
  --rule="security"
```

### build

Generate or update instruction files from rule configurations.

```bash
agent-instruction build
```

Reads all rule files from `.agent-instruction/rules/` and generates:
- `CLAUDE.md` - Instructions for Claude AI assistant
- `AGENTS.md` - Instructions for AI agents in monorepo workspaces

**Examples:**

```bash
# Rebuild after adding rules
agent-instruction build

# Propagate rule changes to instruction files
agent-instruction build
```

### list

Display all available rule files and their contents.

```bash
agent-instruction list [--verbose]
```

**Parameters:**
- `--verbose` (optional): Show full rule content with formatting

**Examples:**

```bash
# List rule files with summaries
agent-instruction list

# Show detailed rule content
agent-instruction list --verbose
```

## Natural Language Triggers

Claude Code will recognize these phrases and invoke the skill:

- "add rule"
- "update CLAUDE.md"
- "update AGENTS.md"
- "list instruction rules"
- "build instruction files"
- "show available rules"

## Usage Examples

### Scenario 1: Discovered Pattern During Code Review

**Conversation:**

> User: "Claude keeps using `fmt.Sprintf` for simple string concatenation."
>
> Claude: "I'll add a rule to remember this for Go development."

```bash
agent-instruction add "Use string concatenation or strings.Builder instead of fmt.Sprintf for simple string building" \
  --title="Performance Guidelines" \
  --rule="golang"
```

### Scenario 2: Adding Multiple Rules

**Workflow:**

1. Add rules during development:
   ```bash
   agent-instruction add "Write tests before implementation" --rule="testing"
   agent-instruction add "Mock external dependencies" --rule="testing"
   agent-instruction add "Use bcrypt for password hashing" --rule="security"
   ```

2. Review what was added:
   ```bash
   agent-instruction list
   ```

3. Build instruction files:
   ```bash
   agent-instruction build
   ```

4. Verify files updated:
   ```bash
   git diff CLAUDE.md AGENTS.md
   ```

### Scenario 3: Checking Before Adding

Before adding a new rule, check if it already exists:

```bash
# List all rules
agent-instruction list --verbose

# Add only if not duplicate
agent-instruction add "New rule content..." --rule="category"
```

## Best Practices

### When to Add Rules

- User says "remember this for next time"
- Discovered a pattern that should be codified
- Found a mistake that keeps recurring
- Established a new standard for the project

### Organizing Rules

- Use descriptive, clear titles
- Group related rules in the same file:
  - `testing.json` - Test patterns and standards
  - `security.json` - Security requirements
  - `golang.json` - Go-specific conventions
  - `global.json` - Cross-cutting concerns

### Workflow

1. **Add** rules as you discover patterns
2. **List** to review before adding more
3. **Build** to propagate to instruction files
4. **Commit** changes with meaningful message

### Avoiding Duplicates

Always list existing rules before adding:

```bash
agent-instruction list --verbose | grep "pattern you want to add"
```

## Integration with Claude Code

The skill integrates seamlessly with Claude Code sessions:

```
You: "Add a rule that we always validate email addresses"

Claude: I'll add that rule to the security configuration.
        [Executes: agent-instruction add "Validate email addresses..." --rule="security"]

        Rule added successfully. Run `agent-instruction build` to update instruction files.
```

## Troubleshooting

### Skill Not Recognized

```bash
# Check Claude Code skills directory
ls ~/.local/share/claude/skills/

# Reinstall if missing
bash skill/install.sh
```

### CLI Not Found

```bash
# Check if agent-instruction is in PATH
which agent-instruction

# Install CLI if missing
go install github.com/validkeys/agent-instruction@latest
```

### Build Fails

```bash
# Ensure .agent-instruction/config.yaml exists
agent-instruction init

# Check rule file syntax
agent-instruction list --verbose
```

### Permission Denied

```bash
# Make install script executable
chmod +x skill/install.sh

# Check directory permissions
ls -la ~/.local/share/claude/skills/
```

## Contributing

Issues and pull requests welcome at [github.com/validkeys/agent-instruction](https://github.com/validkeys/agent-instruction).

## License

MIT License - see LICENSE file for details.

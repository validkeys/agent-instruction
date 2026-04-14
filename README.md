# agent-instruction

A command-line tool for managing AI instruction files (CLAUDE.md and AGENTS.md) in monorepos. Centralize your AI instructions, eliminate manual copy-paste, and build consistent behavior across all packages.

## Overview

Managing CLAUDE.md and AGENTS.md files across a monorepo is error-prone and time-consuming. `agent-instruction` provides a single source of truth for AI instructions, automatically propagating rules across all packages while preserving custom content.

### Features

- **Centralized Rule Management** - Define rules once in `.agent-instruction/rules/`
- **Automatic Propagation** - Generate instruction files for all packages with one command
- **Import System** - Compose instructions from global and package-level rules
- **Safe Updates** - Preserve custom content outside managed sections
- **Multiple Frameworks** - Support for both Claude and Agents AI frameworks
- **Quick Rule Addition** - Add rules during development with simple commands
- **Parallel Processing** - Fast builds across large monorepos

## Installation

### Prerequisites

- Go 1.21 or later
- macOS (initial version)

### Install from Source

```bash
# Clone the repository
git clone https://github.com/validkeys/agent-instruction.git
cd agent-instruction

# Build and install
go install

# Verify installation
agent-instruction --version
```

### Install via Go

```bash
go install github.com/validkeys/agent-instruction@latest
```

## Quick Start

### 1. Initialize Your Repository

```bash
cd /path/to/your/monorepo
agent-instruction init
```

This creates:
```
.agent-instruction/
├── config.json          # Main configuration
└── rules/
    └── global.json      # Global instruction rules
```

### 2. Edit Your Rules

Edit `.agent-instruction/rules/global.json`:

```json
{
  "title": "Global Instructions",
  "instructions": [
    {
      "heading": "Code Standards",
      "rule": "Always use explicit error handling in Go. Never ignore errors."
    },
    {
      "heading": "Testing",
      "rule": "Write table-driven tests for all new functions."
    }
  ]
}
```

### 3. Build Instruction Files

```bash
agent-instruction build
```

This generates `CLAUDE.md` and/or `AGENTS.md` files in your root and package directories, with content wrapped in managed markers:

```markdown
<!-- BEGIN AGENT-INSTRUCTION -->
# Global Instructions

## Code Standards
Always use explicit error handling in Go. Never ignore errors.

## Testing
Write table-driven tests for all new functions.
<!-- END AGENT-INSTRUCTION -->
```

### 4. Add Rules During Development

```bash
agent-instruction add "Use bcrypt for password hashing" \
  --title "Security" \
  --rule security
```

Then rebuild:
```bash
agent-instruction build
```

## Commands

### init

Initialize agent-instruction in your repository.

```bash
agent-instruction init [flags]
```

**Flags:**
- `--non-interactive` - Skip prompts and use defaults
- `--frameworks string` - Comma-separated frameworks: `claude`, `agents` (default: both)
- `--packages string` - Comma-separated package paths or `auto` (default: auto)

**Examples:**

```bash
# Interactive setup with prompts
agent-instruction init

# Non-interactive with defaults
agent-instruction init --non-interactive

# Only Claude framework
agent-instruction init --non-interactive --frameworks claude

# Specific packages
agent-instruction init --non-interactive --packages app,lib,services
```

**Behavior:**
- Creates `.agent-instruction/` directory structure
- Backs up existing CLAUDE.md/AGENTS.md files (`.backup` extension)
- Creates default `config.json` and `rules/global.json`

### build

Generate instruction files from rule configurations.

```bash
agent-instruction build [flags]
```

**Flags:**
- `--dry-run` - Preview changes without writing files
- `--verbose` - Show detailed progress output
- `--no-parallel` - Disable parallel package processing

**Examples:**

```bash
# Build all packages
agent-instruction build

# Preview what would be generated
agent-instruction build --dry-run

# Verbose output with timing
agent-instruction build --verbose

# Sequential processing (for debugging)
agent-instruction build --no-parallel
```

**Behavior:**
- Discovers packages based on `config.json`
- Composes instructions from global and package-level rules
- Generates files for each configured framework
- Preserves content outside managed markers
- Processes packages in parallel by default

### add

Add a new instruction rule to a rule file.

```bash
agent-instruction add <rule-content> [flags]
```

**Flags:**
- `--title string` - Optional heading for the rule
- `--rule string` - Target rule file name (prompts if omitted)

**Examples:**

```bash
# Add with title and specific rule file
agent-instruction add "Always validate user input" \
  --title "Security Standards" \
  --rule security

# Interactive rule file selection
agent-instruction add "Use dependency injection for services"

# Add to global rules
agent-instruction add "Document all exported functions" --rule global
```

**Behavior:**
- Appends instruction to specified rule file
- Creates rule file if it doesn't exist
- Prompts for file selection if `--rule` omitted

### list

Display all instruction rules.

```bash
agent-instruction list [flags]
```

**Flags:**
- `--verbose` - Show full rule content

**Examples:**

```bash
# Summary view
agent-instruction list

# Full content with all details
agent-instruction list --verbose
```

**Output:**
```
📄 global.json
   3 instruction(s)
   - Code Standards
   - Testing
   - Security

📄 testing.json
   2 instruction(s)
   - Test Structure
   - Coverage Requirements
```

## Configuration

### config.json

Located at `.agent-instruction/config.json`:

```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

**Fields:**

- `version` (required) - Config format version (currently "1.0")
- `frameworks` (required) - Array of frameworks: `["claude"]`, `["agents"]`, or both
- `packages` (required) - Package discovery:
  - `["auto"]` - Automatically discover packages
  - `["path1", "path2"]` - Specific package paths relative to repo root
  - `[]` - Root only (no package discovery)

### Rule Files

Rule files are JSON documents in `.agent-instruction/rules/`:

```json
{
  "title": "Security Guidelines",
  "imports": ["global"],
  "instructions": [
    {
      "heading": "Authentication",
      "rule": "Use bcrypt for password hashing with cost factor 12+",
      "references": [
        "https://pkg.go.dev/golang.org/x/crypto/bcrypt"
      ]
    },
    {
      "rule": "Always validate and sanitize user input"
    }
  ]
}
```

**Fields:**

- `title` (required) - Display title for the rule file
- `imports` (optional) - Array of rule files to import (without .json extension)
- `instructions` (required) - Array of instruction objects:
  - `heading` (optional) - Section heading
  - `rule` (required) - The instruction text
  - `references` (optional) - Array of reference URLs or file paths

### Import System

Rule files can import other rule files to compose instructions:

```json
{
  "title": "Backend Service Rules",
  "imports": ["global", "testing", "security"],
  "instructions": [
    {
      "heading": "API Design",
      "rule": "Use REST conventions for HTTP APIs"
    }
  ]
}
```

**Import Resolution:**

1. Imports are resolved recursively
2. Circular dependencies are detected and prevented
3. Each rule file is included only once
4. Instructions maintain import order

### Package-Level Configuration

Create `agent-instruction.json` in any package directory:

```json
{
  "title": "API Service Instructions",
  "imports": ["global", "security"],
  "instructions": [
    {
      "heading": "Rate Limiting",
      "rule": "Implement rate limiting on all public endpoints"
    }
  ]
}
```

This configuration:
- Imports global and security rules
- Adds package-specific instructions
- Generates instruction files in that package directory

## Workflows

### Initial Setup for Existing Monorepo

```bash
# 1. Initialize
cd /path/to/monorepo
agent-instruction init

# 2. Review generated files
cat .agent-instruction/config.json
cat .agent-instruction/rules/global.json

# 3. Add your rules
agent-instruction add "Your first rule" --rule global

# 4. Build instruction files
agent-instruction build --verbose

# 5. Verify generated files
git status
git diff CLAUDE.md AGENTS.md
```

### Adding Rules During Development

When you discover a pattern that should be codified:

```bash
# Add the rule
agent-instruction add "Always use context.Context for cancellation" \
  --title "Concurrency Patterns" \
  --rule golang

# Rebuild immediately
agent-instruction build

# Commit with the code changes
git add .agent-instruction/ CLAUDE.md AGENTS.md
git commit -m "Add concurrency pattern rule"
```

### Organizing Rules by Category

```bash
# Create category-specific rule files
agent-instruction add "Mock external dependencies" \
  --title "Test Isolation" \
  --rule testing

agent-instruction add "Sanitize all database queries" \
  --title "SQL Injection Prevention" \
  --rule security

agent-instruction add "Use explicit types instead of interface{}" \
  --title "Type Safety" \
  --rule golang

# List to see organization
agent-instruction list
```

### Package-Specific Instructions

For packages with unique requirements:

```bash
cd packages/api-service

# Create package-level config
cat > agent-instruction.json <<EOF
{
  "title": "API Service Instructions",
  "imports": ["global", "security"],
  "instructions": [
    {
      "heading": "OpenAPI",
      "rule": "Generate OpenAPI spec from code, never write by hand"
    }
  ]
}
EOF

# Build from root to update all packages
cd ../..
agent-instruction build
```

### Reviewing Before Adding

Avoid duplicate rules:

```bash
# List existing rules
agent-instruction list --verbose | grep -i "pattern you want to add"

# Add only if unique
agent-instruction add "Your new rule" --rule category
```

## File Structure

### Managed Sections

Instruction files use HTML-style comments to mark managed sections:

```markdown
# Custom content here is preserved

<!-- BEGIN AGENT-INSTRUCTION -->
Generated content here is replaced on each build
<!-- END AGENT-INSTRUCTION -->

# More custom content is also preserved
```

**Important:**
- Content inside markers is **replaced** on each build
- Content outside markers is **preserved**
- Manual edits inside markers will be lost

### Typical Repository Layout

```
monorepo/
├── .agent-instruction/
│   ├── config.json
│   └── rules/
│       ├── global.json
│       ├── testing.json
│       ├── security.json
│       └── golang.json
├── CLAUDE.md              # Root-level instruction file
├── AGENTS.md              # Root-level instruction file
├── packages/
│   ├── api/
│   │   ├── agent-instruction.json
│   │   ├── CLAUDE.md
│   │   └── AGENTS.md
│   ├── lib/
│   │   ├── agent-instruction.json
│   │   ├── CLAUDE.md
│   │   └── AGENTS.md
│   └── worker/
│       ├── agent-instruction.json
│       ├── CLAUDE.md
│       └── AGENTS.md
└── .gitignore
```

## Claude Code Integration

A Claude Code skill is available for natural language interaction:

```bash
# Install the skill
cd agent-instruction
bash skill/install.sh
```

Then in Claude Code sessions:

```
You: "Add a rule that we always use explicit error returns"
Claude: [Executes agent-instruction add command]

You: "Build the instruction files"
Claude: [Executes agent-instruction build]
```

See `skill/agent-instruction/README.md` for details.

## Best Practices

### Rule Organization

- **global.json** - Cross-cutting concerns, project-wide standards
- **testing.json** - Test patterns, coverage requirements, mocking strategies
- **security.json** - Security requirements, authentication, authorization
- **[language].json** - Language-specific conventions (golang.json, typescript.json)

### Writing Effective Rules

**Good Rules:**
```json
{
  "heading": "Error Handling",
  "rule": "Always wrap errors with context using fmt.Errorf with %w verb",
  "references": ["https://go.dev/blog/go1.13-errors"]
}
```

**Avoid Vague Rules:**
```json
{
  "rule": "Write good code"  // Too vague
}
```

### Commit Strategy

```bash
# Commit rules and generated files together
git add .agent-instruction/ CLAUDE.md AGENTS.md packages/*/CLAUDE.md packages/*/AGENTS.md
git commit -m "Add [category] instruction rules"
```

### Maintenance

- Review rules quarterly to remove outdated ones
- Update rules when conventions change
- Use `--verbose` to verify rule composition
- Keep rule files focused on single concerns

## Troubleshooting

### Already Initialized Error

```
Error: already initialized: .agent-instruction directory exists
```

**Solution:** Use `build` command instead of `init`:
```bash
agent-instruction build
```

### Not Initialized Error

```
Error: not initialized: run 'agent-instruction init' first
```

**Solution:** Run init in repository root:
```bash
agent-instruction init
```

### Circular Import Detected

```
Error: circular import detected: global -> testing -> global
```

**Solution:** Remove circular import from rule files. Check imports with:
```bash
grep -r "imports" .agent-instruction/rules/
```

### Invalid Config

```
Error: invalid config: at least one framework is required
```

**Solution:** Edit `.agent-instruction/config.json` and add frameworks:
```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

### Permission Denied

```
Error: write CLAUDE.md: permission denied
```

**Solution:** Check file permissions:
```bash
chmod 644 CLAUDE.md
```

## Development

### Building from Source

```bash
git clone https://github.com/validkeys/agent-instruction.git
cd agent-instruction
go build -o agent-instruction ./cmd/agent-instruction
./agent-instruction --version
```

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/rules
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Run linter: `golangci-lint run`
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

- **Issues:** [github.com/validkeys/agent-instruction/issues](https://github.com/validkeys/agent-instruction/issues)
- **Documentation:** See `docs/` directory for additional resources
- **Examples:** See `examples/` directory for sample configurations

## Changelog

### v1.0.0 (2026-04-13)

- Initial release
- Core commands: init, add, build, list
- Support for Claude and Agents frameworks
- Import system for rule composition
- Parallel package processing
- Claude Code skill integration
- Comprehensive test coverage (84%+)

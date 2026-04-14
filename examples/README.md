# Examples

This directory contains working example configurations demonstrating how to use `agent-instruction` in different repository structures.

## Overview

Each example is a complete, runnable configuration that you can:
1. Copy to your own repository as a starting point
2. Test with `agent-instruction build` to see generated output
3. Study to learn patterns and best practices

All examples use realistic rule content based on common development standards rather than placeholder text.

## Examples

### 1. Simple Repository (`simple-repo/`)

**Use case:** Single-package repository with basic global rules

**Structure:**
```
simple-repo/
└── .agent-instruction/
    ├── config.json          # Basic config for Claude framework
    └── rules/
        └── global.json      # Global development guidelines
```

**Features demonstrated:**
- Minimal configuration setup
- Global rules applied to entire repository
- Single framework (Claude)
- Basic instruction structure with headings and rules

**What gets generated:**
- `CLAUDE.md` in repository root
- Contains all rules from `global.json`

**Try it:**
```bash
cd simple-repo
agent-instruction build
cat CLAUDE.md
```

**Customize it:**
Edit `.agent-instruction/rules/global.json` to add your own coding standards, then rebuild.

---

### 2. Monorepo (`monorepo/`)

**Use case:** Multi-package repository with shared rules and package-specific overrides

**Structure:**
```
monorepo/
├── .agent-instruction/
│   ├── config.json              # Config for multiple packages
│   └── rules/
│       ├── global.json          # Repository-wide standards
│       ├── testing.json         # Testing requirements
│       └── security.json        # Security guidelines
└── packages/
    └── api/
        └── agent-instruction.json    # API-specific rules + imports
```

**Features demonstrated:**
- Multi-package configuration
- Multiple rule files organized by concern
- Import system (package imports global rules)
- Both Claude and Agents frameworks
- Package-specific instructions with references
- Rule composition via imports

**What gets generated:**
- `CLAUDE.md` and `AGENTS.md` in repository root (with global rules only)
- `CLAUDE.md` and `AGENTS.md` in `packages/api/` (with global + API-specific rules)
- Same for `packages/web/` and `packages/shared/` (if they existed)

**Import chain:**
The `packages/api/agent-instruction.json` imports additional rule files:
```json
"imports": [
  "../../.agent-instruction/rules/testing.json",
  "../../.agent-instruction/rules/security.json"
]
```

This means the API package's instruction files will contain:
1. Global development standards (automatically included from global.json)
2. Testing standards (imported from testing.json)
3. Security requirements (imported from security.json)
4. API-specific instructions (defined in the package config)

**Note:** Package configs do NOT need to import `global.json` because the global rules are automatically included in all packages by the build system.

**Try it:**
```bash
cd monorepo
agent-instruction build
cat packages/api/CLAUDE.md  # See combined output
```

**Customize it:**
1. Add your packages to `.agent-instruction/config.json`
2. Create package-specific configs in each package directory
3. Import only the rule files relevant to each package
4. Add package-specific instructions as needed

---

## Key Concepts Illustrated

### 1. Rule Organization
Break rules into logical files:
- `global.json` - Repository-wide standards
- `testing.json` - Test requirements
- `security.json` - Security guidelines
- `performance.json` - Performance best practices
- `docs.json` - Documentation standards

### 2. Import System
Packages can import additional rule files beyond the global rules:
```json
{
  "title": "API Instructions",
  "imports": [
    "../../.agent-instruction/rules/testing.json",
    "../../.agent-instruction/rules/security.json"
  ],
  "instructions": [/* API-specific rules */]
}
```

**Important:** Do NOT import `global.json` from package configs. The build system automatically includes global rules in all packages. Only import additional rule files that are specific to the package's needs.

### 3. References
Link rules to actual files in your codebase:
```json
{
  "rule": "Follow error handling pattern",
  "references": [
    {
      "title": "Error Handler",
      "path": "src/middleware/errorHandler.ts"
    }
  ]
}
```

### 4. Multiple Frameworks
Generate instructions for different AI frameworks:
```json
{
  "frameworks": ["claude", "agents"]
}
```
This creates both `CLAUDE.md` and `AGENTS.md` files.

---

## Validation

Each example can be validated using the build command:

```bash
# Validate simple-repo
cd simple-repo
agent-instruction build

# Validate monorepo
cd monorepo
agent-instruction build
```

Successful builds confirm:
- Valid JSON syntax
- Required fields present
- Import paths resolve correctly
- Rule structure is correct

---

## Creating Your Own

### Starting from scratch:
```bash
agent-instruction init
```

### Starting from an example:
```bash
# Copy example to your repo
cp -r examples/monorepo/.agent-instruction /path/to/your/repo/

# Customize config.json with your packages
vim .agent-instruction/config.json

# Customize rules for your needs
vim .agent-instruction/rules/global.json

# Build and test
agent-instruction build
```

---

## Common Patterns

### Pattern 1: Shared + Specific
Have shared rules that all packages use, plus package-specific additions:
- Global rules in `.agent-instruction/rules/`
- Package configs import shared rules
- Package configs add package-specific rules

### Pattern 2: Feature-Based Rules
Organize rules by feature rather than by package:
- `testing.json` - Imported by all packages
- `api.json` - Imported by API packages only
- `ui.json` - Imported by UI packages only
- `database.json` - Imported by packages that access database

### Pattern 3: Environment-Specific
Different rules for different environments (optional advanced usage):
- Use separate rule files: `development.json`, `production.json`
- Import based on deployment target
- Override settings per environment

---

## Questions?

- Read the main [README](../README.md) for full documentation
- Check [command documentation](../docs/commands/) for all available commands
- File issues at https://github.com/validkeys/agent-instruction/issues

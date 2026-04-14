# Configuration Reference

Complete reference for agent-instruction configuration files.

## Overview

agent-instruction uses JSON configuration files to manage instruction generation:

- **config.json** - Main configuration (frameworks, packages)
- **rule files** - Instruction definitions (global.json, testing.json, etc.)
- **package configs** - Package-specific overrides (agent-instruction.json)

## config.json

Main configuration file located at `.agent-instruction/config.json`.

### Format

```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

### Fields

#### version

**Type:** `string`
**Required:** Yes
**Values:** `"1.0"`

Configuration format version. Currently only `"1.0"` is supported.

```json
{
  "version": "1.0"
}
```

#### frameworks

**Type:** `string[]`
**Required:** Yes
**Values:** `["claude"]`, `["agents"]`, or `["claude", "agents"]`

AI frameworks to generate instruction files for.

```json
{
  "frameworks": ["claude", "agents"]
}
```

Generates:
- `CLAUDE.md` when `"claude"` included
- `AGENTS.md` when `"agents"` included

#### packages

**Type:** `string[]`
**Required:** Yes

Package discovery configuration.

**Automatic Discovery:**
```json
{
  "packages": ["auto"]
}
```

Finds all directories containing `agent-instruction.json`.

**Specific Packages:**
```json
{
  "packages": ["api", "lib", "services/worker"]
}
```

Generates instruction files only in specified directories (relative to repo root).

**Root Only:**
```json
{
  "packages": []
}
```

Generates instruction files only in repository root, no package discovery.

### Examples

#### Single Package Project

```json
{
  "version": "1.0",
  "frameworks": ["claude"],
  "packages": []
}
```

#### Monorepo with Auto-Discovery

```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

#### Specific Packages

```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": [
    "packages/api",
    "packages/lib",
    "services/auth",
    "services/worker"
  ]
}
```

### Validation

The config is validated on load:

```bash
agent-instruction build
```

**Validation Rules:**
- `version` must be present and `"1.0"`
- `frameworks` must contain at least one valid framework
- Valid frameworks: `"claude"`, `"agents"`
- `packages` must be an array (can be empty)

**Error Examples:**

```
Error: invalid config: config version is required
Error: invalid config: at least one framework is required
Error: invalid config: invalid framework: xyz (must be 'claude' or 'agents')
```

## Rule Files

Rule files define instructions, located in `.agent-instruction/rules/`.

### Format

```json
{
  "title": "Rule Category Title",
  "imports": ["other-rule-file"],
  "instructions": [
    {
      "heading": "Section Heading",
      "rule": "Instruction text",
      "references": ["https://example.com"]
    }
  ]
}
```

### Fields

#### title

**Type:** `string`
**Required:** Yes

Display title for the rule file.

```json
{
  "title": "Security Guidelines"
}
```

#### imports

**Type:** `string[]`
**Required:** No
**Default:** `[]`

List of rule files to import (without `.json` extension).

```json
{
  "imports": ["global", "testing"]
}
```

See [Import System](#import-system) for details.

#### instructions

**Type:** `Instruction[]`
**Required:** Yes
**Default:** `[]`

Array of instruction objects.

### Instruction Object

#### heading

**Type:** `string`
**Required:** No

Optional section heading for the instruction.

```json
{
  "heading": "Error Handling"
}
```

Renders as:

```markdown
## Error Handling
```

#### rule

**Type:** `string`
**Required:** Yes

The instruction text content.

```json
{
  "rule": "Always wrap errors with context using fmt.Errorf with %w verb"
}
```

Supports multi-line text:

```json
{
  "rule": "When writing tests:\n- Use table-driven patterns\n- Test both paths\n- Mock dependencies"
}
```

#### references

**Type:** `string[]`
**Required:** No
**Default:** `[]`

List of reference URLs or file paths.

```json
{
  "references": [
    "https://go.dev/blog/error-handling-and-go",
    "docs/error-handling.md"
  ]
}
```

Renders as markdown links in generated files.

### Examples

#### Basic Rule File

```json
{
  "title": "Global Instructions",
  "instructions": [
    {
      "heading": "Code Style",
      "rule": "Follow project conventions and style guide"
    },
    {
      "heading": "Testing",
      "rule": "Write tests for all new features"
    }
  ]
}
```

#### With Imports

```json
{
  "title": "Backend Service Rules",
  "imports": ["global", "testing", "security"],
  "instructions": [
    {
      "heading": "API Design",
      "rule": "Use REST conventions for HTTP APIs",
      "references": ["docs/api-design.md"]
    }
  ]
}
```

#### With References

```json
{
  "title": "Security Guidelines",
  "instructions": [
    {
      "heading": "Password Hashing",
      "rule": "Use bcrypt with cost factor 12+",
      "references": [
        "https://pkg.go.dev/golang.org/x/crypto/bcrypt",
        "OWASP Password Storage Cheat Sheet"
      ]
    }
  ]
}
```

### Validation

Rule files are validated on load:

**Validation Rules:**
- `title` must be present and non-empty
- `instructions` must be an array
- Each instruction must have `rule` field
- `imports` must be valid rule file names
- No circular imports

**Error Examples:**

```
Error: invalid rule: title is required
Error: invalid rule: instruction rule is required
Error: circular import detected: global -> testing -> global
```

## Package Configuration

Package-specific configuration in `<package>/agent-instruction.json`.

### Format

Same format as rule files:

```json
{
  "title": "Package-Specific Instructions",
  "imports": ["global", "testing"],
  "instructions": [
    {
      "heading": "Package Rules",
      "rule": "Package-specific instruction"
    }
  ]
}
```

### Behavior

For each package with `agent-instruction.json`:

1. Loads global rules from `.agent-instruction/rules/global.json`
2. Loads package config from `<package>/agent-instruction.json`
3. Composes instructions (package rules after global rules)
4. Generates instruction files in package directory

### Example

```
.agent-instruction/
└── rules/
    ├── global.json
    └── testing.json

packages/
└── api/
    ├── agent-instruction.json
    ├── CLAUDE.md          # Generated
    └── AGENTS.md          # Generated
```

**packages/api/agent-instruction.json:**

```json
{
  "title": "API Service Instructions",
  "imports": ["global", "testing"],
  "instructions": [
    {
      "heading": "API Endpoints",
      "rule": "All endpoints must implement rate limiting"
    },
    {
      "heading": "OpenAPI",
      "rule": "Generate OpenAPI spec from code annotations"
    }
  ]
}
```

**Result:**

The generated `packages/api/CLAUDE.md` contains:
1. Instructions from `global.json`
2. Instructions from `testing.json` (imported by package config)
3. Package-specific instructions

## Import System

Rule files can import other rule files to compose instructions.

### Import Resolution

Imports are resolved recursively:

1. Load specified rule file
2. Resolve its imports recursively
3. Collect all instructions in order
4. Each file included only once
5. Detect and reject circular dependencies

### Import Order

Instructions maintain import order:

```json
{
  "imports": ["a", "b", "c"]
}
```

Final instruction order:
1. Instructions from `a.json`
2. Instructions from `b.json`
3. Instructions from `c.json`
4. Instructions from current file

### Circular Dependency Detection

Circular imports are detected and rejected:

**Example:**

```
global.json imports testing.json
testing.json imports security.json
security.json imports global.json  ← Circular!
```

**Error:**

```
Error: circular import detected: global -> testing -> security -> global
```

**Solution:** Remove circular reference.

### Import Examples

#### Linear Imports

**global.json:**
```json
{
  "title": "Global Rules",
  "instructions": [{"rule": "Rule 1"}]
}
```

**testing.json:**
```json
{
  "title": "Testing Rules",
  "imports": ["global"],
  "instructions": [{"rule": "Rule 2"}]
}
```

**security.json:**
```json
{
  "title": "Security Rules",
  "imports": ["testing"],
  "instructions": [{"rule": "Rule 3"}]
}
```

**Result:** Rule 1, Rule 2, Rule 3

#### Multiple Imports

**api.json:**
```json
{
  "title": "API Rules",
  "imports": ["global", "testing", "security"],
  "instructions": [{"rule": "Rule 4"}]
}
```

**Result:** Rule 1, Rule 2, Rule 3, Rule 4

#### Diamond Imports

```
    global
   /      \
testing  security
   \      /
     api
```

**Result:** global imported only once, then testing, then security, then api rules.

### Best Practices

**DO:**
- Keep import chains shallow (max 3-5 levels)
- Import only what you need
- Organize by concern (testing, security, language)
- Document import relationships

**DON'T:**
- Create circular imports
- Import everything into everything
- Create deep import chains (>5 levels)
- Duplicate content across files

## File Organization

### Recommended Structure

```
.agent-instruction/
├── config.json
└── rules/
    ├── global.json          # Cross-cutting concerns
    ├── testing.json         # Test patterns
    ├── security.json        # Security requirements
    ├── golang.json          # Go-specific rules
    ├── typescript.json      # TypeScript rules
    └── api-design.json      # API conventions
```

### Naming Conventions

**Rule Files:**
- Use lowercase
- Use hyphens for multi-word names: `api-design.json`
- Be specific: `golang.json` not `code.json`
- Group by concern: `testing.json`, `security.json`

**Package Configs:**
- Always named `agent-instruction.json`
- Located in package root directory

## Generated Files

### Managed Sections

Generated instruction files contain managed sections:

```markdown
# Custom Header

Your custom content here.

<!-- BEGIN AGENT-INSTRUCTION -->
# Generated Content

This section is managed by agent-instruction.
<!-- END AGENT-INSTRUCTION -->

# Custom Footer

More custom content.
```

### Markers

**Begin Marker:**
```html
<!-- BEGIN AGENT-INSTRUCTION -->
```

**End Marker:**
```html
<!-- END AGENT-INSTRUCTION -->
```

### Content Preservation

- Content **outside markers**: Preserved on rebuild
- Content **inside markers**: Replaced on rebuild
- Manual edits inside markers: **Lost** on rebuild

**Best Practice:** Keep custom content outside markers.

### First Build

On first build of existing file:
- No markers present
- All content wrapped in markers
- Content preserved
- Markers added for future rebuilds

## Best Practices

### Configuration

- **Start simple**: Use `"packages": ["auto"]` initially
- **Be explicit**: Switch to specific package list when needed
- **Version control**: Commit `.agent-instruction/` directory
- **Document changes**: Update rules when conventions change

### Rule Files

- **Single responsibility**: One concern per file
- **Clear titles**: Descriptive file and instruction headings
- **Specific rules**: Actionable, not vague
- **References**: Link to documentation
- **Examples**: Include code examples in rules

### Imports

- **Minimal imports**: Import only needed files
- **Logical order**: Common rules first, specific rules last
- **No cycles**: Keep import graph acyclic
- **Document relationships**: Comment why imports exist

### Package Configs

- **Inherit from global**: Import global rules
- **Package-specific only**: Don't duplicate global rules
- **Keep focused**: Only package-specific concerns
- **Document purpose**: Explain why package differs

## Troubleshooting

See individual command documentation:
- [init troubleshooting](commands/init.md#troubleshooting)
- [build troubleshooting](commands/build.md#troubleshooting)
- [add troubleshooting](commands/add.md#troubleshooting)
- [list troubleshooting](commands/list.md#troubleshooting)

## See Also

- [Main README](../README.md)
- [Command Reference](commands/)
- [Examples](examples.md)

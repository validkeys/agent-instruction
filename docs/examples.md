# Examples

Complete workflow examples for common use cases.

## Quick Start

### Single Package Project

Initialize and set up agent-instruction for a single package repository.

```bash
# 1. Initialize in project root
cd /path/to/project
agent-instruction init --non-interactive --packages ""

# 2. Add your first rule
agent-instruction add "Follow Go best practices and conventions" \
  --title "Code Style" \
  --rule global

# 3. Build instruction files
agent-instruction build

# 4. Verify generated files
ls -la CLAUDE.md AGENTS.md

# 5. Commit
git add .agent-instruction/ CLAUDE.md AGENTS.md
git commit -m "Initialize agent-instruction"
```

### Monorepo Project

Set up for monorepo with automatic package discovery.

```bash
# 1. Initialize at repo root
cd /path/to/monorepo
agent-instruction init --non-interactive

# 2. Add global rules
agent-instruction add "Use explicit error handling" --rule global
agent-instruction add "Write table-driven tests" --rule testing

# 3. Build all packages
agent-instruction build --verbose

# 4. Verify structure
tree -L 2 -I 'node_modules'

# 5. Commit
git add .agent-instruction/ **/CLAUDE.md **/AGENTS.md
git commit -m "Set up agent-instruction for monorepo"
```

## Complete Workflows

### Workflow 1: Initial Setup

Complete setup from scratch with custom rules.

```bash
# Step 1: Initialize
cd my-project
agent-instruction init

# Answer prompts:
# - Backup existing files? Yes
# - Frameworks? claude, agents
# - Packages? auto

# Step 2: Create category rule files
agent-instruction add "Follow project coding standards" \
  --title "Standards" \
  --rule global

agent-instruction add "Use bcrypt for passwords" \
  --title "Password Security" \
  --rule security

agent-instruction add "Mock external dependencies" \
  --title "Test Isolation" \
  --rule testing

# Step 3: Build
agent-instruction build --verbose

# Step 4: Review generated files
cat CLAUDE.md

# Step 5: Commit
git add .agent-instruction/ CLAUDE.md AGENTS.md
git commit -m "Initialize agent-instruction with initial rules"
```

### Workflow 2: Adding Rules During Development

Capture patterns as you discover them.

```bash
# While coding, discover a pattern
# Add it immediately:
agent-instruction add "Use context.Context for cancellation in long-running operations" \
  --title "Concurrency Patterns" \
  --rule golang

# Continue working...

# Discover another pattern
agent-instruction add "Validate all user input at API boundaries" \
  --title "Input Validation" \
  --rule security

# At end of session, rebuild
agent-instruction build

# Review changes
git diff CLAUDE.md

# Commit with your code changes
git add .agent-instruction/ CLAUDE.md AGENTS.md
git commit -m "Add concurrency and security patterns"
```

### Workflow 3: Organizing Rules

Restructure rules into logical categories.

```bash
# Step 1: List current rules
agent-instruction list --verbose

# Step 2: Create new category files
agent-instruction add "Use semantic versioning for releases" \
  --title "Versioning" \
  --rule deployment

agent-instruction add "Run migrations before deploying" \
  --title "Database Migrations" \
  --rule deployment

agent-instruction add "Use feature flags for gradual rollouts" \
  --title "Feature Flags" \
  --rule deployment

# Step 3: Review organization
agent-instruction list

# Step 4: Build and commit
agent-instruction build
git add .agent-instruction/ CLAUDE.md AGENTS.md
git commit -m "Organize deployment rules"
```

### Workflow 4: Package-Specific Rules

Add rules that apply only to specific packages.

```bash
# Navigate to package
cd packages/api

# Create package config
cat > agent-instruction.json <<EOF
{
  "title": "API Service Instructions",
  "imports": ["global", "security"],
  "instructions": [
    {
      "heading": "Rate Limiting",
      "rule": "All public endpoints must implement rate limiting"
    },
    {
      "heading": "OpenAPI",
      "rule": "Generate OpenAPI spec from code, never write by hand"
    }
  ]
}
EOF

# Build from root (processes all packages)
cd ../..
agent-instruction build --verbose

# Verify package-specific file
cat packages/api/CLAUDE.md

# Commit
git add packages/api/agent-instruction.json packages/api/CLAUDE.md
git commit -m "Add API-specific instruction rules"
```

## Configuration Examples

### Example 1: Simple Repository

**.agent-instruction/config.json:**
```json
{
  "version": "1.0",
  "frameworks": ["claude"],
  "packages": []
}
```

**.agent-instruction/rules/global.json:**
```json
{
  "title": "Project Instructions",
  "instructions": [
    {
      "heading": "Code Style",
      "rule": "Follow Go best practices. Use gofmt and golint."
    },
    {
      "heading": "Testing",
      "rule": "Write table-driven tests. Achieve 80%+ coverage."
    },
    {
      "heading": "Documentation",
      "rule": "Document all exported functions and types."
    }
  ]
}
```

**Result:** Single `CLAUDE.md` file in repository root.

### Example 2: Monorepo with Multiple Rule Files

**.agent-instruction/config.json:**
```json
{
  "version": "1.0",
  "frameworks": ["claude", "agents"],
  "packages": ["auto"]
}
```

**.agent-instruction/rules/global.json:**
```json
{
  "title": "Global Instructions",
  "instructions": [
    {
      "heading": "Project Overview",
      "rule": "This is a TypeScript monorepo using pnpm workspaces"
    },
    {
      "heading": "Code Standards",
      "rule": "Use ESLint and Prettier for code formatting"
    }
  ]
}
```

**.agent-instruction/rules/testing.json:**
```json
{
  "title": "Testing Standards",
  "imports": ["global"],
  "instructions": [
    {
      "heading": "Test Framework",
      "rule": "Use Vitest for unit tests, Playwright for E2E"
    },
    {
      "heading": "Coverage",
      "rule": "Maintain 80%+ test coverage on all packages"
    },
    {
      "heading": "Test Structure",
      "rule": "Group tests by feature, use describe/it blocks"
    }
  ]
}
```

**.agent-instruction/rules/security.json:**
```json
{
  "title": "Security Guidelines",
  "imports": ["global"],
  "instructions": [
    {
      "heading": "Authentication",
      "rule": "Use JWT tokens with 15-minute expiry and refresh tokens",
      "references": ["docs/auth-architecture.md"]
    },
    {
      "heading": "Input Validation",
      "rule": "Use Zod for runtime type validation on all API inputs"
    },
    {
      "heading": "SQL Injection",
      "rule": "Always use parameterized queries, never string interpolation"
    }
  ]
}
```

**Result:** `CLAUDE.md` and `AGENTS.md` in root and each discovered package.

### Example 3: Package-Specific Configuration

**packages/web-app/agent-instruction.json:**
```json
{
  "title": "Web App Instructions",
  "imports": ["global", "testing"],
  "instructions": [
    {
      "heading": "React Patterns",
      "rule": "Use functional components with hooks, avoid class components"
    },
    {
      "heading": "State Management",
      "rule": "Use Zustand for global state, React hooks for local state"
    },
    {
      "heading": "Performance",
      "rule": "Lazy load routes and heavy components. Use React.memo for expensive renders"
    }
  ]
}
```

**packages/api-server/agent-instruction.json:**
```json
{
  "title": "API Server Instructions",
  "imports": ["global", "testing", "security"],
  "instructions": [
    {
      "heading": "API Design",
      "rule": "Use RESTful conventions. Implement HATEOAS for discoverability"
    },
    {
      "heading": "Error Handling",
      "rule": "Return RFC 7807 Problem Details for all errors"
    },
    {
      "heading": "Rate Limiting",
      "rule": "Implement rate limiting with Redis: 100 req/min per user, 1000 req/min per IP"
    }
  ]
}
```

**Result:** Each package gets its own `CLAUDE.md`/`AGENTS.md` with composed instructions.

## Import Chain Examples

### Example 1: Linear Import Chain

```
global.json → testing.json → api.json
```

**global.json:**
```json
{
  "title": "Global Rules",
  "instructions": [
    {"heading": "Style", "rule": "Use consistent formatting"}
  ]
}
```

**testing.json:**
```json
{
  "title": "Testing Rules",
  "imports": ["global"],
  "instructions": [
    {"heading": "Coverage", "rule": "Maintain 80%+ coverage"}
  ]
}
```

**api.json:**
```json
{
  "title": "API Rules",
  "imports": ["testing"],
  "instructions": [
    {"heading": "REST", "rule": "Follow REST conventions"}
  ]
}
```

**Result Order:**
1. Style (from global)
2. Coverage (from testing)
3. REST (from api)

### Example 2: Multiple Imports

```
      global
     /      \
testing    security
     \      /
       api
```

**api.json:**
```json
{
  "title": "API Rules",
  "imports": ["testing", "security"],
  "instructions": [
    {"heading": "API Design", "rule": "Use REST conventions"}
  ]
}
```

**Result Order:**
1. Global rules (imported by testing)
2. Testing rules
3. Security rules
4. API rules

### Example 3: Avoiding Circular Imports

**❌ Wrong (Circular):**

**global.json:**
```json
{
  "imports": ["testing"]
}
```

**testing.json:**
```json
{
  "imports": ["global"]
}
```

**Error:** `circular import detected: global -> testing -> global`

**✓ Correct (Hierarchical):**

**global.json:**
```json
{
  "title": "Global Rules",
  "instructions": [...]
}
```

**testing.json:**
```json
{
  "title": "Testing Rules",
  "imports": ["global"],
  "instructions": [...]
}
```

## Use Case Scenarios

### Scenario 1: Converting Existing Repository

You have existing `CLAUDE.md` with 50+ lines of instructions.

```bash
# 1. Initialize (creates backups)
agent-instruction init --non-interactive

# 2. Review backup
cat CLAUDE.md.backup

# 3. Migrate content by category
# Extract testing rules
agent-instruction add "Write comprehensive tests..." --rule testing

# Extract security rules
agent-instruction add "Validate all inputs..." --rule security

# Extract general rules
agent-instruction add "Follow coding standards..." --rule global

# 4. Build
agent-instruction build

# 5. Compare
diff CLAUDE.md.backup CLAUDE.md

# 6. Adjust as needed, then commit
git add .agent-instruction/ CLAUDE.md
git commit -m "Migrate to agent-instruction"
```

### Scenario 2: Team Onboarding

New team member needs to understand project rules.

```bash
# Show all rules
agent-instruction list --verbose

# Build fresh instruction files
agent-instruction build

# Open in editor for review
code CLAUDE.md

# Generate documentation
agent-instruction list --verbose > docs/ai-instructions.md
```

### Scenario 3: Rule Evolution

Project conventions change over time.

```bash
# Review current rules
agent-instruction list --verbose

# Edit rule files to update conventions
vim .agent-instruction/rules/testing.json

# Or add new rules
agent-instruction add "New pattern discovered..." --rule patterns

# Rebuild
agent-instruction build

# Review changes
git diff CLAUDE.md

# Commit
git add .agent-instruction/ CLAUDE.md
git commit -m "Update testing conventions"
```

### Scenario 4: Multi-Framework Support

Support both Claude and other AI agents.

```bash
# Initialize with both frameworks
agent-instruction init \
  --non-interactive \
  --frameworks claude,agents

# Add rules (applies to both)
agent-instruction add "Project-specific rules..." --rule global

# Build both
agent-instruction build

# Verify both files
ls -la CLAUDE.md AGENTS.md

# Both files have same content, just different names
diff CLAUDE.md AGENTS.md
# (Should be identical)
```

## Advanced Patterns

### Pattern 1: Modular Rule Organization

```
.agent-instruction/rules/
├── global.json          # Core standards
├── golang.json          # Go-specific
├── typescript.json      # TS-specific
├── testing.json         # Test patterns
├── security.json        # Security rules
├── performance.json     # Performance guidelines
├── deployment.json      # Deploy procedures
└── api-design.json      # API conventions
```

Each package imports only relevant files:

```json
{
  "imports": ["global", "golang", "testing", "security"]
}
```

### Pattern 2: Tiered Rule System

**Tier 1: Universal Rules (global.json)**
- Code style
- Documentation standards
- Version control practices

**Tier 2: Domain Rules**
- testing.json
- security.json
- performance.json

**Tier 3: Technology Rules**
- golang.json
- typescript.json
- react.json

**Tier 4: Package-Specific**
- api/agent-instruction.json
- web/agent-instruction.json

### Pattern 3: Progressive Enhancement

Start minimal, add over time:

**Week 1:**
```json
{
  "title": "Basic Rules",
  "instructions": [
    {"rule": "Follow existing code style"}
  ]
}
```

**Week 2:**
```json
{
  "title": "Basic Rules",
  "instructions": [
    {"heading": "Style", "rule": "Use gofmt..."},
    {"heading": "Testing", "rule": "Write tests..."}
  ]
}
```

**Month 2:**
```json
{
  "title": "Comprehensive Rules",
  "imports": ["testing", "security"],
  "instructions": [
    {"heading": "Style", "rule": "Use gofmt, golint..."},
    {"heading": "Architecture", "rule": "Follow hexagonal..."}
  ]
}
```

## Integration Examples

### CI/CD Integration

**.github/workflows/ci.yml:**
```yaml
jobs:
  verify-instructions:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Install agent-instruction
        run: go install github.com/validkeys/agent-instruction@latest

      - name: Build instruction files
        run: agent-instruction build

      - name: Verify files are up-to-date
        run: |
          if git diff --exit-code; then
            echo "✓ Instruction files up to date"
          else
            echo "✗ Instruction files out of date"
            echo "Run: agent-instruction build"
            exit 1
          fi
```

### Pre-commit Hook

**.git/hooks/pre-commit:**
```bash
#!/bin/bash

# Rebuild instruction files
agent-instruction build

# Stage any changes
git add **/CLAUDE.md **/AGENTS.md

echo "✓ Instruction files updated"
```

### Makefile Integration

**Makefile:**
```makefile
.PHONY: rules rules-list rules-build

rules-list:
	@agent-instruction list --verbose

rules-build:
	@agent-instruction build --verbose

rules-verify:
	@agent-instruction build --dry-run

# Include in default build
build: rules-build
	@go build ./...
```

## See Also

- [Configuration Reference](configuration.md)
- [Command Reference](commands/)
- [Main README](../README.md)

# Example: Complete Build Workflow

## Scenario

You've been working on a new feature and discovered several patterns that should become rules. Now you want to add them all, review what you've added, and propagate the changes to instruction files.

## Step 1: Add Multiple Rules

During development, you add rules as you discover patterns:

```bash
# Testing pattern discovered
agent-instruction add "Use table-driven tests for testing multiple inputs" \
  --title="Table-Driven Tests" \
  --rule="testing"
```

```
✓ Rule added to testing.json
```

```bash
# Security requirement identified
agent-instruction add "Validate all user input before passing to database queries. Use parameterized queries to prevent SQL injection" \
  --title="Input Validation" \
  --rule="security"
```

```
✓ Rule added to security.json
```

```bash
# Go convention clarified
agent-instruction add "Interface names should be short and end with 'er' (e.g., Reader, Writer, Handler)" \
  --title="Interface Naming" \
  --rule="golang"
```

```
✓ Rule added to golang.json
```

## Step 2: Review What Was Added

Before building, list all rules to review:

```bash
agent-instruction list
```

**Output:**

```
Rule Files (3):

1. testing.json (1 rule)
   - Table-Driven Tests

2. security.json (1 rule)
   - Input Validation

3. golang.json (1 rule)
   - Interface Naming

Run with --verbose to see full content
```

Verify details with verbose flag:

```bash
agent-instruction list --verbose
```

**Output:**

```
Rule Files (3):

1. testing.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Table-Driven Tests
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Use table-driven tests for testing multiple inputs

2. security.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Input Validation
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Validate all user input before passing to database queries.
   Use parameterized queries to prevent SQL injection

3. golang.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Interface Naming
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Interface names should be short and end with 'er'
   (e.g., Reader, Writer, Handler)
```

## Step 3: Build Instruction Files

Propagate rules to CLAUDE.md and AGENTS.md:

```bash
agent-instruction build
```

**Output:**

```
Building instruction files...

✓ Read 3 rule files
✓ Generated CLAUDE.md (1,234 bytes)
✓ Generated AGENTS.md (987 bytes)

Instruction files updated successfully.
```

## Step 4: Verify Changes

Check what changed in the instruction files:

```bash
git diff CLAUDE.md
```

**Output:**

```diff
diff --git a/CLAUDE.md b/CLAUDE.md
index abc123..def456 100644
--- a/CLAUDE.md
+++ b/CLAUDE.md
@@ -45,6 +45,26 @@ When writing tests:

 - Write tests before implementation (TDD)
 - Mock external dependencies
+
+## Table-Driven Tests
+
+Use table-driven tests for testing multiple inputs
+
+## Security
+
+### Input Validation
+
+Validate all user input before passing to database queries.
+Use parameterized queries to prevent SQL injection
+
+## Go Guidelines
+
+### Interface Naming
+
+Interface names should be short and end with 'er'
+(e.g., Reader, Writer, Handler)
```

## Step 5: Commit Changes

```bash
# Stage all changes
git add .agent-instruction/ CLAUDE.md AGENTS.md

# Commit with descriptive message
git commit -m "docs: Add testing, security, and Go naming rules"
```

**Output:**

```
[main a1b2c3d] docs: Add testing, security, and Go naming rules
 5 files changed, 45 insertions(+)
```

## Summary

This workflow demonstrates the complete cycle:

1. ✓ Add rules as patterns are discovered
2. ✓ List to review before building
3. ✓ Build to propagate changes
4. ✓ Verify with git diff
5. ✓ Commit with clear message

## Best Practices

- **Add incrementally**: Add rules as you discover them, don't batch up
- **Review before building**: Use `list` to check what's been added
- **Verify changes**: Always `git diff` before committing
- **Descriptive commits**: Explain what rules were added and why

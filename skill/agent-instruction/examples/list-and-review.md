# Example: Listing and Reviewing Rules

## Scenario

Before adding a new rule, you want to see what rules already exist to avoid creating duplicates. You also want to review existing rules to ensure they're still relevant.

## Basic List Command

Show all rule files with summaries:

```bash
agent-instruction list
```

**Output:**

```
Rule Files (4):

1. global.json (2 rules)
   - Code Quality
   - Documentation Standards

2. testing.json (3 rules)
   - Table-Driven Tests
   - Test Naming
   - Mock External Dependencies

3. security.json (2 rules)
   - Input Validation
   - Password Hashing

4. golang.json (4 rules)
   - Error Handling
   - Interface Naming
   - Package Organization
   - Receiver Names

Total: 11 rules across 4 files

Run with --verbose to see full content
```

## Verbose List

Show full rule content with formatting:

```bash
agent-instruction list --verbose
```

**Output:**

```
Rule Files (4):

1. global.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Code Quality
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Keep functions small and focused. Each function
   should do one thing well. Prefer composition over
   large monolithic functions.

   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Documentation Standards
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Document public APIs with clear examples. Include
   usage examples in godoc comments for exported
   functions and types.

2. testing.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Table-Driven Tests
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Use table-driven tests for testing multiple inputs.
   Example:

   tests := []struct {
     name string
     input int
     want string
   }{
     {"zero", 0, "zero"},
     {"one", 1, "one"},
   }

   [... more rules ...]
```

## Understanding the Output

### Summary View (`list`)

- Shows count of rules per file
- Lists rule titles for quick scanning
- Fast overview of what's configured
- Good for checking if a category exists

### Detailed View (`list --verbose`)

- Shows full rule content
- Includes all text and examples
- Visual separators for readability
- Good for reviewing specific rules

## Use Cases

### 1. Check Before Adding

Search for keywords to avoid duplicates:

```bash
# Check if error handling rule exists
agent-instruction list --verbose | grep -i "error"
```

**Output:**

```
   Error Handling
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Always wrap errors with fmt.Errorf to add context.
   Example: return fmt.Errorf("failed to read: %w", err)
```

Found it! No need to add a duplicate.

### 2. Review Category Contents

See what's in a specific rule file:

```bash
# Show all security rules
agent-instruction list --verbose | grep -A 10 "security.json"
```

**Output:**

```
2. security.json
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Input Validation
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Validate all user input before passing to database
   queries. Use parameterized queries to prevent SQL
   injection

   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Password Hashing
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Always use bcrypt for password hashing. Never store
   passwords in plain text or use weak hashing like MD5.
```

### 3. Count Rules by Category

Quick stats on rule distribution:

```bash
agent-instruction list | grep "rules)"
```

**Output:**

```
1. global.json (2 rules)
2. testing.json (3 rules)
3. security.json (2 rules)
4. golang.json (4 rules)
```

Identifies which categories need more coverage.

### 4. Export for Documentation

Generate a reference document:

```bash
# Export all rules to a file
agent-instruction list --verbose > docs/instruction-rules.txt
```

Good for team documentation or onboarding.

## When to Use --verbose

### Use verbose when:

- Adding a new rule (check for duplicates)
- Reviewing rule content for accuracy
- Sharing rules with team members
- Debugging instruction file issues
- Writing documentation

### Use summary when:

- Quick check of what exists
- Counting rules per category
- Seeing rule titles at a glance
- Checking if a category exists

## Best Practices

**Before adding rules:**
```bash
# Always list first to avoid duplicates
agent-instruction list --verbose | grep -i "keyword"

# If not found, safe to add
agent-instruction add "New rule..." --rule="category"
```

**Periodic review:**
```bash
# Monthly review of all rules
agent-instruction list --verbose > review-$(date +%Y-%m).txt

# Discuss with team, remove outdated rules
```

**Team onboarding:**
```bash
# Show new developers what's configured
agent-instruction list --verbose

# Export to shared docs
agent-instruction list --verbose > docs/current-rules.md
```

## Troubleshooting

### No rules listed

```bash
agent-instruction list
```

```
No rule files found in .agent-instruction/rules/

Initialize with: agent-instruction init
```

**Solution:** Run `agent-instruction init` to create structure.

### Empty rule file

```
1. testing.json (0 rules)
```

**Solution:** Rule file exists but has no rules. Add some or delete the file.

### Malformed JSON error

```
Error reading golang.json: invalid JSON syntax at line 12
```

**Solution:** Check JSON syntax in the rule file. Fix manually or remove and recreate.

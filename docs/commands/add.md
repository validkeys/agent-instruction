# add Command

Add a new instruction rule to a rule file.

## Synopsis

```bash
agent-instruction add <rule-content> [flags]
```

## Description

The `add` command appends a new instruction to an existing rule file in `.agent-instruction/rules/`. This is the fastest way to add rules during development when you discover patterns that should be codified.

After adding rules, run `build` to regenerate instruction files.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `<rule-content>` | Yes | The instruction text to add |

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | - | Optional heading for the rule section |
| `--rule` | string | - | Target rule file name (without .json extension) |

## Behavior

### Rule File Selection

**With `--rule` flag:**
- Adds to specified rule file
- Creates file if it doesn't exist
- Example: `--rule global` adds to `global.json`

**Without `--rule` flag:**
- Lists available rule files
- Prompts for interactive selection
- User chooses from menu

### Rule Content

The rule content is added as a new instruction:

```json
{
  "heading": "Title if --title provided",
  "rule": "Your rule content here"
}
```

### File Updates

1. Loads existing rule file
2. Appends new instruction to `instructions` array
3. Validates JSON structure
4. Writes updated file
5. Prints success message with next steps

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success - rule added |
| `1` | Not initialized (run `init` first) |
| `1` | Empty rule content |
| `1` | Invalid rule file name |
| `1` | File write error |
| `1` | JSON validation error |

## Examples

### Quick Add to Global Rules

```bash
agent-instruction add "Always use explicit error handling in Go" --rule global
```

### Add with Title

```bash
agent-instruction add "Use bcrypt for password hashing with cost factor 12+" \
  --title "Password Security" \
  --rule security
```

### Interactive File Selection

```bash
agent-instruction add "Always validate user input before processing"
```

Prompts:
```
Select rule file:
  1. global
  2. testing
  3. security
Enter number: 3
✓ Added instruction to security.json
```

### Multi-Line Content

```bash
agent-instruction add "When writing tests:
- Use table-driven patterns
- Test both success and error paths
- Mock external dependencies
- Cover edge cases" --title "Testing Guidelines" --rule testing
```

### Add and Rebuild

```bash
# Add rule
agent-instruction add "Document all exported functions" --rule global

# Rebuild immediately
agent-instruction build
```

### Create New Rule File

If the rule file doesn't exist, it's created:

```bash
# Creates golang.json if it doesn't exist
agent-instruction add "Use explicit types instead of interface{}" \
  --title "Type Safety" \
  --rule golang
```

## Common Patterns

### During Development Session

When AI makes repeated mistakes:

```bash
# Discovered pattern
agent-instruction add "Never use global variables for state management" \
  --title "State Management" \
  --rule global

# Apply immediately
agent-instruction build
```

### Category Organization

Group related rules:

```bash
# Testing rules
agent-instruction add "Write tests before implementation" --rule testing
agent-instruction add "Mock external dependencies" --rule testing
agent-instruction add "Achieve 80%+ coverage" --rule testing

# Security rules
agent-instruction add "Validate all user input" --rule security
agent-instruction add "Use parameterized queries" --rule security
agent-instruction add "Sanitize error messages" --rule security

# Rebuild once
agent-instruction build
```

### Language-Specific Rules

```bash
# Go conventions
agent-instruction add "Use camelCase for unexported names" --rule golang

# TypeScript conventions
agent-instruction add "Enable strict mode in tsconfig.json" --rule typescript

# Python conventions
agent-instruction add "Follow PEP 8 style guide" --rule python
```

### Quick Capture Workflow

```bash
# 1. Discover pattern during development
# 2. Add rule immediately
agent-instruction add "Your insight here" --rule category

# 3. Continue working (build later)
# 4. At end of session, rebuild all
agent-instruction build

# 5. Commit with changes
git add .agent-instruction/ **/CLAUDE.md
git commit -m "Add new instruction rules"
```

## Rule Content Guidelines

### Good Rule Examples

**Specific and Actionable:**
```bash
agent-instruction add "Use fmt.Errorf with %w verb to wrap errors, preserving stack traces" \
  --title "Error Wrapping" \
  --rule golang
```

**With Context:**
```bash
agent-instruction add "All HTTP endpoints must implement rate limiting using golang.org/x/time/rate" \
  --title "Rate Limiting" \
  --rule security
```

**Clear Constraints:**
```bash
agent-instruction add "Table-driven tests required for functions with 3+ code paths" \
  --title "Test Coverage" \
  --rule testing
```

### Avoid Vague Rules

**Too General:**
```bash
agent-instruction add "Write good code" --rule global  # Not helpful
```

**Better:**
```bash
agent-instruction add "Functions should do one thing. Split functions over 50 lines" \
  --title "Function Size" \
  --rule global
```

## Troubleshooting

### Not Initialized

**Error:**
```
Error: not initialized: run 'agent-instruction init' first
```

**Solution:**
Initialize first:

```bash
agent-instruction init
```

### Empty Rule Content

**Error:**
```
Error: rule content cannot be empty
```

**Solution:**
Provide rule content as argument:

```bash
agent-instruction add "Your rule here" --rule global
```

### Invalid Rule File

**Error:**
```
Error: add instruction: invalid rule: title is required
```

**Solution:**
The rule file is corrupted. Check JSON syntax:

```bash
# Validate JSON
cat .agent-instruction/rules/global.json | jq .

# Fix manually if needed
vim .agent-instruction/rules/global.json
```

### Permission Denied

**Error:**
```
Error: add instruction: write rule file: permission denied
```

**Solution:**
Check file permissions:

```bash
# Check permissions
ls -la .agent-instruction/rules/

# Fix if needed
chmod u+w .agent-instruction/rules/global.json
```

### No Rule Files Found

**Error:**
```
Error: no rule files found in .agent-instruction/rules
```

**Solution:**
Create at least one rule file:

```bash
cat > .agent-instruction/rules/global.json <<EOF
{
  "title": "Global Instructions",
  "instructions": []
}
EOF
```

## Rule File Structure

After adding, the rule file contains:

```json
{
  "title": "Category Name",
  "instructions": [
    {
      "heading": "Previous Rule",
      "rule": "Previous content"
    },
    {
      "heading": "Your Title",
      "rule": "Your new rule content"
    }
  ]
}
```

## Next Steps

After adding rules:

1. **Review**: Check the rule file was updated correctly
   ```bash
   cat .agent-instruction/rules/[file].json
   ```

2. **Build**: Regenerate instruction files
   ```bash
   agent-instruction build
   ```

3. **Verify**: Check generated files
   ```bash
   git diff CLAUDE.md
   ```

4. **Commit**: Save changes
   ```bash
   git add .agent-instruction/ CLAUDE.md AGENTS.md
   git commit -m "Add [category] instruction rules"
   ```

## Related Commands

- [`build`](build.md) - Regenerate files after adding rules
- [`list`](list.md) - View all rules before adding
- [`init`](init.md) - Initialize before adding rules

## See Also

- [Configuration Documentation](../configuration.md)
- [Rule File Format](../configuration.md#rule-files)
- [Examples](../examples.md)
- [Main README](../../README.md)

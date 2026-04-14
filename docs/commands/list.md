# list Command

Display all instruction rules.

## Synopsis

```bash
agent-instruction list [flags]
```

## Description

The `list` command displays all rule files and their contents from `.agent-instruction/rules/` in a readable format. Use this to review existing rules before adding new ones or to understand the current instruction set.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--verbose` | boolean | `false` | Show full rule content with all details |

## Behavior

### Default Output (Summary)

Shows overview of each rule file:
- File name
- Number of instructions
- Instruction headings (if present)

### Verbose Output

Shows complete details:
- File name
- Title
- Import count
- Each instruction with:
  - Heading
  - Full rule text
  - Reference count

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Not initialized (run `init` first) |
| `1` | No rule files found |
| `1` | File read error |

## Examples

### Basic List

```bash
agent-instruction list
```

Output:
```
📄 global.json
   3 instruction(s)
   - Code Standards
   - Testing Practices
   - Error Handling

📄 security.json
   2 instruction(s)
   - Authentication
   - Input Validation

📄 testing.json
   4 instruction(s)
   - Test Structure
   - Coverage Requirements
   - Mocking Strategy
   - Edge Cases
```

### Verbose Output

```bash
agent-instruction list --verbose
```

Output:
```
📄 global.json
   Title: Global Instructions
   Imports: 0

   [1] Code Standards
       Follow Go best practices. Use gofmt, enable linters, write idiomatic code.

   [2] Testing Practices
       Write table-driven tests. Test both success and error paths.
       References: 1

   [3] Error Handling
       Always wrap errors with context using fmt.Errorf with %w verb.

📄 security.json
   Title: Security Guidelines
   Imports: 1

   [1] Authentication
       Use bcrypt for password hashing with cost factor 12+.
       References: 1

   [2] Input Validation
       Validate and sanitize all user input before processing.
```

### Filter with grep

```bash
# Find specific rules
agent-instruction list --verbose | grep -A 5 "Testing"

# Count total instructions
agent-instruction list | grep "instruction(s)" | wc -l

# Check if rule exists
agent-instruction list --verbose | grep -i "bcrypt"
```

### Preview Before Adding

```bash
# Check if rule already exists
agent-instruction list --verbose | grep -i "error handling"

# Add only if not duplicate
agent-instruction add "New error handling rule" --rule global
```

## Common Patterns

### Review Before Changes

```bash
# See current state
agent-instruction list

# Make changes
agent-instruction add "New rule" --rule global

# Verify changes
agent-instruction list --verbose
```

### Documentation

```bash
# Generate rule documentation
agent-instruction list --verbose > docs/current-rules.txt

# Share with team
cat docs/current-rules.txt
```

### Audit Rules

```bash
# Review all rules in detail
agent-instruction list --verbose

# Look for inconsistencies, duplicates, or outdated rules
# Edit rule files as needed
vim .agent-instruction/rules/global.json

# Rebuild
agent-instruction build
```

### Find Rule Location

```bash
# Where is a specific rule?
agent-instruction list --verbose | grep -B 3 "password"
```

Output shows file name and rule number:
```
📄 security.json
   Title: Security Guidelines

   [1] Authentication
       Use bcrypt for password hashing with cost factor 12+.
```

### Check Imports

```bash
agent-instruction list --verbose | grep "Imports:"
```

Shows which rule files import others:
```
   Imports: 0
   Imports: 1
   Imports: 2
```

## Output Format

### Summary Format

```
📄 [filename]
   [count] instruction(s)
   - [heading 1]
   - [heading 2]
   ...

📄 [filename]
   ...
```

### Verbose Format

```
📄 [filename]
   Title: [rule file title]
   Imports: [count]

   [1] [heading]
       [rule text]
       References: [count]

   [2] [heading]
       [rule text]

   ...
```

### Empty Rules Directory

```
No rule files found in .agent-instruction/rules
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

### No Rule Files

**Output:**
```
No rule files found in .agent-instruction/rules
```

**Solution:**
Create rule files:

```bash
# Create global rules
cat > .agent-instruction/rules/global.json <<EOF
{
  "title": "Global Instructions",
  "instructions": [
    {
      "heading": "Getting Started",
      "rule": "Follow project conventions"
    }
  ]
}
EOF

# List again
agent-instruction list
```

### Malformed JSON

**Output:**
```
Warning: failed to load global.json: invalid character '}' looking for beginning of object key string
```

**Solution:**
Fix JSON syntax:

```bash
# Validate JSON
cat .agent-instruction/rules/global.json | jq .

# Edit to fix
vim .agent-instruction/rules/global.json
```

### Permission Denied

**Error:**
```
Error: list rule files: open .agent-instruction/rules: permission denied
```

**Solution:**
Check directory permissions:

```bash
# Check permissions
ls -la .agent-instruction/

# Fix if needed
chmod u+r .agent-instruction/rules/
```

## Use Cases

### Before Adding Rules

Avoid duplicates:

```bash
agent-instruction list --verbose | grep -i "pattern you want to add"
```

### After Editing Rule Files

Verify changes:

```bash
vim .agent-instruction/rules/testing.json
agent-instruction list --verbose
```

### Documentation Generation

```bash
# Create team documentation
echo "# Current AI Instructions" > docs/ai-rules.md
agent-instruction list --verbose >> docs/ai-rules.md
```

### Code Review Prep

```bash
# Show rules being applied
agent-instruction list --verbose

# Show what will be generated
agent-instruction build --dry-run
```

### Onboarding New Developers

```bash
# Show project-specific instructions
agent-instruction list --verbose | less
```

## Integration

### With Other Commands

```bash
# Check rules before building
agent-instruction list
agent-instruction build

# Add rule and verify
agent-instruction add "New rule" --rule global
agent-instruction list | grep "global"

# Review before committing
agent-instruction list --verbose
git diff .agent-instruction/
```

### In Scripts

```bash
#!/bin/bash
# Validate rules exist before build

rule_count=$(agent-instruction list | grep "instruction(s)" | wc -l)

if [ "$rule_count" -eq 0 ]; then
  echo "No rules found"
  exit 1
fi

agent-instruction build
```

### In CI/CD

```yaml
# .github/workflows/ci.yml
- name: Show current rules
  run: agent-instruction list --verbose

- name: Build instruction files
  run: agent-instruction build

- name: Verify files generated
  run: test -f CLAUDE.md
```

## Related Commands

- [`add`](add.md) - Add rules after reviewing with list
- [`build`](build.md) - Generate files from listed rules
- [`init`](init.md) - Initialize before listing

## See Also

- [Configuration Documentation](../configuration.md)
- [Rule File Format](../configuration.md#rule-files)
- [Examples](../examples.md)
- [Main README](../../README.md)

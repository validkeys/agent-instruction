# Example: Adding a Rule

## Scenario

During a code review, you discover that Claude keeps making the same mistake with error handling in Go. You want to add a rule so this pattern is remembered for future sessions.

## Conversation

**User:** I've noticed you keep returning errors without wrapping them with context. In our codebase, we always wrap errors with `fmt.Errorf` to add context about where they came from.

**Claude:** You're right, that's an important pattern. Let me add a rule to remember this for Go development.

```bash
agent-instruction add "Always wrap errors with fmt.Errorf to add context. Example: return fmt.Errorf(\"failed to read config: %w\", err)" \
  --title="Error Handling" \
  --rule="golang"
```

**Expected Output:**

```
Rule added to golang.json
Title: Error Handling
Content: Always wrap errors with fmt.Errorf to add context. Example: return fmt.Errorf("failed to read config: %w", err)

To apply changes, run: agent-instruction build
```

## What Happened

1. **Rule content** was added to `.agent-instruction/rules/golang.json`
2. **Title** helps organize and find the rule later
3. **Rule file** (`golang`) groups this with other Go-specific patterns

## Next Steps

```bash
# Review what was added
agent-instruction list --verbose

# Build to update CLAUDE.md
agent-instruction build

# Commit the changes
git add .agent-instruction/ CLAUDE.md
git commit -m "docs: Add error wrapping rule for Go"
```

## Interactive Selection

If you don't specify `--rule`, you'll be prompted:

```bash
agent-instruction add "Always wrap errors..." --title="Error Handling"
```

**Prompts:**

```
Available rule files:
  1. global
  2. golang
  3. testing
  4. security

Select a rule file (1-4): 2

Rule added to golang.json
```

## Best Practice

Check existing rules first to avoid duplicates:

```bash
# Search for existing error handling rules
agent-instruction list --verbose | grep -i "error"

# Add only if not duplicate
agent-instruction add "New rule..." --rule="golang"
```

<!-- BEGIN AGENT-INSTRUCTION -->
## Code Quality

Write clean, readable code following single responsibility principle. Functions should do one thing well. Keep functions under 50 lines when possible.

## Error Handling

Always handle errors explicitly. Never use underscore to ignore errors unless you have a specific reason documented in a comment. Wrap errors with context using fmt.Errorf with %w verb.

## Testing Requirements

Write table-driven tests for all exported functions. Test files must be placed alongside implementation files with _test.go suffix. Aim for 80% code coverage on new code.

## Documentation

Add godoc comments to all exported types, functions, and constants. Comments should explain why, not what. Include usage examples for complex APIs.

## Git Commits

Use conventional commit format: type(scope): description. Types: feat, fix, docs, test, refactor. Keep first line under 72 characters. Add detailed explanation in commit body when needed.

## Project Context

This is a single-package Go application. All source code is in the root directory. Follow standard Go project conventions.

<!-- END AGENT-INSTRUCTION -->
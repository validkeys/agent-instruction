# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-04-14

### Added

- `init` command — initializes `.agent-instruction/` directory structure, config, and global rule template; backs up existing `CLAUDE.md`/`AGENTS.md` files
- `add` command — appends instructions to a rule file without rebuilding
- `build` command — generates `CLAUDE.md` and/or `AGENTS.md` for all configured packages; preserves user-authored content outside managed sections; supports `--dry-run`, `--verbose`, and `--no-parallel` flags
- `list` command — displays all rule files and instruction summaries; `--verbose` shows full content
- Import system — rule files can import other rule files, enabling composition across global and package-level rules; cycle detection prevents infinite loops
- Parallel package processing — worker pool for fast builds across large monorepos
- Atomic file writes with backup support
- Claude Code skill integration (`skill/agent-instruction.md`)
- GitHub Actions CI/CD pipeline (test, lint, security scan, cross-platform release builds)
- golangci-lint configuration
- Example configurations for simple-repo and monorepo layouts (`examples/`)
- Command reference documentation (`docs/commands/`)
- E2E test suite (`test/e2e/full_workflow_test.sh`)
- 84%+ test coverage across all packages

### Known Limitations

- macOS and Linux only (no Windows support in v1.0.0)
- `init --packages auto` discovers packages by `agent-instruction.json` presence; directories without this file are skipped

[1.0.0]: https://github.com/validkeys/agent-instruction/releases/tag/v1.0.0

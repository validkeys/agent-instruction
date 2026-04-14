# Contributing to agent-instruction

Thank you for considering contributing to agent-instruction! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Getting Help](#getting-help)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Push to your fork
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- A GitHub account

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/agent-instruction.git
cd agent-instruction

# Add upstream remote
git remote add upstream https://github.com/validkeys/agent-instruction.git

# Install dependencies
go mod download

# Build the project
go build

# Run the binary
./agent-instruction --version
```

## Running Tests

We maintain high test coverage (target: 80%+). All contributions should include tests.

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
go test -v ./internal/builder
```

### Run Tests with Race Detection

```bash
go test -race ./...
```

## Code Style

### Go Standards

- Follow standard Go conventions and idioms
- Use `gofmt` to format code (enforced by CI)
- Run `golangci-lint` before committing (enforced by CI)

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Auto-fix issues where possible
golangci-lint run --fix
```

### Code Guidelines

- **Error Handling**: Always handle errors explicitly. Never ignore errors.
- **Testing**: Write table-driven tests for all functions
- **Documentation**: Add godoc comments for exported functions and types
- **Naming**: Use clear, descriptive names. Avoid abbreviations unless widely understood
- **Package Structure**: Keep packages focused and cohesive

## Commit Messages

We follow conventional commit message format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks, dependency updates
- `ci`: CI/CD pipeline changes

### Examples

```
feat(builder): add parallel processing for large monorepos

Implement worker pool pattern to process packages concurrently,
improving build time by ~3x for repos with 20+ packages.

Closes #42
```

```
fix(config): handle missing config file gracefully

Previously crashed with nil pointer when config.json was missing.
Now returns clear error message and suggests running init.

Fixes #38
```

## Pull Request Process

### Before Submitting

1. **Run tests**: Ensure all tests pass
2. **Run linter**: Fix all linting issues
3. **Update docs**: Update README or docs if adding features
4. **Add tests**: Include tests for new functionality
5. **Check coverage**: Maintain or improve test coverage

### PR Checklist

- [ ] Tests pass locally
- [ ] Linter passes (`golangci-lint run`)
- [ ] Added tests for new functionality
- [ ] Updated documentation if needed
- [ ] Commit messages follow conventional format
- [ ] Branch is up to date with main

### Submitting

1. Push your branch to your fork
2. Open a PR against `main` branch
3. Fill out the PR template completely
4. Link any related issues
5. Wait for CI checks to pass
6. Address review feedback promptly

### Review Process

- All PRs require at least one approval
- CI must pass (tests, linting, coverage)
- Maintainers may request changes
- Once approved and CI passes, maintainers will merge

### After Merge

- Delete your branch
- Pull latest changes from upstream
- Your contribution will be in the next release!

## Getting Help

- **Questions**: Open a [GitHub Discussion](https://github.com/validkeys/agent-instruction/discussions)
- **Bugs**: Open a [Bug Report](https://github.com/validkeys/agent-instruction/issues/new?template=bug_report.md)
- **Features**: Open a [Feature Request](https://github.com/validkeys/agent-instruction/issues/new?template=feature_request.md)
- **Security**: Email security@validkeys.com (do not open public issues)

## Development Tips

### Testing Locally

Test the tool against a real monorepo:

```bash
# Build
go build

# Test in a sample repo
cd /tmp
mkdir test-repo && cd test-repo
git init
/path/to/agent-instruction init
/path/to/agent-instruction build
```

### Debugging

Use Go's built-in debugging:

```bash
# Run with verbose output
go run main.go build --verbose

# Use delve debugger
dlv debug -- build
```

### Adding New Commands

1. Create command file in `cmd/`
2. Implement command logic in `internal/`
3. Add tests in `*_test.go` files
4. Update documentation
5. Add examples if applicable

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

## Recognition

Contributors will be recognized in:
- Release notes
- GitHub contributors page
- Special thanks in major releases

Thank you for contributing! 🙏

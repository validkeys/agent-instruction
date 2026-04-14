#!/usr/bin/env bash
# E2E test suite for agent-instruction CLI
# Tests the full workflow: init -> add -> build -> list
#
# Usage:
#   bash test/e2e/full_workflow_test.sh
#
# Exit codes:
#   0 - all tests passed
#   1 - one or more tests failed

set -euo pipefail

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

PASS=0
FAIL=0
BINARY=""

red()   { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n' "$*"; }

pass() {
  PASS=$((PASS + 1))
  green "  PASS: $1"
}

fail() {
  FAIL=$((FAIL + 1))
  red "  FAIL: $1"
  red "        $2"
}

# assert_exit_zero: run command; pass if exit 0
assert_success() {
  local label="$1"; shift
  if "$@" > /tmp/e2e_out 2>&1; then
    pass "$label"
  else
    fail "$label" "command failed (exit $?): $*"
    cat /tmp/e2e_out >&2
  fi
}

# assert_exit_nonzero: run command; pass if exit non-zero
assert_failure() {
  local label="$1"; shift
  if ! "$@" > /tmp/e2e_out 2>&1; then
    pass "$label"
  else
    fail "$label" "expected failure but command succeeded: $*"
  fi
}

# assert_output: run command; pass if stdout contains pattern
assert_output() {
  local label="$1"; local pattern="$2"; shift 2
  local out
  out=$("$@" 2>&1) || true
  if echo "$out" | grep -q "$pattern"; then
    pass "$label"
  else
    fail "$label" "expected output to contain '$pattern'; got: $out"
  fi
}

# assert_file_exists: pass if file exists
assert_file_exists() {
  local label="$1"; local path="$2"
  if [[ -f "$path" ]]; then
    pass "$label"
  else
    fail "$label" "file not found: $path"
  fi
}

# assert_file_contains: pass if file contains pattern
assert_file_contains() {
  local label="$1"; local path="$2"; local pattern="$3"
  if [[ -f "$path" ]] && grep -q "$pattern" "$path"; then
    pass "$label"
  else
    fail "$label" "file '$path' does not contain '$pattern'"
  fi
}

# assert_file_not_contains: pass if file does NOT contain pattern
assert_file_not_contains() {
  local label="$1"; local path="$2"; local pattern="$3"
  if [[ -f "$path" ]] && ! grep -q "$pattern" "$path"; then
    pass "$label"
  else
    fail "$label" "file '$path' unexpectedly contains '$pattern'"
  fi
}

# Accumulate temp dirs for cleanup at exit
E2E_TMPDIRS=()
cleanup_tmpdirs() {
  for d in "${E2E_TMPDIRS[@]+"${E2E_TMPDIRS[@]}"}"; do
    rm -rf "$d"
  done
}
trap cleanup_tmpdirs EXIT

# make_tmpdir: create temp dir, register for cleanup
make_tmpdir() {
  local d
  d=$(mktemp -d)
  E2E_TMPDIRS+=("$d")
  echo "$d"
}

# ---------------------------------------------------------------------------
# Build binary
# ---------------------------------------------------------------------------

build_binary() {
  local repo_root
  repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
  local bin_dir="$repo_root/build"
  mkdir -p "$bin_dir"
  go build -o "$bin_dir/agent-instruction" "$repo_root/cmd/agent-instruction" \
    || { red "Failed to build binary"; exit 1; }
  BINARY="$bin_dir/agent-instruction"
}

# ---------------------------------------------------------------------------
# Test suites
# ---------------------------------------------------------------------------

test_init_basic() {
  bold "Suite: init (basic)"
  local dir; dir=$(make_tmpdir)

  # init succeeds
  assert_success "init creates .agent-instruction dir" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive"

  # creates config.json
  assert_file_exists "init creates config.json" \
    "$dir/.agent-instruction/config.json"

  # creates global.json
  assert_file_exists "init creates rules/global.json" \
    "$dir/.agent-instruction/rules/global.json"

  # config contains frameworks
  assert_file_contains "config has claude framework" \
    "$dir/.agent-instruction/config.json" "claude"

  # running init again fails (already initialized)
  assert_failure "second init fails with helpful error" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive"
}

test_init_with_frameworks_flag() {
  bold "Suite: init --frameworks"
  local dir; dir=$(make_tmpdir)

  assert_success "init with --frameworks claude" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude"

  assert_file_contains "config only has claude" \
    "$dir/.agent-instruction/config.json" '"claude"'

  assert_file_not_contains "config does not have agents" \
    "$dir/.agent-instruction/config.json" '"agents"'
}

test_init_invalid_framework() {
  bold "Suite: init invalid framework"
  local dir; dir=$(make_tmpdir)

  assert_failure "init rejects unknown framework" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks badframework"
}

test_init_backs_up_existing_files() {
  bold "Suite: init backs up existing instruction files"
  local dir; dir=$(make_tmpdir)

  echo "# existing content" > "$dir/CLAUDE.md"

  assert_success "init with existing CLAUDE.md succeeds" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive"

  assert_file_exists "backup created for CLAUDE.md" \
    "$dir/CLAUDE.md.backup"

  assert_file_contains "backup preserves original content" \
    "$dir/CLAUDE.md.backup" "existing content"
}

test_build_basic() {
  bold "Suite: build (basic)"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages ." \
    > /dev/null
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  assert_success "build succeeds after init" \
    bash -c "cd '$dir' && '$BINARY' build"

  assert_file_exists "build generates CLAUDE.md at root" \
    "$dir/CLAUDE.md"

  assert_file_contains "CLAUDE.md has managed section header" \
    "$dir/CLAUDE.md" "BEGIN AGENT-INSTRUCTION"
}

test_build_dry_run() {
  bold "Suite: build --dry-run"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages ." \
    > /dev/null
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  assert_success "dry-run exits 0" \
    bash -c "cd '$dir' && '$BINARY' build --dry-run"

  # Dry run should NOT create actual files
  if [[ -f "$dir/CLAUDE.md" ]]; then
    fail "dry-run must not write files" "CLAUDE.md was created"
  else
    pass "dry-run does not write CLAUDE.md"
  fi
}

test_build_preserves_custom_content() {
  bold "Suite: build preserves user content"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages ." \
    > /dev/null
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"
  bash -c "cd '$dir' && '$BINARY' build" > /dev/null

  # Append user content outside managed section
  echo "" >> "$dir/CLAUDE.md"
  echo "## My Custom Section" >> "$dir/CLAUDE.md"
  echo "User-authored content here." >> "$dir/CLAUDE.md"

  # Rebuild
  bash -c "cd '$dir' && '$BINARY' build" > /dev/null

  assert_file_contains "custom content preserved after rebuild" \
    "$dir/CLAUDE.md" "User-authored content here."
}

test_build_without_init() {
  bold "Suite: build without init"
  local dir; dir=$(make_tmpdir)

  assert_failure "build fails when not initialized" \
    bash -c "cd '$dir' && '$BINARY' build"

  assert_output "build error mentions init" "init" \
    bash -c "cd '$dir' && '$BINARY' build 2>&1 || true"
}

test_build_verbose() {
  bold "Suite: build --verbose"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages ." \
    > /dev/null
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  assert_output "verbose output shows processing" "Processing\|Found\|package" \
    bash -c "cd '$dir' && '$BINARY' build --verbose"
}

test_add_basic() {
  bold "Suite: add (basic)"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive" > /dev/null

  assert_success "add rule to global file" \
    bash -c "cd '$dir' && '$BINARY' add 'Always write unit tests.' --rule global"

  assert_file_contains "rule content saved to global.json" \
    "$dir/.agent-instruction/rules/global.json" "Always write unit tests."
}

test_add_with_title() {
  bold "Suite: add --title"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive" > /dev/null

  assert_success "add rule with --title" \
    bash -c "cd '$dir' && '$BINARY' add 'Validate all inputs.' --title 'Input Validation' --rule global"

  assert_file_contains "heading saved to global.json" \
    "$dir/.agent-instruction/rules/global.json" "Input Validation"

  assert_file_contains "rule content saved to global.json" \
    "$dir/.agent-instruction/rules/global.json" "Validate all inputs."
}

test_add_without_init() {
  bold "Suite: add without init"
  local dir; dir=$(make_tmpdir)

  assert_failure "add fails when not initialized" \
    bash -c "cd '$dir' && '$BINARY' add 'some rule' --rule global"
}

test_add_empty_rule() {
  bold "Suite: add empty rule content"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive" > /dev/null

  assert_failure "add rejects empty rule content" \
    bash -c "cd '$dir' && '$BINARY' add '' --rule global"
}

test_list_basic() {
  bold "Suite: list (basic)"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive" > /dev/null

  assert_success "list succeeds after init" \
    bash -c "cd '$dir' && '$BINARY' list"

  assert_output "list shows global.json" "global.json" \
    bash -c "cd '$dir' && '$BINARY' list"
}

test_list_verbose() {
  bold "Suite: list --verbose"
  local dir; dir=$(make_tmpdir)

  bash -c "cd '$dir' && '$BINARY' init --non-interactive" > /dev/null
  bash -c "cd '$dir' && '$BINARY' add 'Use explicit errors.' --title 'Error Handling' --rule global" \
    > /dev/null

  assert_output "verbose list shows instruction heading" "Error Handling" \
    bash -c "cd '$dir' && '$BINARY' list --verbose"

  assert_output "verbose list shows rule content" "Use explicit errors." \
    bash -c "cd '$dir' && '$BINARY' list --verbose"
}

test_list_without_init() {
  bold "Suite: list without init"
  local dir; dir=$(make_tmpdir)

  assert_failure "list fails when not initialized" \
    bash -c "cd '$dir' && '$BINARY' list"
}

test_full_workflow_simple_repo() {
  bold "Suite: full workflow (simple repo)"
  local dir; dir=$(make_tmpdir)

  # Step 1: init
  assert_success "init" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages ."
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  # Step 2: add a rule
  assert_success "add rule" \
    bash -c "cd '$dir' && '$BINARY' add 'Write table-driven tests.' --title 'Testing' --rule global"

  # Step 3: build
  assert_success "build" \
    bash -c "cd '$dir' && '$BINARY' build"

  # Step 4: verify output
  assert_file_exists "CLAUDE.md generated" "$dir/CLAUDE.md"

  assert_file_contains "CLAUDE.md contains rule text" \
    "$dir/CLAUDE.md" "Write table-driven tests."

  # Step 5: list shows the rule
  assert_output "list shows Testing heading" "Testing" \
    bash -c "cd '$dir' && '$BINARY' list --verbose"

  # Step 6: rebuild is idempotent
  assert_success "second build succeeds" \
    bash -c "cd '$dir' && '$BINARY' build"

  assert_file_contains "CLAUDE.md rule still present after second build" \
    "$dir/CLAUDE.md" "Write table-driven tests."
}

test_full_workflow_monorepo() {
  bold "Suite: full workflow (monorepo)"
  local dir; dir=$(make_tmpdir)

  # Create package directories with agent-instruction.json so build discovers them
  mkdir -p "$dir/packages/api" "$dir/packages/worker"
  echo '{"title":"API","instructions":[{"rule":"placeholder"}]}' > "$dir/packages/api/agent-instruction.json"
  echo '{"title":"Worker","instructions":[{"rule":"placeholder"}]}' > "$dir/packages/worker/agent-instruction.json"

  # Init with explicit packages
  assert_success "init with monorepo packages" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude --packages packages/api,packages/worker"

  # Add a global rule
  assert_success "add global rule" \
    bash -c "cd '$dir' && '$BINARY' add 'Never ignore errors.' --title 'Error Handling' --rule global"

  # Build
  assert_success "build all packages" \
    bash -c "cd '$dir' && '$BINARY' build"

  # Both packages get CLAUDE.md
  assert_file_exists "api/CLAUDE.md generated" "$dir/packages/api/CLAUDE.md"
  assert_file_exists "worker/CLAUDE.md generated" "$dir/packages/worker/CLAUDE.md"

  # Both contain the global rule
  assert_file_contains "api/CLAUDE.md has error handling rule" \
    "$dir/packages/api/CLAUDE.md" "Never ignore errors."

  assert_file_contains "worker/CLAUDE.md has error handling rule" \
    "$dir/packages/worker/CLAUDE.md" "Never ignore errors."
}

test_full_workflow_agents_framework() {
  bold "Suite: full workflow (agents framework)"
  local dir; dir=$(make_tmpdir)

  assert_success "init with agents framework" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks agents --packages ."
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  assert_success "build generates AGENTS.md" \
    bash -c "cd '$dir' && '$BINARY' build"

  assert_file_exists "AGENTS.md generated" "$dir/AGENTS.md"
}

test_full_workflow_both_frameworks() {
  bold "Suite: full workflow (both frameworks)"
  local dir; dir=$(make_tmpdir)

  assert_success "init with both frameworks" \
    bash -c "cd '$dir' && '$BINARY' init --non-interactive --frameworks claude,agents --packages ."
  echo '{"title":"Root","instructions":[{"rule":"placeholder"}]}' > "$dir/agent-instruction.json"

  assert_success "build generates both files" \
    bash -c "cd '$dir' && '$BINARY' build"

  assert_file_exists "CLAUDE.md generated" "$dir/CLAUDE.md"
  assert_file_exists "AGENTS.md generated" "$dir/AGENTS.md"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

main() {
  bold "=== agent-instruction E2E test suite ==="
  echo ""

  build_binary

  test_init_basic
  test_init_with_frameworks_flag
  test_init_invalid_framework
  test_init_backs_up_existing_files

  test_build_basic
  test_build_dry_run
  test_build_preserves_custom_content
  test_build_without_init
  test_build_verbose

  test_add_basic
  test_add_with_title
  test_add_without_init
  test_add_empty_rule

  test_list_basic
  test_list_verbose
  test_list_without_init

  test_full_workflow_simple_repo
  test_full_workflow_monorepo
  test_full_workflow_agents_framework
  test_full_workflow_both_frameworks

  echo ""
  bold "=== Results ==="
  green "Passed: $PASS"
  if [[ $FAIL -gt 0 ]]; then
    red "Failed: $FAIL"
    exit 1
  else
    green "Failed: $FAIL"
    echo ""
    green "All tests passed."
  fi
}

main "$@"

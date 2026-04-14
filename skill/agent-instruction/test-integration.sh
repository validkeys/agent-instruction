#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Print functions
print_test() {
    echo -e "${BLUE}▶${NC} $1"
}

print_pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

print_fail() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

print_info() {
    echo "  ℹ $1"
}

# Test helpers
run_test() {
    ((TESTS_RUN++))
    print_test "$1"
}

# Create temporary test directory
TEST_DIR=$(mktemp -d)
trap "rm -rf ${TEST_DIR}" EXIT

echo "Integration Tests for agent-instruction Skill"
echo "=============================================="
echo

# Check if agent-instruction CLI is available
run_test "Check agent-instruction CLI is available"
if command -v agent-instruction &> /dev/null; then
    VERSION=$(agent-instruction --version 2>&1 || echo "unknown")
    print_pass "CLI found: ${VERSION}"
else
    print_fail "CLI not found in PATH"
    print_info "Install with: go install github.com/yourusername/agent-instruction@latest"
    exit 1
fi

# Test 1: Initialize in temp directory
run_test "Initialize agent-instruction in test directory"
cd "${TEST_DIR}"
if agent-instruction init > /dev/null 2>&1; then
    if [[ -d ".agent-instruction" ]] && [[ -f ".agent-instruction/config.yaml" ]]; then
        print_pass "Initialization successful"
    else
        print_fail "Initialization did not create expected files"
    fi
else
    print_fail "Initialization failed"
fi

# Test 2: Add a rule
run_test "Add a test rule"
OUTPUT=$(agent-instruction add "Test rule content" --title="Test Rule" --rule="testing" 2>&1)
if [[ $? -eq 0 ]]; then
    if [[ -f ".agent-instruction/rules/testing.json" ]]; then
        print_pass "Rule added successfully"
    else
        print_fail "Rule file not created"
        print_info "Output: ${OUTPUT}"
    fi
else
    print_fail "Add command failed"
    print_info "Output: ${OUTPUT}"
fi

# Test 3: List rules
run_test "List rules"
OUTPUT=$(agent-instruction list 2>&1)
if [[ $? -eq 0 ]]; then
    if echo "${OUTPUT}" | grep -q "testing.json"; then
        print_pass "List shows expected rule file"
    else
        print_fail "List output missing expected content"
        print_info "Output: ${OUTPUT}"
    fi
else
    print_fail "List command failed"
    print_info "Output: ${OUTPUT}"
fi

# Test 4: List with verbose flag
run_test "List rules with --verbose"
OUTPUT=$(agent-instruction list --verbose 2>&1)
if [[ $? -eq 0 ]]; then
    if echo "${OUTPUT}" | grep -q "Test Rule" && echo "${OUTPUT}" | grep -q "Test rule content"; then
        print_pass "Verbose list shows rule details"
    else
        print_fail "Verbose list missing expected content"
        print_info "Output: ${OUTPUT}"
    fi
else
    print_fail "List --verbose failed"
    print_info "Output: ${OUTPUT}"
fi

# Test 5: Build instruction files
run_test "Build instruction files"
OUTPUT=$(agent-instruction build 2>&1)
if [[ $? -eq 0 ]]; then
    if [[ -f "CLAUDE.md" ]] && [[ -f "AGENTS.md" ]]; then
        print_pass "Build generated instruction files"
    else
        print_fail "Build did not create expected files"
        print_info "Output: ${OUTPUT}"
    fi
else
    print_fail "Build command failed"
    print_info "Output: ${OUTPUT}"
fi

# Test 6: Verify CLAUDE.md content
run_test "Verify CLAUDE.md contains rule content"
if [[ -f "CLAUDE.md" ]]; then
    if grep -q "Test Rule" CLAUDE.md && grep -q "Test rule content" CLAUDE.md; then
        print_pass "CLAUDE.md contains expected rule"
    else
        print_fail "CLAUDE.md missing expected content"
        print_info "CLAUDE.md preview:"
        head -20 CLAUDE.md | sed 's/^/    /'
    fi
else
    print_fail "CLAUDE.md not found"
fi

# Test 7: Add multiple rules
run_test "Add multiple rules to different categories"
SUCCESS=true
agent-instruction add "Security rule 1" --title="Security Test" --rule="security" > /dev/null 2>&1 || SUCCESS=false
agent-instruction add "Global rule 1" --title="Global Test" --rule="global" > /dev/null 2>&1 || SUCCESS=false

if [[ "$SUCCESS" == true ]]; then
    if [[ -f ".agent-instruction/rules/security.json" ]] && [[ -f ".agent-instruction/rules/global.json" ]]; then
        print_pass "Multiple rules added to different categories"
    else
        print_fail "Not all rule files created"
    fi
else
    print_fail "Failed to add multiple rules"
fi

# Test 8: List shows all categories
run_test "List shows all rule categories"
OUTPUT=$(agent-instruction list 2>&1)
if echo "${OUTPUT}" | grep -q "testing.json" && \
   echo "${OUTPUT}" | grep -q "security.json" && \
   echo "${OUTPUT}" | grep -q "global.json"; then
    print_pass "List shows all three categories"
else
    print_fail "List missing some categories"
    print_info "Output: ${OUTPUT}"
fi

# Test 9: Build with multiple rules
run_test "Build updates files with all rules"
agent-instruction build > /dev/null 2>&1
if grep -q "Test Rule" CLAUDE.md && \
   grep -q "Security Test" CLAUDE.md && \
   grep -q "Global Test" CLAUDE.md; then
    print_pass "Build includes all rules"
else
    print_fail "Build missing some rules"
fi

# Test 10: Error handling - invalid rule file
run_test "Handle invalid rule file gracefully"
OUTPUT=$(agent-instruction add "Test" --rule="nonexistent-file-that-does-not-exist-yet" 2>&1)
if [[ $? -eq 0 ]]; then
    # Command should either create the file or prompt user
    print_pass "Handled new rule file appropriately"
else
    # Some implementations may error on non-existent files
    if echo "${OUTPUT}" | grep -qi "error\|failed"; then
        print_pass "Reported error for invalid input"
    else
        print_fail "Unexpected behavior with invalid rule file"
        print_info "Output: ${OUTPUT}"
    fi
fi

# Print summary
echo
echo "=============================================="
echo "Test Summary"
echo "=============================================="
echo "Tests run:    ${TESTS_RUN}"
echo -e "Tests passed: ${GREEN}${TESTS_PASSED}${NC}"
if [[ ${TESTS_FAILED} -gt 0 ]]; then
    echo -e "Tests failed: ${RED}${TESTS_FAILED}${NC}"
else
    echo -e "Tests failed: ${TESTS_FAILED}"
fi
echo

if [[ ${TESTS_FAILED} -eq 0 ]]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi

#!/bin/bash
# Automated regression test suite for GoShell

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

print_test() { echo -e "${YELLOW}[TEST]${NC} $1"; }
print_pass() { echo -e "${GREEN}[PASS]${NC} $1"; ((TESTS_PASSED++)); }
print_fail() { echo -e "${RED}[FAIL]${NC} $1"; ((TESTS_FAILED++)); }

# Build
echo "Building GoShell..."
go build -o /tmp/goshell-test app/*.go || exit 1
SHELL="/tmp/goshell-test"

# Cleanup
cleanup() {
    rm -f /tmp/test_* /tmp/goshell-test
    pkill -9 -f goshell-test 2>/dev/null || true
}
trap cleanup EXIT

# Run command in shell
run_cmd() {
    printf "%s\nexit\n" "$1" | $SHELL 2>&1 | grep -v '^\$' | grep -v '^$' | head -1
}

echo ""
echo "=========================================="
echo "  GoShell Regression Test Suite"
echo "=========================================="
echo ""

# TEST: Echo
print_test "Testing echo command"
((TESTS_RUN++))
output=$(run_cmd "echo hello world")
if [ "$output" = "hello world" ]; then
    print_pass "echo works"
else
    print_fail "echo: got '$output'"
fi

# TEST: PWD
print_test "Testing pwd command"
((TESTS_RUN++))
output=$(run_cmd "pwd")
if [[ "$output" =~ ^/ ]]; then
    print_pass "pwd works"
else
    print_fail "pwd: got '$output'"
fi

# TEST: Type builtin
print_test "Testing type for builtin"
((TESTS_RUN++))
output=$(run_cmd "type echo")
if [[ "$output" =~ "builtin" ]]; then
    print_pass "type identifies builtins"
else
    print_fail "type builtin: got '$output'"
fi

# TEST: Type external
print_test "Testing type for external"
((TESTS_RUN++))
output=$(run_cmd "type ls")
if [[ "$output" =~ ^ls\ is\ / ]]; then
    print_pass "type identifies externals"
else
    print_fail "type external: got '$output'"
fi

# TEST: Single quotes
print_test "Testing single quotes"
((TESTS_RUN++))
output=$(run_cmd "echo 'hello   world'")
if [ "$output" = "hello   world" ]; then
    print_pass "single quotes work"
else
    print_fail "single quotes: got '$output'"
fi

# TEST: Double quotes
print_test "Testing double quotes"
((TESTS_RUN++))
output=$(run_cmd 'echo "hello   world"')
if [ "$output" = "hello   world" ]; then
    print_pass "double quotes work"
else
    print_fail "double quotes: got '$output'"
fi

# TEST: Stdout redirect
print_test "Testing stdout redirection"
((TESTS_RUN++))
printf "echo test123 > /tmp/test_out.txt\nexit\n" | $SHELL >/dev/null 2>&1
if [ -f /tmp/test_out.txt ] && grep -q "test123" /tmp/test_out.txt; then
    print_pass "stdout redirect works"
else
    print_fail "stdout redirect failed"
fi

# TEST: Append
print_test "Testing stdout append"
((TESTS_RUN++))
printf "echo line1 > /tmp/test_append.txt\necho line2 >> /tmp/test_append.txt\nexit\n" | $SHELL >/dev/null 2>&1
lines=$(wc -l < /tmp/test_append.txt 2>/dev/null | tr -d ' ')
if [ "$lines" = "2" ]; then
    print_pass "append works"
else
    print_fail "append: got $lines lines"
fi

# TEST: Stderr redirect
print_test "Testing stderr redirection"
((TESTS_RUN++))
printf "ls /nonexistent_xyz 2> /tmp/test_err.txt\nexit\n" | $SHELL >/dev/null 2>&1
if [ -f /tmp/test_err.txt ] && [ -s /tmp/test_err.txt ]; then
    print_pass "stderr redirect works"
else
    print_fail "stderr redirect failed"
fi

# TEST: Simple pipeline
print_test "Testing simple pipeline"
((TESTS_RUN++))
printf "apple\nbanana\ncherry\n" > /tmp/test_pipe.txt
output=$(run_cmd "cat /tmp/test_pipe.txt | grep banana")
if [ "$output" = "banana" ]; then
    print_pass "pipeline works"
else
    print_fail "pipeline: got '$output'"
fi

# TEST: Multi-command pipeline (FIXED)
print_test "Testing multi-command pipeline"
((TESTS_RUN++))
printf "first\nsecond\nthird\n" > /tmp/test_multi.txt
output=$(run_cmd "cat /tmp/test_multi.txt | cat | grep second")
if [ "$output" = "second" ]; then
    print_pass "multi-pipeline works"
else
    print_fail "multi-pipeline: got '$output'"
fi

# TEST: External command
print_test "Testing external command"
((TESTS_RUN++))
output=$(run_cmd "ls /tmp")
if [ -n "$output" ]; then
    print_pass "external commands work"
else
    print_fail "external command failed"
fi

# TEST: Invalid command
print_test "Testing invalid command"
((TESTS_RUN++))
output=$(run_cmd "xyz_nonexistent" 2>&1)
if [[ "$output" =~ "not found" ]]; then
    print_pass "invalid command handled"
else
    print_fail "invalid: got '$output'"
fi

# TEST: CD absolute
print_test "Testing cd absolute path"
((TESTS_RUN++))
output=$(printf "cd /tmp\npwd\nexit\n" | $SHELL 2>&1 | grep -v '^\$' | tail -1)
if [[ "$output" =~ "/tmp" ]]; then
    print_pass "cd absolute works"
else
    print_fail "cd absolute: got '$output'"
fi

# TEST: CD home
print_test "Testing cd home"
((TESTS_RUN++))
output=$(printf "cd ~\npwd\nexit\n" | $SHELL 2>&1 | grep -v '^\$' | tail -1)
if [[ "$output" =~ ^/home/ ]] || [[ "$output" =~ ^/Users/ ]]; then
    print_pass "cd ~ works"
else
    print_fail "cd ~: got '$output'"
fi

# TEST: History
print_test "Testing history command"
((TESTS_RUN++))
export HISTFILE=/tmp/test_hist.txt
echo "echo test1" > $HISTFILE
echo "echo test2" >> $HISTFILE
output=$(printf "history\nexit\n" | HISTFILE=/tmp/test_hist.txt $SHELL 2>&1 | grep test1)
if [[ "$output" =~ "test1" ]]; then
    print_pass "history works"
else
    print_fail "history failed"
fi

# Summary
echo ""
echo "=========================================="
echo "  Test Results"
echo "=========================================="
echo "Total Tests:  $TESTS_RUN"
echo -e "${GREEN}Passed:       $TESTS_PASSED${NC}"
echo -e "${RED}Failed:       $TESTS_FAILED${NC}"
echo "=========================================="
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!!${NC}"
    exit 0
else
    echo -e "${RED}✗ $TESTS_FAILED test(s) failed${NC}"
    exit 1
fi
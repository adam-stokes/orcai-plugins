#!/usr/bin/env bash
# run_all.sh — iterate all tests/test-*.sh scripts, report PASS/SKIP/FAIL.
# Exit 1 if any test FAILs.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pass=0
skip=0
fail=0
failed_tests=()

for test_script in "$SCRIPT_DIR"/test-*.sh; do
    name="$(basename "$test_script" .sh)"
    bash "$test_script"
    rc=$?
    case $rc in
        0)
            echo "PASS  $name"
            ((pass++)) || true
            ;;
        77)
            echo "SKIP  $name"
            ((skip++)) || true
            ;;
        *)
            echo "FAIL  $name (exit $rc)"
            ((fail++)) || true
            failed_tests+=("$name")
            ;;
    esac
done

echo ""
echo "Results: ${pass} passed, ${skip} skipped, ${fail} failed"

if [[ ${fail} -gt 0 ]]; then
    echo "Failed tests: ${failed_tests[*]}"
    exit 1
fi

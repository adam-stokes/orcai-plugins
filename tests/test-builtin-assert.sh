#!/usr/bin/env bash
# test-builtin-assert.sh — run an inline assert pipeline using builtin.http_get + builtin.assert.
# Exits 77 (SKIP) if orcai is not available.

set -euo pipefail

# Check prerequisites.
if ! command -v orcai &>/dev/null; then
    echo "SKIP: orcai not found in PATH"
    exit 77
fi

# Write a temporary inline pipeline.
tmpfile="$(mktemp /tmp/orcai-assert-test-XXXXXX.pipeline.yaml)"
trap 'rm -f "$tmpfile"' EXIT

cat > "$tmpfile" << 'PIPELINE'
name: builtin-assert-test
version: "1.0"

steps:
  - id: fetch
    executor: builtin.http_get
    args:
      url: "https://httpbin.org/json"

  - id: check
    executor: builtin.assert
    needs: [fetch]
    args:
      condition: "contains:Sample Slide Show"
      value: "{{step.fetch.data.body}}"
PIPELINE

output="$(orcai pipeline run "$tmpfile" 2>&1)"
exit_code=$?

if [[ $exit_code -ne 0 ]]; then
    echo "FAIL: pipeline exited $exit_code"
    echo "$output"
    exit 1
fi

echo "PASS: builtin assert pipeline succeeded"

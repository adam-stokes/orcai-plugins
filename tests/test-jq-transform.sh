#!/usr/bin/env bash
# test-jq-transform.sh — run the jq-transform pipeline, assert exit 0 and non-empty output.
# Exits 77 (SKIP) if orcai or jq are not available.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Check prerequisites.
for tool in orcai jq; do
    if ! command -v "$tool" &>/dev/null; then
        echo "SKIP: $tool not found in PATH"
        exit 77
    fi
done

output="$(orcai pipeline run "$REPO_ROOT/examples/jq-transform.pipeline.yaml" 2>&1)"
if [[ -z "$output" ]]; then
    echo "FAIL: empty output from jq-transform pipeline"
    exit 1
fi

echo "Output: $output"

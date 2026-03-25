#!/usr/bin/env bash
# test-opencode-local.sh — run the opencode-local pipeline, assert exit 0 and non-empty output.
# Exits 77 (SKIP) if orcai, orcai-opencode, opencode, or ollama are not available.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Check prerequisites.
for tool in orcai orcai-opencode opencode ollama; do
    if ! command -v "$tool" &>/dev/null; then
        echo "SKIP: $tool not found in PATH"
        exit 77
    fi
done

# Check ollama is running.
if ! ollama list &>/dev/null; then
    echo "SKIP: ollama is not running"
    exit 77
fi

output="$(orcai pipeline run "$REPO_ROOT/examples/opencode-local.pipeline.yaml" 2>&1)"
if [[ -z "$output" ]]; then
    echo "FAIL: empty output from opencode-local pipeline"
    exit 1
fi

echo "Output preview: ${output:0:120}"

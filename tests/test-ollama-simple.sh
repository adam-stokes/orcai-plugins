#!/usr/bin/env bash
# test-ollama-simple.sh — run the ollama llama3.2 pipeline, assert non-empty output.
# Exits 77 (SKIP) if orcai, orcai-ollama, or ollama are not available.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Check prerequisites.
for tool in orcai orcai-ollama ollama; do
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

output="$(orcai pipeline run "$REPO_ROOT/examples/llama3.2-prompt.pipeline.yaml" 2>&1)"
if [[ -z "$output" ]]; then
    echo "FAIL: empty output from llama3.2-prompt pipeline"
    exit 1
fi

echo "Output preview: ${output:0:120}"

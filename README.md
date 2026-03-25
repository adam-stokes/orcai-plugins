# orcai-plugins

Official plugins for [ORCAI — the Agentic Bulletin Board System](https://github.com/adam-stokes/orcai).

Each plugin lives in `plugins/<name>/` with its own `go.mod`, `Makefile`, and sidecar YAML.
Install a plugin by running `make install` in its directory.

## Available Plugins

| Plugin | Description | Directory |
|--------|-------------|-----------|
| `ollama` | Run local models via the Ollama daemon | [`plugins/ollama/`](plugins/ollama/) |
| `opencode` | Agentic AI coding via `opencode run` (supports Ollama local models) | [`plugins/opencode/`](plugins/opencode/) |
| `jq` | Apply a jq filter expression to JSON input (pure YAML sidecar, no binary) | [`plugins/jq/`](plugins/jq/) |

## Plugin Structure

```
plugins/
  <name>/
    main.go        # CLI adapter binary source
    main_test.go   # Unit tests
    go.mod         # Standalone Go module
    Makefile       # build / install / test targets
    <name>.yaml    # Sidecar descriptor for orcai
```

## Adding a New Plugin

1. Create `plugins/<name>/` with a `go.mod` declaring module `github.com/adam-stokes/orcai-plugins/plugins/<name>`.
2. Implement a binary that reads prompt from stdin, reads config from `ORCAI_*` env vars, writes output to stdout, exits non-zero on error.
3. Create `<name>.yaml` sidecar descriptor (see existing plugins for the format).
4. Add a `Makefile` with `build`, `install`, and `test` targets.
5. Open a PR.

## End-to-End Tests

The `tests/` directory contains shell scripts that run real pipelines against real tools.
Each test skips gracefully (exit 77 → SKIP) when a required tool is not installed.

```
make test-e2e
```

**Prerequisites** (only needed for tests that use them — missing tools are SKIPped):

- `orcai` — the core CLI
- `orcai-ollama` — `cd plugins/ollama && make install`
- `orcai-opencode` — `cd plugins/opencode && make install`
- `ollama` — running daemon with `llama3.2` and/or `qwen3.5` pulled
- `opencode` — `brew install opencode` or `npm i -g opencode-ai`
- `jq` — system jq (`brew install jq`)

## Prerequisites

- [ORCAI](https://github.com/adam-stokes/orcai) installed and on `$PATH`.
- Go 1.22+ for building plugins from source.
- Each plugin may have additional prerequisites (see its README).

# orcai-plugins

Official plugins for [ORCAI — the Agentic Bulletin Board System](https://github.com/adam-stokes/orcai).

Each plugin lives in `plugins/<name>/` with its own `go.mod`, `Makefile`, and sidecar YAML.
Install a plugin by running `make install` in its directory.

## Available Plugins

| Plugin | Description | Directory |
|--------|-------------|-----------|
| `claude` | Claude AI assistant via the `claude` CLI | [`plugins/claude/`](plugins/claude/) |
| `codex` | OpenAI Codex coding agent via the `codex` CLI | [`plugins/codex/`](plugins/codex/) |
| `gemini` | Google Gemini AI via the `gemini` CLI | [`plugins/gemini/`](plugins/gemini/) |
| `github-copilot` | GitHub Copilot suggestions via the `gh` CLI | [`plugins/github-copilot/`](plugins/github-copilot/) |
| `ollama` | Run local models via the Ollama daemon | [`plugins/ollama/`](plugins/ollama/) |
| `opencode` | Agentic AI coding via `opencode run` (supports Ollama local models) | [`plugins/opencode/`](plugins/opencode/) |

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
3. Implement `--list-models` to print a JSON array of `{"id": "...", "label": "..."}` objects and exit 0.
4. Create `<name>.yaml` sidecar descriptor (see existing plugins for the format).
5. Add a `Makefile` with `build`, `install`, and `test` targets.
6. Open a PR.

## End-to-End Tests

The `tests/` directory contains shell scripts that run real pipelines against real tools.
Each test skips gracefully (exit 77 → SKIP) when a required tool is not installed.

```
make test-e2e
```

**Prerequisites** (only needed for tests that use them — missing tools are SKIPped):

- `orcai` — the core CLI
- `orcai-claude` — `cd plugins/claude && make install`
- `orcai-codex` — `cd plugins/codex && make install`
- `orcai-gemini` — `cd plugins/gemini && make install`
- `orcai-github-copilot` — `cd plugins/github-copilot && make install`
- `orcai-ollama` — `cd plugins/ollama && make install`
- `orcai-opencode` — `cd plugins/opencode && make install`
- `claude` — Claude Code CLI (`npm i -g @anthropic-ai/claude-code`)
- `codex` — OpenAI Codex CLI (`npm install -g @openai/codex`)
- `gemini` — Google Gemini CLI (`npm install -g @google/gemini-cli`)
- `gh` — GitHub CLI with `gh extension install github/gh-copilot`
- `ollama` — running daemon with desired models pulled
- `opencode` — `brew install opencode` or `npm i -g opencode-ai`

## Prerequisites

- [ORCAI](https://github.com/adam-stokes/orcai) installed and on `$PATH`.
- Go 1.22+ for building plugins from source.
- Each plugin may have additional prerequisites (see its README or sidecar YAML).

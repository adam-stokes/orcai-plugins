.PHONY: test-e2e

# Run all end-to-end tests.
# Prerequisites: orcai, orcai-ollama, orcai-opencode, ollama, jq, opencode
# Tests that are missing a required tool will be automatically SKIPPED.
test-e2e:
	@echo "Running e2e tests (missing tools are SKIPped automatically)..."
	bash tests/run_all.sh

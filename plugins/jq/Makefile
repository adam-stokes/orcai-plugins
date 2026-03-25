SIDECAR      := jq.yaml
WRAPPERS_DIR := $(HOME)/.config/orcai/wrappers

.PHONY: install test

install:
	@mkdir -p $(WRAPPERS_DIR)
	install -m 0644 $(SIDECAR) $(WRAPPERS_DIR)/$(SIDECAR)
	@echo "Installed $(SIDECAR) → $(WRAPPERS_DIR)/$(SIDECAR)"
	@echo "Prerequisite: jq must be on PATH (brew install jq)"

test:
	@command -v jq >/dev/null 2>&1 || { echo "SKIP: jq not on PATH"; exit 0; }
	@echo '{"name":"orcai"}' | sh -c 'ORCAI_FILTER=.name jq "$$ORCAI_FILTER"' | grep -q '"orcai"' && echo "PASS: jq filter works" || { echo "FAIL"; exit 1; }

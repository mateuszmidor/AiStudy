# Random City BAML demo — Makefile
BAML_VERSION    ?= 0.226.2
OPENAPI_CLI_VER ?= 2.41.0
BAML_PORT       ?= 2024

# shell used by recipes
SHELL := /bin/bash

.PHONY: help install-tools serve serve-prod ping generate tidy build run all clean

help: ## Show available targets
	@echo "BAML Random City demo — targets:"
	@echo "  install-tools  Install openapi-generator-cli (npm, global) [once]"
	@echo "  serve          DEV server: hot reload + regenerates baml_client/ (.baml changes)"
	@echo "  serve-prod     PROD server: compiled runtime, no codegen/file watching"
	@echo "  ping           Check the BAML server is up"
	@echo "  generate       Regenerate Go client from baml_src/ — server-free (for CI/deploys)"
	@echo "  tidy           go mod tidy"
	@echo "  build          Ensure client + go build ./..."
	@echo "  run            build + go run main.go (needs a server: serve or serve-prod)"
	@echo "  all            install-tools + run"
	@echo "  clean          Remove generated baml_client/"

install-tools: ## Install codegen tooling
	npm install -g @openapitools/openapi-generator-cli@$(OPENAPI_CLI_VER)

serve: ## Start BAML dev server (regenerates baml_client/ on file changes)
	@test -n "$$OPENAI_API_KEY" || { echo "OPENAI_API_KEY is not set"; exit 1; }
	npx @boundaryml/baml@$(BAML_VERSION) dev

serve-prod: ## Start the production BAML server (compiled runtime, no codegen)
	@test -n "$$OPENAI_API_KEY" || { echo "OPENAI_API_KEY is not set"; exit 1; }
	npx @boundaryml/baml@$(BAML_VERSION) serve --port $(BAML_PORT) --dotenv

ping: ## Verify the BAML server responds
	curl -s http://localhost:$(BAML_PORT)/_debug/ping; echo

# baml_client/ is produced from baml_src/ by `baml generate` (server-free).
# Make re-runs this recipe only when the directory is missing.
baml_client:
	@if [ ! -d baml_client ]; then \
		echo "baml_client/ missing — generating via baml generate..."; \
		npx @boundaryml/baml@$(BAML_VERSION) generate; \
	fi

generate: ## Regenerate the Go client from baml_src/ (no server needed; use in CI)
	npx @boundaryml/baml@$(BAML_VERSION) generate

tidy: ## Resolve Go module deps
	go mod tidy

build: baml_client ## Compile the Go app
	go build ./...

run: build ## Run the demo (start `make serve` or `make serve-prod` in another terminal first)
	go run main.go

all: install-tools ## Install tools, then run the demo
	$(MAKE) run

clean: ## Remove the generated client/schema
	rm -rf baml_client
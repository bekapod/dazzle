.PHONY: help all build test lint lint-fix clean run dev fmt vet

help: ## Show available commands
	@echo "Available commands:"
	@echo
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo

all: build test lint ## Build, test, and lint

build: ## Build binary to bin/dazzle
	go build -o bin/dazzle .

test: ## Run all tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

lint: vet ## Run golangci-lint
	golangci-lint run

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	go fmt ./...

run: build ## Run built binary (SPEC=<path-or-url>)
	@test -n "$(SPEC)" || (echo "Usage: make run SPEC=<spec-file-or-url>" && exit 1)
	./bin/dazzle $(SPEC)

dev: ## Run without building (SPEC=<path-or-url>)
	@test -n "$(SPEC)" || (echo "Usage: make dev SPEC=<spec-file-or-url>" && exit 1)
	go run . $(SPEC)

clean: ## Clean build artifacts
	rm -rf bin/ coverage.txt

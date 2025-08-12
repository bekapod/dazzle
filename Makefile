.PHONY: help all build test lint clean run dev mod-tidy vet

# Default target
help: ## Show available commands
	@echo "Available make commands:"
	@echo
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-10s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo

all: build ## Build, test, and lint
	$(MAKE) test
	$(MAKE) lint

build: ## Build
	@echo "Building dazzle application..."
	go build -o bin/dazzle .

test: build ## Run tests
	@echo "Running tests..."
	go test -v ./...

lint: vet ## Run linting checks
	@echo "Running linting checks..."
	golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy and verify go modules
	@echo "Tidying go modules..."
	go mod tidy
	go mod verify

run: build ## Run the built app
	@echo "Running dazzle application..."
	@echo "Usage: make run URL=<url>"
	@test -n "$(URL)" || (echo "Error: URL argument required. Use: make run URL=<url>" && exit 1)
	./bin/dazzle $(URL)

dev: ## Run app (no build)
	@echo "Running dazzle in development mode..."
	@echo "Usage: make dev URL=<url>"
	@test -n "$(URL)" || (echo "Error: URL argument required. Use: make dev URL=<url>" && exit 1)
	go run . $(URL)

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

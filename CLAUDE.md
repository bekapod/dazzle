# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Dazzle is a Go CLI application that parses OpenAPI specifications and provides an interactive terminal UI for browsing API operations. It uses the Charmbracelet TUI ecosystem (bubbletea, bubbles, lipgloss) and kin-openapi for spec parsing.

## Common Commands

```bash
make build          # Build binary to bin/dazzle
make test           # Run all tests
make lint           # Run golangci-lint
make dev URL=<url>  # Run without building (go run . <url>)
make run URL=<url>  # Run built binary
make all            # Build + test + lint
```

Run a single test:
```bash
go test -v -run TestOutputAfterFiltering ./...
```

## Architecture

The codebase follows Clean Architecture with four layers:

```
main.go                          # Entry point - wires up layers
├── internal/infrastructure/     # Data access (OpenAPI parsing)
│   └── openapi_repository.go    # Implements domain.OperationRepository
├── internal/application/        # Business logic
│   └── operation_service.go     # Implements domain.OperationService
├── internal/domain/             # Core models & interfaces
│   └── operation.go             # Operation, HTTPMethod, interfaces
└── internal/ui/                 # Presentation (TUI)
    ├── app.go                   # Main Bubbletea model
    ├── operation_list.go        # List component with filtering
    └── styles.go                # Catppuccin-based theming
```

**Dependency flow**: UI → Application → Domain ← Infrastructure

Each layer depends only on the domain interfaces, not on concrete implementations. This enables testing with mocks.

## Key Patterns

**Bubbletea Model-Update-View**: The UI uses Elm-inspired architecture where `Update()` handles messages and returns the new model, and `View()` renders the current state. All UI components implement `tea.Model`.

**Repository Pattern**: `OperationRepository` interface abstracts OpenAPI spec loading. The infrastructure layer implements it, supporting both file paths and URLs.

**Mock Testing**: Service tests use inline mock structs implementing domain interfaces (see `operation_service_test.go`).

**Golden File Testing**: TUI tests use teatest with golden files in `testdata/`. Set `lipgloss.SetColorProfile(termenv.Ascii)` in test init to ensure deterministic output.

## Testing

- Unit tests for domain/service logic use standard Go testing with mock repositories
- TUI integration tests use `teatest` with `tea.KeyMsg` for simulating user input
- Golden files in `testdata/` capture expected UI output
- Fixtures in `fixtures/` contain test OpenAPI specs

# CLAUDE.md

## Project Overview

Dazzle is a spec-aware, terminal-native API explorer and test runner. It parses OpenAPI specifications and provides an interactive TUI for browsing operations, making requests, and validating responses.

## Commands

```bash
make build              # Build binary to bin/dazzle
make test               # Run all tests with race detector
make lint               # Run golangci-lint
make fmt                # Format code
make dev SPEC=<path>    # Run without building
make run SPEC=<path>    # Run built binary
make all                # Build + test + lint
```

Debug logging:

```bash
DAZZLE_DEBUG=debug.log ./bin/dazzle spec.yaml
```

Run a single test:

```bash
go test -v -run TestOperationService_SortOperations ./internal/application
```

## Architecture

Clean Architecture with dependency inversion:

```
main.go                                  # Entry point, DI wiring
├── internal/domain/                     # Core entities & interfaces
├── internal/application/                # Business logic (services)
├── internal/infrastructure/openapi/     # kin-openapi integration
└── internal/ui/                         # Bubbletea TUI
    ├── app.go                           # Root model, screen routing
    ├── screen.go                        # Screen interface
    ├── screens/                         # Individual screens
    └── styles/                          # Catppuccin palette & styles
```

**Dependency flow:** UI → Application → Domain ← Infrastructure

### Key Patterns

- **Constructor-based DI**: main.go wires repo → services → UI
- **Bubbletea Model-Update-View**: Elm-inspired architecture for the TUI
- **Screen interface**: Each screen implements `tea.Model` + `Name()`. AppModel routes between screens.
- **Adapter pattern**: `infrastructure/openapi/adapter.go` maps kin-openapi types → domain types
- **Repository pattern**: `SpecRepository` interface abstracts spec loading

## Testing

- Unit tests with mock repositories (implement domain interfaces inline)
- Test fixtures in `testdata/fixtures/`
- All tests run with `-race` flag

## Tech Stack

- Go, Bubbletea/Bubbles/Lipgloss, kin-openapi

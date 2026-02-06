# dazzle

A spec-aware, terminal-native API explorer.

Dazzle parses your OpenAPI spec and provides an interactive terminal UI for browsing endpoints, filtering, and keyboard navigation.

## Usage

```bash
# From a local file
dazzle ./openapi.yaml

# From a URL
dazzle https://petstore3.swagger.io/api/v3/openapi.json
```

## Install

```bash
git clone https://github.com/yourusername/dazzle.git
cd dazzle
make build
```

Requires Go 1.23+.

## Development

```bash
make build    # Build binary
make test     # Run tests
make lint     # Lint
make all      # All of the above
```

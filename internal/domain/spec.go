package domain

// Spec represents a parsed OpenAPI specification.
type Spec struct {
	Info       SpecInfo
	Servers    []Server
	Operations []Operation
}

// SpecInfo contains metadata about the API.
type SpecInfo struct {
	Title       string
	Description string
	Version     string
}

// Server represents an API server endpoint.
type Server struct {
	URL         string
	Description string
}

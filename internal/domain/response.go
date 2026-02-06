package domain

// Response represents an API response for a status code.
type Response struct {
	Description string
	Content     map[string]MediaType
	Headers     map[string]Header
}

// Header represents a response header.
type Header struct {
	Description string
	Schema      *Schema
}

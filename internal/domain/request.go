package domain

// RequestBody represents an operation's request body.
type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]MediaType
}

// MediaType represents a media type with schema and example.
type MediaType struct {
	Schema *Schema
}

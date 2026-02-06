package domain

// Schema represents a JSON Schema definition.
type Schema struct {
	Type        SchemaType
	Format      string
	Description string
	Required    []string
	Properties  map[string]*Schema
	Items       *Schema
	Enum        []any
}

// SchemaType represents the data type of a schema.
type SchemaType string

const (
	SchemaTypeString  SchemaType = "string"
	SchemaTypeNumber  SchemaType = "number"
	SchemaTypeInteger SchemaType = "integer"
	SchemaTypeBoolean SchemaType = "boolean"
	SchemaTypeArray   SchemaType = "array"
	SchemaTypeObject  SchemaType = "object"
)

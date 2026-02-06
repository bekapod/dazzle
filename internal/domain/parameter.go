package domain

// Parameter represents an operation parameter.
type Parameter struct {
	Name        string
	In          ParameterIn
	Description string
	Required    bool
	Schema      *Schema
}

// ParameterIn indicates where the parameter appears.
type ParameterIn string

const (
	ParameterInPath   ParameterIn = "path"
	ParameterInQuery  ParameterIn = "query"
	ParameterInHeader ParameterIn = "header"
	ParameterInCookie ParameterIn = "cookie"
)

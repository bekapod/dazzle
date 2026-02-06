package domain

// Operation represents a single API operation (path + method).
type Operation struct {
	ID          string
	Path        string
	Method      HTTPMethod
	Summary     string
	Description string
	Tags        []string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[int]Response
}

// HTTPMethod represents an HTTP request method.
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

// OperationFilter defines criteria for filtering operations.
type OperationFilter struct {
	Query  string
	Tags   []string
	Method HTTPMethod
}

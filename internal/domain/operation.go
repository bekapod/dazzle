package domain

// Operation represents a single API operation
type Operation struct {
	Path    string
	Method  HTTPMethod
	Summary string
	Tags    []string
}

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
)

// OperationRepository defines the interface for operation data access
type OperationRepository interface {
	GetOperations() ([]Operation, error)
}

// OperationService defines business logic for operations
type OperationService interface {
	ListOperations() ([]Operation, error)
	SortOperations(operations []Operation) []Operation
}
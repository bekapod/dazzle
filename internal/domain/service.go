package domain

import "context"

// SpecService provides business logic for spec operations.
type SpecService interface {
	LoadSpec(ctx context.Context, source string) (*Spec, error)
	GetInfo(spec *Spec) SpecInfo
}

// OperationService provides business logic for operations.
type OperationService interface {
	ListOperations(spec *Spec) []Operation
	FilterOperations(operations []Operation, filter OperationFilter) []Operation
	SortOperations(operations []Operation) []Operation
}

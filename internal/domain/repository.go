package domain

import "context"

// SpecRepository loads OpenAPI specifications from a source.
type SpecRepository interface {
	Load(ctx context.Context, source string) (*Spec, error)
}

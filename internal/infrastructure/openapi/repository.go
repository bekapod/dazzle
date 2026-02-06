package openapi

import (
	"context"
	"fmt"
	"net/url"

	"dazzle/internal/domain"

	oas "github.com/getkin/kin-openapi/openapi3"
)

// Repository loads OpenAPI specs from files or URLs.
type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Load(ctx context.Context, source string) (*domain.Spec, error) {
	loader := oas.NewLoader()
	loader.Context = ctx
	loader.IsExternalRefsAllowed = true

	doc, err := r.loadDoc(loader, source)
	if err != nil {
		return nil, fmt.Errorf("loading spec from %s: %w", source, err)
	}

	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("validating spec: %w", err)
	}

	return adaptSpec(doc), nil
}

func (r *Repository) loadDoc(loader *oas.Loader, source string) (*oas.T, error) {
	if u, err := url.ParseRequestURI(source); err == nil && u.Scheme != "" {
		return loader.LoadFromURI(u)
	}
	return loader.LoadFromFile(source)
}

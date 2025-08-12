package infrastructure

import (
	"context"
	"net/url"

	"dazzle/internal/domain"
	oas "github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIRepository struct {
	doc *oas.T
}

func NewOpenAPIRepository(docURL string) (*OpenAPIRepository, error) {
	parsedURL, err := url.ParseRequestURI(docURL)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	loader := &oas.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromURI(parsedURL)
	if err != nil {
		return nil, err
	}

	return &OpenAPIRepository{doc: doc}, nil
}

func NewOpenAPIRepositoryFromFile(filePath string) (*OpenAPIRepository, error) {
	loader := oas.NewLoader()
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return &OpenAPIRepository{doc: doc}, nil
}

func (r *OpenAPIRepository) GetOperations() ([]domain.Operation, error) {
	var operations []domain.Operation

	for path, pathItem := range r.doc.Paths.Map() {
		operations = append(operations, r.extractOperations(path, pathItem)...)
	}

	return operations, nil
}


func (r *OpenAPIRepository) extractOperations(path string, pathItem *oas.PathItem) []domain.Operation {
	var operations []domain.Operation

	if pathItem.Get != nil {
		operations = append(operations, r.createOperation(path, domain.GET, pathItem.Get))
	}
	if pathItem.Post != nil {
		operations = append(operations, r.createOperation(path, domain.POST, pathItem.Post))
	}
	if pathItem.Put != nil {
		operations = append(operations, r.createOperation(path, domain.PUT, pathItem.Put))
	}
	if pathItem.Patch != nil {
		operations = append(operations, r.createOperation(path, domain.PATCH, pathItem.Patch))
	}
	if pathItem.Delete != nil {
		operations = append(operations, r.createOperation(path, domain.DELETE, pathItem.Delete))
	}

	return operations
}

func (r *OpenAPIRepository) createOperation(path string, method domain.HTTPMethod, op *oas.Operation) domain.Operation {
	return domain.Operation{
		Path:    path,
		Method:  method,
		Summary: op.Summary,
		Tags:    op.Tags,
	}
}
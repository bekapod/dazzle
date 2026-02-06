package openapi

import (
	"dazzle/internal/domain"

	oas "github.com/getkin/kin-openapi/openapi3"
)

func adaptSpec(doc *oas.T) *domain.Spec {
	spec := &domain.Spec{
		Info: domain.SpecInfo{
			Title:       doc.Info.Title,
			Description: doc.Info.Description,
			Version:     doc.Info.Version,
		},
		Servers: adaptServers(doc.Servers),
	}

	if doc.Paths != nil {
		for path, item := range doc.Paths.Map() {
			spec.Operations = append(spec.Operations, extractOperations(path, item)...)
		}
	}

	return spec
}

func adaptServers(servers oas.Servers) []domain.Server {
	result := make([]domain.Server, len(servers))
	for i, s := range servers {
		result[i] = domain.Server{
			URL:         s.URL,
			Description: s.Description,
		}
	}
	return result
}

func extractOperations(path string, item *oas.PathItem) []domain.Operation {
	type entry struct {
		method domain.HTTPMethod
		op     *oas.Operation
	}

	candidates := []entry{
		{domain.GET, item.Get},
		{domain.POST, item.Post},
		{domain.PUT, item.Put},
		{domain.PATCH, item.Patch},
		{domain.DELETE, item.Delete},
		{domain.HEAD, item.Head},
		{domain.OPTIONS, item.Options},
	}

	var ops []domain.Operation
	for _, c := range candidates {
		if c.op != nil {
			ops = append(ops, adaptOperation(path, c.method, c.op))
		}
	}
	return ops
}

func adaptOperation(path string, method domain.HTTPMethod, op *oas.Operation) domain.Operation {
	id := op.OperationID
	if id == "" {
		id = string(method) + " " + path
	}

	return domain.Operation{
		ID:          id,
		Path:        path,
		Method:      method,
		Summary:     op.Summary,
		Description: op.Description,
		Tags:        op.Tags,
		Parameters:  adaptParameters(op.Parameters),
	}
}

func adaptParameters(params oas.Parameters) []domain.Parameter {
	if len(params) == 0 {
		return nil
	}

	result := make([]domain.Parameter, 0, len(params))
	for _, p := range params {
		if p.Value == nil {
			continue
		}
		result = append(result, domain.Parameter{
			Name:        p.Value.Name,
			In:          domain.ParameterIn(p.Value.In),
			Description: p.Value.Description,
			Required:    p.Value.Required,
		})
	}
	return result
}

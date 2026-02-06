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
			ops = append(ops, adaptOperation(path, c.method, item.Parameters, c.op))
		}
	}
	return ops
}

func adaptOperation(path string, method domain.HTTPMethod, pathParams oas.Parameters, op *oas.Operation) domain.Operation {
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
		Parameters:  adaptParameters(mergeParameters(pathParams, op.Parameters)),
		RequestBody: adaptRequestBody(op.RequestBody),
		Responses:   adaptResponses(op.Responses),
	}
}

// mergeParameters combines path-level and operation-level parameters.
// Operation-level parameters override path-level parameters with the same name and location.
// The key uses In+Name directly — In is a lowercase enum per the OpenAPI spec,
// and parameter names are case-sensitive, so no normalisation is needed.
func mergeParameters(pathParams, opParams oas.Parameters) oas.Parameters {
	if len(pathParams) == 0 {
		return opParams
	}
	if len(opParams) == 0 {
		return pathParams
	}

	// Index operation params by name+in for override lookup.
	opSet := make(map[string]struct{}, len(opParams))
	for _, p := range opParams {
		if p.Value != nil {
			opSet[p.Value.In+":"+p.Value.Name] = struct{}{}
		}
	}

	merged := make(oas.Parameters, 0, len(pathParams)+len(opParams))
	for _, p := range pathParams {
		if p.Value == nil {
			continue
		}
		if _, overridden := opSet[p.Value.In+":"+p.Value.Name]; !overridden {
			merged = append(merged, p)
		}
	}
	merged = append(merged, opParams...)
	return merged
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
			Schema:      adaptSchemaRef(p.Value.Schema),
		})
	}
	return result
}

func adaptRequestBody(ref *oas.RequestBodyRef) *domain.RequestBody {
	if ref == nil || ref.Value == nil {
		return nil
	}
	rb := ref.Value
	return &domain.RequestBody{
		Description: rb.Description,
		Required:    rb.Required,
		Content:     adaptContent(rb.Content),
	}
}

func adaptResponses(responses *oas.Responses) map[string]domain.Response {
	if responses == nil {
		return nil
	}
	m := responses.Map()
	if len(m) == 0 {
		return nil
	}

	result := make(map[string]domain.Response, len(m))
	for code, ref := range m {
		if ref.Value == nil {
			continue
		}
		r := ref.Value
		resp := domain.Response{
			Content: adaptContent(r.Content),
			Headers: adaptHeaders(r.Headers),
		}
		if r.Description != nil {
			resp.Description = *r.Description
		}
		result[code] = resp
	}
	return result
}

func adaptContent(content oas.Content) map[string]domain.MediaType {
	if len(content) == 0 {
		return nil
	}
	result := make(map[string]domain.MediaType, len(content))
	for mediaType, mt := range content {
		result[mediaType] = domain.MediaType{
			Schema: adaptSchemaRef(mt.Schema),
		}
	}
	return result
}

// schemaMaxDepth controls how many levels of Properties/Items are expanded.
const schemaMaxDepth = 3

func adaptSchemaRef(ref *oas.SchemaRef) *domain.Schema {
	if ref == nil || ref.Value == nil {
		return nil
	}
	return adaptSchema(ref.Value, schemaMaxDepth)
}

// adaptSchema maps an OAS schema to a domain schema, recursing into Properties
// and Items up to the given depth. At depth 0 only scalar fields are copied.
func adaptSchema(s *oas.Schema, depth int) *domain.Schema {
	ds := &domain.Schema{
		Format:      s.Format,
		Description: s.Description,
		Required:    s.Required,
	}

	if s.Type != nil && len(*s.Type) > 0 {
		ds.Type = domain.SchemaType((*s.Type)[0])
	}

	if len(s.Enum) > 0 {
		ds.Enum = s.Enum
	}

	if depth <= 0 {
		return ds
	}

	if s.Items != nil && s.Items.Value != nil {
		ds.Items = adaptSchema(s.Items.Value, depth-1)
	}

	if len(s.Properties) > 0 {
		ds.Properties = make(map[string]*domain.Schema, len(s.Properties))
		for name, propRef := range s.Properties {
			if propRef.Value != nil {
				ds.Properties[name] = adaptSchema(propRef.Value, depth-1)
			}
		}
	}

	return ds
}

func adaptHeaders(headers oas.Headers) map[string]domain.Header {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]domain.Header, len(headers))
	for name, ref := range headers {
		if ref.Value == nil {
			continue
		}
		h := domain.Header{
			Description: ref.Value.Description,
		}
		if ref.Value.Schema != nil {
			h.Schema = adaptSchemaRef(ref.Value.Schema)
		}
		result[name] = h
	}
	return result
}

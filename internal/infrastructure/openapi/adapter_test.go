package openapi

import (
	"testing"

	"dazzle/internal/domain"

	oas "github.com/getkin/kin-openapi/openapi3"
)

func makeParamRef(name, in string) *oas.ParameterRef {
	return &oas.ParameterRef{
		Value: &oas.Parameter{Name: name, In: in},
	}
}

func TestMergeParameters_OperationOverridesPath(t *testing.T) {
	pathParams := oas.Parameters{
		makeParamRef("petId", "path"),
	}
	opParams := oas.Parameters{
		makeParamRef("petId", "path"),
	}

	merged := mergeParameters(pathParams, opParams)

	if len(merged) != 1 {
		t.Fatalf("expected 1 merged param (op wins), got %d", len(merged))
	}
	// The surviving param should be the operation-level one.
	if merged[0] != opParams[0] {
		t.Error("expected operation-level param to win the override")
	}
}

func TestMergeParameters_DifferentLocationsNotOverridden(t *testing.T) {
	pathParams := oas.Parameters{
		makeParamRef("id", "path"),
	}
	opParams := oas.Parameters{
		makeParamRef("id", "query"),
	}

	merged := mergeParameters(pathParams, opParams)

	if len(merged) != 2 {
		t.Fatalf("expected 2 params (different locations), got %d", len(merged))
	}
}

func TestMergeParameters_NilValuePathParamsSkippedDuringMerge(t *testing.T) {
	// When both path and op params exist, nil Value path params are
	// skipped in the merge loop (they can't be indexed for override).
	pathParams := oas.Parameters{
		{Value: nil},
		makeParamRef("petId", "path"),
	}
	opParams := oas.Parameters{
		makeParamRef("limit", "query"),
	}

	merged := mergeParameters(pathParams, opParams)

	if len(merged) != 2 {
		t.Fatalf("expected 2 params (nil skipped in merge), got %d", len(merged))
	}
}

func TestAdaptParameters_NilValueSkipped(t *testing.T) {
	// adaptParameters filters out any params with nil Value,
	// regardless of how they arrived from mergeParameters.
	params := oas.Parameters{
		{Value: nil},
		makeParamRef("petId", "path"),
	}

	adapted := adaptParameters(params)

	if len(adapted) != 1 {
		t.Fatalf("expected 1 adapted param (nil skipped), got %d", len(adapted))
	}
	if adapted[0].Name != "petId" {
		t.Errorf("expected petId, got %q", adapted[0].Name)
	}
}

func TestMergeParameters_NilValueOpParamsInSet(t *testing.T) {
	pathParams := oas.Parameters{
		makeParamRef("id", "path"),
	}
	opParams := oas.Parameters{
		{Value: nil}, // nil Value op param — shouldn't cause override
	}

	merged := mergeParameters(pathParams, opParams)

	// path param "id" should survive since the nil op param can't override it
	found := false
	for _, p := range merged {
		if p.Value != nil && p.Value.Name == "id" {
			found = true
		}
	}
	if !found {
		t.Error("expected path-level 'id' param to survive when op param has nil Value")
	}
}

func TestAdaptSchema_ObjectWithProperties(t *testing.T) {
	strType := oas.Types{"string"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type:     &objType,
		Required: []string{"name"},
		Properties: oas.Schemas{
			"name": &oas.SchemaRef{
				Value: &oas.Schema{Type: &strType, Description: "Pet name"},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Type != domain.SchemaTypeObject {
		t.Errorf("expected object type, got %q", ds.Type)
	}
	if len(ds.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(ds.Properties))
	}
	nameProp := ds.Properties["name"]
	if nameProp.Type != domain.SchemaTypeString {
		t.Errorf("expected string type for name, got %q", nameProp.Type)
	}
}

func TestAdaptSchema_DepthZeroDropsPropertiesAndItems(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type: &objType,
		Properties: oas.Schemas{
			"name": &oas.SchemaRef{
				Value: &oas.Schema{Type: &strType},
			},
		},
		Items: &oas.SchemaRef{
			Value: &oas.Schema{Type: &arrType},
		},
	}

	ds := adaptSchema(schema, 0)

	if ds.Type != domain.SchemaTypeObject {
		t.Errorf("expected object type, got %q", ds.Type)
	}
	if ds.Properties != nil {
		t.Error("expected nil properties at depth 0")
	}
	if ds.Items != nil {
		t.Error("expected nil items at depth 0")
	}
}

func TestAdaptSchema_ArrayWithItems(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	schema := &oas.Schema{
		Type: &arrType,
		Items: &oas.SchemaRef{
			Value: &oas.Schema{Type: &strType},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Type != domain.SchemaTypeArray {
		t.Errorf("expected array type, got %q", ds.Type)
	}
	if ds.Items == nil {
		t.Fatal("expected items schema")
	}
	if ds.Items.Type != domain.SchemaTypeString {
		t.Errorf("expected string item type, got %q", ds.Items.Type)
	}
}

func TestAdaptSchema_ArrayPropertyRetainsItems(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type: &objType,
		Properties: oas.Schemas{
			"tags": &oas.SchemaRef{
				Value: &oas.Schema{
					Type: &arrType,
					Items: &oas.SchemaRef{
						Value: &oas.Schema{Type: &strType},
					},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	tagsProp := ds.Properties["tags"]
	if tagsProp == nil {
		t.Fatal("expected tags property")
	}
	if tagsProp.Type != domain.SchemaTypeArray {
		t.Errorf("expected array type, got %q", tagsProp.Type)
	}
	if tagsProp.Items == nil {
		t.Fatal("expected items schema for array property")
	}
	if tagsProp.Items.Type != domain.SchemaTypeString {
		t.Errorf("expected string item type, got %q", tagsProp.Items.Type)
	}
}

func TestAdaptSchema_ArrayPropertyWithObjectItems(t *testing.T) {
	strType := oas.Types{"string"}
	intType := oas.Types{"integer"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type: &objType,
		Properties: oas.Schemas{
			"pets": &oas.SchemaRef{
				Value: &oas.Schema{
					Type: &arrType,
					Items: &oas.SchemaRef{
						Value: &oas.Schema{
							Type: &objType,
							Properties: oas.Schemas{
								"id":   &oas.SchemaRef{Value: &oas.Schema{Type: &intType}},
								"name": &oas.SchemaRef{Value: &oas.Schema{Type: &strType}},
							},
						},
					},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	petsProp := ds.Properties["pets"]
	if petsProp == nil {
		t.Fatal("expected pets property")
	}
	if petsProp.Items == nil {
		t.Fatal("expected items for pets array property")
	}
	if len(petsProp.Items.Properties) != 2 {
		t.Fatalf("expected 2 item properties, got %d", len(petsProp.Items.Properties))
	}
	if petsProp.Items.Properties["id"].Type != domain.SchemaTypeInteger {
		t.Errorf("expected integer for id, got %q", petsProp.Items.Properties["id"].Type)
	}
	if petsProp.Items.Properties["name"].Type != domain.SchemaTypeString {
		t.Errorf("expected string for name, got %q", petsProp.Items.Properties["name"].Type)
	}
}

func TestAdaptSchema_ArrayItemObjectWithArrayProperty(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type: &arrType,
		Items: &oas.SchemaRef{
			Value: &oas.Schema{
				Type: &objType,
				Properties: oas.Schemas{
					"tags": &oas.SchemaRef{
						Value: &oas.Schema{
							Type: &arrType,
							Items: &oas.SchemaRef{
								Value: &oas.Schema{Type: &strType},
							},
						},
					},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Items == nil {
		t.Fatal("expected items schema")
	}
	tagsProp := ds.Items.Properties["tags"]
	if tagsProp == nil {
		t.Fatal("expected tags property in item schema")
	}
	if tagsProp.Type != domain.SchemaTypeArray {
		t.Errorf("expected array type for tags, got %q", tagsProp.Type)
	}
	if tagsProp.Items == nil {
		t.Fatal("expected items for tags array property inside object items")
	}
	if tagsProp.Items.Type != domain.SchemaTypeString {
		t.Errorf("expected string item type for tags, got %q", tagsProp.Items.Type)
	}
}

func TestAdaptSchema_PropertyNestedArrayOfObjectsAtDepthBoundary(t *testing.T) {
	intType := oas.Types{"integer"}
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	// object → array → array → object{properties} = 4 levels of nesting.
	// At schemaMaxDepth=3, the innermost object's type is preserved but
	// its properties are beyond the depth limit.
	schema := &oas.Schema{
		Type: &objType,
		Properties: oas.Schemas{
			"matrix": &oas.SchemaRef{
				Value: &oas.Schema{
					Type: &arrType,
					Items: &oas.SchemaRef{
						Value: &oas.Schema{
							Type: &arrType,
							Items: &oas.SchemaRef{
								Value: &oas.Schema{
									Type: &objType,
									Properties: oas.Schemas{
										"id":   &oas.SchemaRef{Value: &oas.Schema{Type: &intType}},
										"name": &oas.SchemaRef{Value: &oas.Schema{Type: &strType}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	matrixProp := ds.Properties["matrix"]
	if matrixProp == nil || matrixProp.Items == nil || matrixProp.Items.Items == nil {
		t.Fatal("expected matrix → items → items chain")
	}
	innerObj := matrixProp.Items.Items
	if innerObj.Type != domain.SchemaTypeObject {
		t.Errorf("expected object type, got %q", innerObj.Type)
	}
	// At depth 0 the innermost object's properties are not expanded.
	if innerObj.Properties != nil {
		t.Error("expected nil properties at depth boundary")
	}
}

func TestAdaptSchema_PropertyNestedArrayRetainsInnerItems(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	objType := oas.Types{"object"}
	schema := &oas.Schema{
		Type: &objType,
		Properties: oas.Schemas{
			"matrix": &oas.SchemaRef{
				Value: &oas.Schema{
					Type: &arrType,
					Items: &oas.SchemaRef{
						Value: &oas.Schema{
							Type: &arrType,
							Items: &oas.SchemaRef{
								Value: &oas.Schema{Type: &strType},
							},
						},
					},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	matrixProp := ds.Properties["matrix"]
	if matrixProp == nil {
		t.Fatal("expected matrix property")
	}
	if matrixProp.Items == nil {
		t.Fatal("expected items for matrix array property")
	}
	if matrixProp.Items.Type != domain.SchemaTypeArray {
		t.Errorf("expected array inner type, got %q", matrixProp.Items.Type)
	}
	if matrixProp.Items.Items == nil {
		t.Fatal("expected inner items for array[array[string]] property")
	}
	if matrixProp.Items.Items.Type != domain.SchemaTypeString {
		t.Errorf("expected string innermost type, got %q", matrixProp.Items.Items.Type)
	}
}

func TestAdaptSchema_NestedArrayItems(t *testing.T) {
	strType := oas.Types{"string"}
	arrType := oas.Types{"array"}
	schema := &oas.Schema{
		Type: &arrType,
		Items: &oas.SchemaRef{
			Value: &oas.Schema{
				Type: &arrType,
				Items: &oas.SchemaRef{
					Value: &oas.Schema{Type: &strType},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Items == nil {
		t.Fatal("expected items schema")
	}
	if ds.Items.Type != domain.SchemaTypeArray {
		t.Errorf("expected array item type, got %q", ds.Items.Type)
	}
	if ds.Items.Items == nil {
		t.Fatal("expected nested items schema for array[array[string]]")
	}
	if ds.Items.Items.Type != domain.SchemaTypeString {
		t.Errorf("expected string inner item type, got %q", ds.Items.Items.Type)
	}
}

func TestAdaptSchema_NilType(t *testing.T) {
	schema := &oas.Schema{
		Description: "no type",
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Type != "" {
		t.Errorf("expected empty type, got %q", ds.Type)
	}
	if ds.Description != "no type" {
		t.Errorf("unexpected description: %q", ds.Description)
	}
}

func TestAdaptSchemaRef_Nil(t *testing.T) {
	if adaptSchemaRef(nil) != nil {
		t.Error("expected nil for nil SchemaRef")
	}

	if adaptSchemaRef(&oas.SchemaRef{Value: nil}) != nil {
		t.Error("expected nil for SchemaRef with nil Value")
	}
}

func TestAdaptRequestBody_Nil(t *testing.T) {
	if adaptRequestBody(nil) != nil {
		t.Error("expected nil for nil RequestBodyRef")
	}

	if adaptRequestBody(&oas.RequestBodyRef{Value: nil}) != nil {
		t.Error("expected nil for RequestBodyRef with nil Value")
	}
}

func TestAdaptSchema_ArrayWithObjectItems(t *testing.T) {
	strType := oas.Types{"string"}
	intType := oas.Types{"integer"}
	objType := oas.Types{"object"}
	arrType := oas.Types{"array"}
	schema := &oas.Schema{
		Type: &arrType,
		Items: &oas.SchemaRef{
			Value: &oas.Schema{
				Type: &objType,
				Properties: oas.Schemas{
					"id":   &oas.SchemaRef{Value: &oas.Schema{Type: &intType}},
					"name": &oas.SchemaRef{Value: &oas.Schema{Type: &strType}},
				},
			},
		},
	}

	ds := adaptSchema(schema, schemaMaxDepth)

	if ds.Items == nil {
		t.Fatal("expected items schema")
	}
	if ds.Items.Type != domain.SchemaTypeObject {
		t.Errorf("expected object item type, got %q", ds.Items.Type)
	}
	if len(ds.Items.Properties) != 2 {
		t.Fatalf("expected 2 item properties, got %d", len(ds.Items.Properties))
	}
	if ds.Items.Properties["id"].Type != domain.SchemaTypeInteger {
		t.Errorf("expected integer for id, got %q", ds.Items.Properties["id"].Type)
	}
	if ds.Items.Properties["name"].Type != domain.SchemaTypeString {
		t.Errorf("expected string for name, got %q", ds.Items.Properties["name"].Type)
	}
}

func TestAdaptHeaders_Empty(t *testing.T) {
	if adaptHeaders(nil) != nil {
		t.Error("expected nil for nil headers")
	}
	if adaptHeaders(oas.Headers{}) != nil {
		t.Error("expected nil for empty headers")
	}
}

func TestAdaptHeaders_NilValue(t *testing.T) {
	headers := oas.Headers{
		"X-Foo": &oas.HeaderRef{Value: nil},
	}
	result := adaptHeaders(headers)
	if len(result) != 0 {
		t.Errorf("expected 0 headers (nil value skipped), got %d", len(result))
	}
}

func TestAdaptResponses_Nil(t *testing.T) {
	if adaptResponses(nil) != nil {
		t.Error("expected nil for nil Responses")
	}
}

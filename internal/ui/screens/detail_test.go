package screens_test

import (
	"regexp"
	"strings"
	"testing"

	"dazzle/internal/domain"
	"dazzle/internal/ui/screens"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func fullOperation() domain.Operation {
	return domain.Operation{
		ID:          "createPet",
		Path:        "/pets",
		Method:      domain.POST,
		Summary:     "Create a pet",
		Description: "Creates a new pet in the store",
		Parameters: []domain.Parameter{
			{
				Name:        "X-Request-ID",
				In:          domain.ParameterInHeader,
				Description: "Trace ID",
				Required:    false,
				Schema:      &domain.Schema{Type: domain.SchemaTypeString},
			},
		},
		RequestBody: &domain.RequestBody{
			Required: true,
			Content: map[string]domain.MediaType{
				"application/json": {
					Schema: &domain.Schema{
						Type:     domain.SchemaTypeObject,
						Required: []string{"name"},
						Properties: map[string]*domain.Schema{
							"name": {Type: domain.SchemaTypeString},
						},
					},
				},
			},
		},
		Responses: map[string]domain.Response{
			"201": {Description: "Pet created"},
			"400": {
				Description: "Validation error",
				Content: map[string]domain.MediaType{
					"application/json": {
						Schema: &domain.Schema{
							Type: domain.SchemaTypeObject,
							Properties: map[string]*domain.Schema{
								"message": {Type: domain.SchemaTypeString},
							},
						},
					},
				},
			},
		},
	}
}

func minimalOperation() domain.Operation {
	return domain.Operation{
		ID:     "healthCheck",
		Path:   "/health",
		Method: domain.GET,
	}
}

func renderDetail(op domain.Operation) string {
	d := screens.NewDetailPanel(80, 40)
	d.SetOperation(op)
	return d.View()
}

func TestDetailPanel_RendersMethodAndPath(t *testing.T) {
	view := renderDetail(fullOperation())

	if !strings.Contains(view, "POST") {
		t.Error("expected POST in view")
	}
	if !strings.Contains(view, "/pets") {
		t.Error("expected /pets in view")
	}
}

func TestDetailPanel_RendersSummaryAndDescription(t *testing.T) {
	view := renderDetail(fullOperation())

	if !strings.Contains(view, "Create a pet") {
		t.Error("expected summary in view")
	}
	// Glamour inserts ANSI codes between words; strip them before checking.
	plain := ansiRe.ReplaceAllString(view, "")
	if !strings.Contains(plain, "Creates a new pet in the store") {
		t.Error("expected description in view")
	}
}

func TestDetailPanel_RendersParameters(t *testing.T) {
	view := renderDetail(fullOperation())

	if !strings.Contains(view, "X-Request-ID") {
		t.Error("expected parameter name in view")
	}
	if !strings.Contains(view, "header") {
		t.Error("expected parameter location in view")
	}
	if !strings.Contains(view, "string") {
		t.Error("expected parameter type in view")
	}
	// Glamour inserts ANSI codes between words; strip them before checking.
	plain := ansiRe.ReplaceAllString(view, "")
	if !strings.Contains(plain, "Trace ID") {
		t.Error("expected parameter description in view")
	}
}

func TestDetailPanel_RendersRequestBody(t *testing.T) {
	view := renderDetail(fullOperation())

	if !strings.Contains(view, "Request Body") {
		t.Error("expected Request Body section in view")
	}
	if !strings.Contains(view, "required") {
		t.Error("expected required indicator in view")
	}
	if !strings.Contains(view, "application/json") {
		t.Error("expected content type in view")
	}
	if !strings.Contains(view, "name") {
		t.Error("expected schema property name in view")
	}
}

func TestDetailPanel_RendersResponses(t *testing.T) {
	view := renderDetail(fullOperation())

	if !strings.Contains(view, "201") {
		t.Error("expected 201 status code in view")
	}
	if !strings.Contains(view, "Pet created") {
		t.Error("expected 201 description in view")
	}
	if !strings.Contains(view, "400") {
		t.Error("expected 400 status code in view")
	}
	if !strings.Contains(view, "Validation error") {
		t.Error("expected 400 description in view")
	}
}

func TestDetailPanel_MinimalOperation(t *testing.T) {
	view := renderDetail(minimalOperation())

	if !strings.Contains(view, "GET") {
		t.Error("expected GET in view")
	}
	if !strings.Contains(view, "/health") {
		t.Error("expected /health in view")
	}
	if !strings.Contains(view, "None") {
		t.Error("expected None for empty sections")
	}
}

func TestDetailPanel_NoOperationShowsPlaceholder(t *testing.T) {
	d := screens.NewDetailPanel(80, 40)
	view := d.View()

	if !strings.Contains(view, "Select an operation") {
		t.Error("expected placeholder text when no operation is set")
	}
}

func TestDetailPanel_ScrollbarVisibility(t *testing.T) {
	// Tall panel — all content fits, no scrollbar thumb.
	d := screens.NewDetailPanel(80, 200)
	d.SetOperation(fullOperation())
	view := d.View()
	if strings.Contains(view, "┃") {
		t.Error("expected no scrollbar thumb when all content is visible")
	}

	// Short panel — content overflows, scrollbar thumb appears.
	d.SetSize(80, 5)
	view = d.View()
	if !strings.Contains(view, "┃") {
		t.Error("expected scrollbar thumb when content overflows")
	}
}

func TestDetailPanel_SetSizeUpdatesContent(t *testing.T) {
	d := screens.NewDetailPanel(80, 40)
	op := fullOperation()
	d.SetOperation(op)

	d.SetSize(120, 50)
	view := d.View()

	if !strings.Contains(view, "/pets") {
		t.Error("expected content after resize")
	}
}

func TestDetailPanel_RequiredParameter(t *testing.T) {
	op := domain.Operation{
		ID:     "getPet",
		Path:   "/pets/{id}",
		Method: domain.GET,
		Parameters: []domain.Parameter{
			{
				Name:     "id",
				In:       domain.ParameterInPath,
				Required: true,
				Schema:   &domain.Schema{Type: domain.SchemaTypeInteger, Format: "int64"},
			},
		},
	}

	view := renderDetail(op)

	if !strings.Contains(view, "required") {
		t.Error("expected required indicator for path param")
	}
	if !strings.Contains(view, "integer") {
		t.Error("expected integer type")
	}
}

func TestDetailPanel_ArrayResponseSchema(t *testing.T) {
	op := domain.Operation{
		ID:     "listPets",
		Path:   "/pets",
		Method: domain.GET,
		Responses: map[string]domain.Response{
			"200": {
				Description: "A list of pets",
				Content: map[string]domain.MediaType{
					"application/json": {
						Schema: &domain.Schema{
							Type: domain.SchemaTypeArray,
							Items: &domain.Schema{
								Type: domain.SchemaTypeObject,
								Properties: map[string]*domain.Schema{
									"id":   {Type: domain.SchemaTypeInteger},
									"name": {Type: domain.SchemaTypeString},
								},
							},
						},
					},
				},
			},
		},
	}

	view := renderDetail(op)

	if !strings.Contains(view, "array[object]") {
		t.Error("expected array[object] type in view")
	}
	if !strings.Contains(view, "id") {
		t.Error("expected id property in view")
	}
}

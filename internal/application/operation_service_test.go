package application_test

import (
	"testing"

	"dazzle/internal/application"
	"dazzle/internal/domain"
)

func newTestOperations() []domain.Operation {
	return []domain.Operation{
		{ID: "createUser", Path: "/users", Method: domain.POST, Summary: "Create user", Tags: []string{"users"}},
		{ID: "listPets", Path: "/pets", Method: domain.GET, Summary: "List pets", Tags: []string{"pets"}},
		{ID: "listUsers", Path: "/users", Method: domain.GET, Summary: "List users", Tags: []string{"users"}},
		{ID: "deletePet", Path: "/pets/{id}", Method: domain.DELETE, Summary: "Delete a pet", Tags: []string{"pets"}},
	}
}

func TestOperationService_ListOperations(t *testing.T) {
	svc := application.NewOperationService()
	ops := newTestOperations()
	spec := &domain.Spec{Operations: ops}

	result := svc.ListOperations(spec)
	if len(result) != 4 {
		t.Errorf("expected 4 operations, got %d", len(result))
	}
}

func TestOperationService_SortOperations(t *testing.T) {
	svc := application.NewOperationService()
	sorted := svc.SortOperations(newTestOperations())

	expected := []struct {
		path   string
		method domain.HTTPMethod
	}{
		{"/pets", domain.GET},
		{"/pets/{id}", domain.DELETE},
		{"/users", domain.GET},
		{"/users", domain.POST},
	}

	if len(sorted) != len(expected) {
		t.Fatalf("expected %d operations, got %d", len(expected), len(sorted))
	}

	for i, want := range expected {
		if sorted[i].Path != want.path || sorted[i].Method != want.method {
			t.Errorf("position %d: expected %s %s, got %s %s",
				i, want.method, want.path, sorted[i].Method, sorted[i].Path)
		}
	}
}

func TestOperationService_FilterByQuery(t *testing.T) {
	svc := application.NewOperationService()
	result := svc.FilterOperations(newTestOperations(), domain.OperationFilter{Query: "pets"})

	if len(result) != 2 {
		t.Errorf("expected 2 operations matching 'pets', got %d", len(result))
	}
}

func TestOperationService_FilterByMethod(t *testing.T) {
	svc := application.NewOperationService()
	result := svc.FilterOperations(newTestOperations(), domain.OperationFilter{Method: domain.GET})

	if len(result) != 2 {
		t.Errorf("expected 2 GET operations, got %d", len(result))
	}
}

func TestOperationService_FilterByTag(t *testing.T) {
	svc := application.NewOperationService()
	result := svc.FilterOperations(newTestOperations(), domain.OperationFilter{Tags: []string{"users"}})

	if len(result) != 2 {
		t.Errorf("expected 2 operations with tag 'users', got %d", len(result))
	}
}

func TestOperationService_FilterCombined(t *testing.T) {
	svc := application.NewOperationService()
	result := svc.FilterOperations(newTestOperations(), domain.OperationFilter{
		Query:  "list",
		Method: domain.GET,
		Tags:   []string{"pets"},
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(result))
	}
	if result[0].ID != "listPets" {
		t.Errorf("expected listPets, got %s", result[0].ID)
	}
}

func TestOperationService_FilterNoMatch(t *testing.T) {
	svc := application.NewOperationService()
	result := svc.FilterOperations(newTestOperations(), domain.OperationFilter{Query: "nonexistent"})

	if len(result) != 0 {
		t.Errorf("expected 0 operations, got %d", len(result))
	}
}

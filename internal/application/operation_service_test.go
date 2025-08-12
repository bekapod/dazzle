package application

import (
	"testing"

	"dazzle/internal/domain"
)

// Mock repository for testing
type mockRepository struct {
	operations []domain.Operation
}

func (m *mockRepository) GetOperations() ([]domain.Operation, error) {
	return m.operations, nil
}


func TestOperationService_ListOperations(t *testing.T) {
	mockOps := []domain.Operation{
		{Path: "/users", Method: domain.GET, Summary: "List users"},
		{Path: "/users", Method: domain.POST, Summary: "Create user"},
		{Path: "/api/v1/health", Method: domain.GET, Summary: "Health check"},
	}

	repo := &mockRepository{operations: mockOps}
	service := NewOperationService(repo)

	operations, err := service.ListOperations()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(operations) != 3 {
		t.Errorf("expected 3 operations, got %d", len(operations))
	}

	// Test sorting: should be by path first, then by method
	expected := []struct {
		path   string
		method domain.HTTPMethod
	}{
		{"/api/v1/health", domain.GET},
		{"/users", domain.GET},
		{"/users", domain.POST},
	}

	for i, exp := range expected {
		if operations[i].Path != exp.path {
			t.Errorf("operation %d: expected path %s, got %s", i, exp.path, operations[i].Path)
		}
		if operations[i].Method != exp.method {
			t.Errorf("operation %d: expected method %s, got %s", i, exp.method, operations[i].Method)
		}
	}
}








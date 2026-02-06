package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"dazzle/internal/domain"
	"dazzle/internal/ui/screens"
)

type stubOpService struct{}

func (s *stubOpService) ListOperations(spec *domain.Spec) []domain.Operation { return spec.Operations }
func (s *stubOpService) FilterOperations(ops []domain.Operation, _ domain.OperationFilter) []domain.Operation {
	return ops
}
func (s *stubOpService) SortOperations(ops []domain.Operation) []domain.Operation { return ops }

func testSpec() *domain.Spec {
	return &domain.Spec{
		Info: domain.SpecInfo{Title: "Petstore API"},
		Operations: []domain.Operation{
			{ID: "listPets", Path: "/pets", Method: domain.GET, Summary: "List all pets"},
			{ID: "createPet", Path: "/pets", Method: domain.POST, Summary: "Create a pet"},
			{ID: "deletePet", Path: "/pets/{id}", Method: domain.DELETE, Summary: "Delete a pet"},
		},
	}
}

func TestOperationsScreen_Name(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	if s.Name() != "operations" {
		t.Errorf("expected name 'operations', got %q", s.Name())
	}
}

func TestOperationsScreen_Render(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()

	if !strings.Contains(view, "/pets") {
		t.Error("expected /pets in view")
	}
	if !strings.Contains(view, "GET") {
		t.Error("expected GET in view")
	}
	if !strings.Contains(view, "POST") {
		t.Error("expected POST in view")
	}
	if !strings.Contains(view, "DELETE") {
		t.Error("expected DELETE in view")
	}
}

func TestOperationsScreen_ShowsTitle(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if !strings.Contains(view, "Petstore API") {
		t.Error("expected spec title in view")
	}
}

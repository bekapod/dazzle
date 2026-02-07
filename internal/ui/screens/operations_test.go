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

// typeFilter enters filter mode and types the given text. It drains the
// command from the final keystroke so that bubbles/list's async filtering
// takes effect. Only the last command is drained — the filter captures the
// model state at creation, so intermediate filter commands are superseded.
func typeFilter(s *screens.OperationsScreen, text string) {
	s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	var lastCmd tea.Cmd
	for _, r := range text {
		_, lastCmd = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	drainCmd(s, lastCmd)
}

// drainCmd executes a command and feeds the result back into the model.
// Batch sub-commands are run concurrently and all results are collected.
func drainCmd(s *screens.OperationsScreen, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	result := cmd()
	if result == nil {
		return
	}
	if batch, ok := result.(tea.BatchMsg); ok {
		ch := make(chan tea.Msg, len(batch))
		n := 0
		for _, c := range batch {
			if c == nil {
				continue
			}
			n++
			go func(fn tea.Cmd) { ch <- fn() }(c)
		}
		for range n {
			if m := <-ch; m != nil {
				s.Update(m)
			}
		}
		return
	}
	s.Update(result)
}

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

func TestOperationsScreen_EmptySummary(t *testing.T) {
	spec := &domain.Spec{
		Info: domain.SpecInfo{Title: "Test"},
		Operations: []domain.Operation{
			{ID: "noSummary", Path: "/empty", Method: domain.GET, Summary: ""},
			{ID: "withSummary", Path: "/full", Method: domain.POST, Summary: "Has a summary"},
		},
	}
	s := screens.NewOperationsScreen(spec, &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if !strings.Contains(view, "/empty") {
		t.Error("expected /empty in view")
	}
	if !strings.Contains(view, "/full") {
		t.Error("expected /full in view")
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

func TestOperationsScreen_DetailShowsSelectedOperation(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	view := s.View()

	// First operation (listPets GET /pets) is selected by default.
	if !strings.Contains(view, "List all pets") {
		t.Error("expected selected operation's summary in detail panel")
	}
}

func TestOperationsScreen_TabSwitchesFocus(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// First operation is selected; detail shows its summary.
	view := s.View()
	if !strings.Contains(view, "List all pets") {
		t.Fatal("expected first operation selected initially")
	}

	// Tab switches focus to detail. Down key should now scroll
	// the viewport rather than changing the list selection.
	s.Update(tea.KeyMsg{Type: tea.KeyTab})
	s.Update(tea.KeyMsg{Type: tea.KeyDown})

	view = s.View()
	// Detail should still show the first operation (list didn't move).
	if !strings.Contains(view, "List all pets") {
		t.Error("expected first operation's summary to remain after down-key in detail focus")
	}

	// Tab back to list. Down key should now move selection.
	s.Update(tea.KeyMsg{Type: tea.KeyTab})
	s.Update(tea.KeyMsg{Type: tea.KeyDown})

	view = s.View()
	if !strings.Contains(view, "Create a pet") {
		t.Error("expected second operation after down-key in list focus")
	}
}

func TestOperationsScreen_NavigationUpdatesDetail(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Move selection down to the second operation.
	s.Update(tea.KeyMsg{Type: tea.KeyDown})
	view := s.View()

	if !strings.Contains(view, "Create a pet") {
		t.Error("expected second operation's summary in detail after navigating down")
	}
}

func TestOperationsScreen_FilterUpdatesDetail(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Initially index 0 shows listPets.
	view := s.View()
	if !strings.Contains(view, "List all pets") {
		t.Fatal("expected first operation initially")
	}

	// Filter to match only deletePet. Index stays 0 but the item changes.
	typeFilter(s, "delete")

	view = s.View()
	if !strings.Contains(view, "Delete a pet") {
		t.Error("expected detail to update when filter changes the item at index 0")
	}
}

func TestOperationsScreen_EmptyFilterClearsDetail(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Detail shows the first operation initially.
	view := s.View()
	if !strings.Contains(view, "List all pets") {
		t.Fatal("expected first operation initially")
	}

	// Filter to a term that matches nothing.
	typeFilter(s, "zzzznotfound")

	view = s.View()
	// Detail should clear and show the placeholder.
	if strings.Contains(view, "List all pets") {
		t.Error("expected detail to clear when filter matches nothing")
	}
	if !strings.Contains(view, "Select an operation") {
		t.Error("expected placeholder text when no operation is selected")
	}
}

func TestOperationsScreen_ResizePropagates(t *testing.T) {
	s := screens.NewOperationsScreen(testSpec(), &stubOpService{})
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view80 := s.View()
	if view80 == "" {
		t.Fatal("expected non-empty view at 80 width")
	}

	// Resize to wider terminal.
	s.Update(tea.WindowSizeMsg{Width: 160, Height: 50})
	view160 := s.View()

	if view160 == "" {
		t.Fatal("expected non-empty view at 160 width")
	}
	if view80 == view160 {
		t.Error("expected view to change after resize")
	}
}

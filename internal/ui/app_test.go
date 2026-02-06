package ui_test

import (
	"context"
	"errors"
	"testing"

	"dazzle/internal/domain"
	"dazzle/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type stubSpecService struct {
	spec *domain.Spec
	err  error
}

func (s *stubSpecService) LoadSpec(_ context.Context, _ string) (*domain.Spec, error) {
	return s.spec, s.err
}

func (s *stubSpecService) GetInfo(spec *domain.Spec) domain.SpecInfo {
	return spec.Info
}

type stubOperationService struct{}

func (s *stubOperationService) ListOperations(spec *domain.Spec) []domain.Operation {
	return spec.Operations
}

func (s *stubOperationService) FilterOperations(ops []domain.Operation, _ domain.OperationFilter) []domain.Operation {
	return ops
}

func (s *stubOperationService) SortOperations(ops []domain.Operation) []domain.Operation {
	return ops
}

func TestAppModel_Init(t *testing.T) {
	svc := &stubSpecService{spec: &domain.Spec{}}
	app := ui.NewAppModel(context.Background(), svc, &stubOperationService{}, "test.yaml")

	cmd := app.Init()
	if cmd == nil {
		t.Fatal("expected init to return a command")
	}
}

func TestAppModel_Quit(t *testing.T) {
	svc := &stubSpecService{spec: &domain.Spec{}}
	app := ui.NewAppModel(context.Background(), svc, &stubOperationService{}, "test.yaml")

	updated, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if updated == nil {
		t.Fatal("expected model after quit")
	}

	// cmd should be tea.Quit
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestAppModel_SpecLoadedError(t *testing.T) {
	svc := &stubSpecService{spec: &domain.Spec{}}
	app := ui.NewAppModel(context.Background(), svc, &stubOperationService{}, "test.yaml")

	// Simulate window size first so view renders properly
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Simulate spec load error
	app.Update(ui.SpecLoadedMsg{Err: errors.New("bad spec")})

	view := app.View()
	if view == "" {
		t.Error("expected non-empty view after error")
	}
}

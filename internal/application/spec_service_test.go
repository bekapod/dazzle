package application_test

import (
	"context"
	"errors"
	"testing"

	"dazzle/internal/application"
	"dazzle/internal/domain"
)

type mockSpecRepo struct {
	spec *domain.Spec
	err  error
}

func (m *mockSpecRepo) Load(_ context.Context, _ string) (*domain.Spec, error) {
	return m.spec, m.err
}

func TestSpecService_LoadSpec(t *testing.T) {
	want := &domain.Spec{
		Info: domain.SpecInfo{Title: "Test API", Version: "1.0.0"},
	}

	svc := application.NewSpecService(&mockSpecRepo{spec: want})

	got, err := svc.LoadSpec(context.Background(), "test.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Info.Title != want.Info.Title {
		t.Errorf("expected title %q, got %q", want.Info.Title, got.Info.Title)
	}
}

func TestSpecService_LoadSpec_Error(t *testing.T) {
	svc := application.NewSpecService(&mockSpecRepo{err: errors.New("not found")})

	_, err := svc.LoadSpec(context.Background(), "bad.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSpecService_GetInfo(t *testing.T) {
	svc := application.NewSpecService(nil)
	spec := &domain.Spec{
		Info: domain.SpecInfo{
			Title:       "My API",
			Description: "A test API",
			Version:     "2.0.0",
		},
	}

	info := svc.GetInfo(spec)
	if info.Title != "My API" {
		t.Errorf("expected title 'My API', got %q", info.Title)
	}
	if info.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got %q", info.Version)
	}
}

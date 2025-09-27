package main

import (
	"io"
	"testing"
	"time"

	"dazzle/internal/application"
	"dazzle/internal/infrastructure"
	"dazzle/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func createTestModel(t *testing.T) *ui.AppModel {
	// Use file-based loading for tests
	repo, err := infrastructure.NewOpenAPIRepositoryFromFile("fixtures/openapi-spec.yaml")
	if err != nil {
		t.Fatal(err)
	}

	service := application.NewOperationService(repo)

	model, err := ui.NewAppModel(service)
	if err != nil {
		t.Fatal(err)
	}

	return model
}

func TestInitialOutput(t *testing.T) {
	m := createTestModel(t)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})
	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*2)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestOutputAfterScrolling(t *testing.T) {
	m := createTestModel(t)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("j"),
	})
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*2)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestOutputAfterFiltering(t *testing.T) {
	m := createTestModel(t)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))
	for _, r := range "/user" {
		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
		})
	}
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*5)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestDomainLogic(t *testing.T) {
	repo, err := infrastructure.NewOpenAPIRepositoryFromFile("fixtures/openapi-spec.yaml")
	if err != nil {
		t.Fatal(err)
	}

	service := application.NewOperationService(repo)

	operations, err := service.ListOperations()
	if err != nil {
		t.Fatal(err)
	}

	if len(operations) != 19 {
		t.Errorf("expected 19 operations, got %d", len(operations))
	}

	// Test that operations are sorted correctly
	if len(operations) > 5 {
		sample := operations[5]
		if sample.Path != "/pet/{petId}" {
			t.Errorf("expected path '/pet/{petId}', got '%s'", sample.Path)
		}

		if string(sample.Method) != "POST" {
			t.Errorf("expected method 'POST', got '%s'", sample.Method)
		}

		if sample.Summary != "Updates a pet in the store with form data" {
			t.Errorf("expected summary 'Updates a pet in the store with form data', got '%s'", sample.Summary)
		}
	}
}


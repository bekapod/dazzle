package main

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	oas "github.com/getkin/kin-openapi/openapi3"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func loadDoc() *oas.T {
	loader := oas.NewLoader()
	doc, err := loader.LoadFromFile("fixtures/openapi-spec.yaml")
	if err != nil {
		panic(err)
	}
	return doc
}

func TestInitialOutput(t *testing.T) {
	doc := loadDoc()
	m := createModel(doc)
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
	doc := loadDoc()
	m := createModel(doc)
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
	doc := loadDoc()
	m := createModel(doc)
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

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*2)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestModel(t *testing.T) {
	doc := loadDoc()
	tm := teatest.NewTestModel(t, createModel(doc), teatest.WithInitialTermSize(300, 100))
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})
	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second*2))
	m, ok := fm.(model)
	if !ok {
		t.Fatalf("final model has the wrong type: %T", fm)
	}
	if len(m.operations.Items()) != 19 {
		t.Errorf("m.operations.Items len != 19: %d", len(m.operations.Items()))
	}

	sample := m.operations.Items()[5].(operation)
	if sample.path != "/pet/{petId}" {
		t.Errorf("m.operations.Items[5].path != /pet/{petId}: %s", sample.path)
	}

	if sample.method != "POST" {
		t.Errorf("m.operations.Items[5].method != POST: %s", sample.method)
	}

	if sample.summary != "Updates a pet in the store with form data" {
		t.Errorf("m.operations.Items[5].summary != Updates a pet in the store with form data: %s", sample.summary)
	}

	if len(sample.tags) != 1 {
		t.Errorf("m.operations.Items[5].tags len != 1: %d", len(sample.tags))
	}

	if sample.tags[0] != "pet" {
		t.Errorf("m.operations.Items[5].tags[0] != pet: %s", sample.tags[0])
	}
}

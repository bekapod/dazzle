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
	if len(m.paths.Items()) != 13 {
		t.Errorf("m.paths.Items len != 13: %d", len(m.paths.Items()))
	}
	if m.paths.Items()[5].(path).path != "/store/inventory" {
		t.Errorf("m.paths.Items[5].path != /store/inventory: %s", m.paths.Items()[5].(path).path)
	}
	if m.paths.Items()[9].(path).path != "/user/createWithList" {
		t.Errorf("m.paths.Items[9].path != /user/createWithList: %s", m.paths.Items()[9].(path).path)
	}
}

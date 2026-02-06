package screens_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"dazzle/internal/ui/screens"
)

func TestWelcomeScreen_Name(t *testing.T) {
	s := screens.NewWelcomeScreen()
	if s.Name() != "welcome" {
		t.Errorf("expected name 'welcome', got %q", s.Name())
	}
}

func TestWelcomeScreen_Loading(t *testing.T) {
	s := screens.NewWelcomeScreen()
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := s.View()
	if !strings.Contains(view, "Loading") {
		t.Error("expected loading message in view")
	}
}

func TestWelcomeScreen_Error(t *testing.T) {
	s := screens.NewWelcomeScreen()
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	s.SetError(errors.New("file not found"))

	view := s.View()
	if !strings.Contains(view, "file not found") {
		t.Error("expected error message in view")
	}
}

func TestWelcomeScreen_Loaded(t *testing.T) {
	s := screens.NewWelcomeScreen()
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	s.SetLoaded()

	view := s.View()
	if !strings.Contains(view, "Ready") {
		t.Error("expected ready message in view")
	}
}

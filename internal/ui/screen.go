package ui

import (
	"dazzle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents a navigable UI screen.
type Screen interface {
	tea.Model
	Name() string
}

// SpecLoadedMsg is sent when a spec finishes loading.
type SpecLoadedMsg struct {
	Spec *domain.Spec
	Err  error
}

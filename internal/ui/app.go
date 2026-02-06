package ui

import (
	"context"

	"dazzle/internal/domain"
	"dazzle/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// AppModel is the root Bubbletea model managing screen navigation.
type AppModel struct {
	ctx        context.Context
	specSvc    domain.SpecService
	opSvc      domain.OperationService
	specSource string

	spec   *domain.Spec
	screen Screen
	width  int
	height int
}

func NewAppModel(
	ctx context.Context,
	specSvc domain.SpecService,
	opSvc domain.OperationService,
	specSource string,
) *AppModel {
	return &AppModel{
		ctx:        ctx,
		specSvc:    specSvc,
		opSvc:      opSvc,
		specSource: specSource,
		screen:     screens.NewWelcomeScreen(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	return m.loadSpec()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		updated, cmd := m.screen.Update(msg)
		m.screen = updated.(Screen)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "q" && m.screen.Name() == "welcome" {
			return m, tea.Quit
		}

	case SpecLoadedMsg:
		return m.handleSpecLoaded(msg)
	}

	updated, cmd := m.screen.Update(msg)
	m.screen = updated.(Screen)
	return m, cmd
}

func (m *AppModel) View() string {
	return m.screen.View()
}

func (m *AppModel) handleSpecLoaded(msg SpecLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		if ws, ok := m.screen.(*screens.WelcomeScreen); ok {
			ws.SetError(msg.Err)
		}
		return m, nil
	}

	m.spec = msg.Spec

	opsScreen := screens.NewOperationsScreen(m.spec, m.opSvc)
	m.screen = opsScreen

	// Send the current window size to the new screen
	return m.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
}

func (m *AppModel) loadSpec() tea.Cmd {
	return func() tea.Msg {
		spec, err := m.specSvc.LoadSpec(m.ctx, m.specSource)
		return SpecLoadedMsg{Spec: spec, Err: err}
	}
}

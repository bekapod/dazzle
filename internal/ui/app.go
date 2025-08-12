package ui

import (
	"dazzle/internal/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Padding(0, 0)

// AppModel represents the main application model
type AppModel struct {
	operationList *OperationList
}

// NewAppModel creates a new application model
func NewAppModel(service domain.OperationService) (*AppModel, error) {
	operationList, err := NewOperationList(service)
	if err != nil {
		return nil, err
	}

	return &AppModel{
		operationList: operationList,
	}, nil
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.operationList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.operationList.FilterState() == list.Filtering {
			break
		}

		switch {
		// Add any additional key handling here
		}
	}

	_, cmd := m.operationList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *AppModel) View() string {
	return appStyle.Render(m.operationList.View())
}
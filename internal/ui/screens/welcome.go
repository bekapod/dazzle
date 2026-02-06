package screens

import (
	"fmt"

	"dazzle/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WelcomeScreen displays a loading/error state while the spec loads.
type WelcomeScreen struct {
	width   int
	height  int
	loading bool
	err     error
}

func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{loading: true}
}

func (s *WelcomeScreen) Name() string { return "welcome" }

func (s *WelcomeScreen) Init() tea.Cmd { return nil }

func (s *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		s.width = msg.Width
		s.height = msg.Height
	}
	return s, nil
}

func (s *WelcomeScreen) View() string {
	var content string

	switch {
	case s.err != nil:
		content = lipgloss.JoinVertical(lipgloss.Left,
			styles.Title.Render("dazzle"),
			"",
			styles.Error.Render(fmt.Sprintf("Error: %v", s.err)),
			"",
			styles.Muted.Render("Press q to quit."),
		)
	case s.loading:
		content = lipgloss.JoinVertical(lipgloss.Left,
			styles.Title.Render("dazzle"),
			"",
			styles.Subtitle.Render("Loading spec..."),
		)
	default:
		content = lipgloss.JoinVertical(lipgloss.Left,
			styles.Title.Render("dazzle"),
			"",
			styles.Subtitle.Render("Ready."),
		)
	}

	return lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, content)
}

func (s *WelcomeScreen) SetError(err error) {
	s.err = err
	s.loading = false
}

func (s *WelcomeScreen) SetLoaded() {
	s.loading = false
}

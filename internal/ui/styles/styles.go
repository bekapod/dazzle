package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Text)

	Subtitle = lipgloss.NewStyle().
			Foreground(Subtext1).
			Italic(true)

	Error = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)

	Muted = lipgloss.NewStyle().
		Foreground(Overlay1)
)

// Method returns a styled string for an HTTP method.
func Method(method string) string {
	return lipgloss.NewStyle().
		Foreground(MethodColor(method)).
		Bold(true).
		Render(method)
}

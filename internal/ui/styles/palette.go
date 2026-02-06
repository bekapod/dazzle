package styles

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha (dark) / Latte (light) adaptive palette.
var (
	Text     = lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cad3f5"}
	Subtext1 = lipgloss.AdaptiveColor{Light: "#5c5f77", Dark: "#b8c0e0"}
	Subtext0 = lipgloss.AdaptiveColor{Light: "#6c6f85", Dark: "#a5adcb"}
	Overlay1 = lipgloss.AdaptiveColor{Light: "#8c8fa1", Dark: "#8087a2"}
	Surface1 = lipgloss.AdaptiveColor{Light: "#bcc0cc", Dark: "#494d64"}
	Base     = lipgloss.AdaptiveColor{Light: "#eff1f5", Dark: "#24273a"}

	Green  = lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6da95"}
	Purple = lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#c6a0f6"}
	Orange = lipgloss.AdaptiveColor{Light: "#fe640b", Dark: "#f5a97f"}
	Yellow = lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#eed49f"}
	Red    = lipgloss.AdaptiveColor{Light: "#e64553", Dark: "#ed8796"}
	Blue   = lipgloss.AdaptiveColor{Light: "#1e66f5", Dark: "#8aadf4"}
)

// MethodColor returns the color for an HTTP method.
func MethodColor(method string) lipgloss.AdaptiveColor {
	switch method {
	case "GET":
		return Green
	case "POST":
		return Purple
	case "PUT":
		return Orange
	case "PATCH":
		return Yellow
	case "DELETE":
		return Red
	default:
		return Text
	}
}

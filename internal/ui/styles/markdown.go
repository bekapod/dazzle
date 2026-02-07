package styles

import (
	"github.com/charmbracelet/glamour/ansi"
	glamourstyles "github.com/charmbracelet/glamour/styles"
)

// MarkdownStyle returns a glamour style using the Catppuccin Mocha palette with
// document margins stripped so rendered markdown fits inline in detail panels.
func MarkdownStyle() ansi.StyleConfig {
	s := glamourstyles.DarkStyleConfig

	// Strip document margins/padding for inline use.
	s.Document.BlockPrefix = ""
	s.Document.BlockSuffix = ""
	s.Document.Margin = uintPtr(0)
	s.Document.Color = nil // inherit terminal foreground

	// Remap accent colours to Catppuccin Mocha.
	s.Link.Color = stringPtr(Blue.Dark)
	s.LinkText.Color = stringPtr(Blue.Dark)
	s.Code.Color = stringPtr(Red.Dark)
	s.Code.BackgroundColor = stringPtr(Surface0.Dark)
	s.CodeBlock.Margin = uintPtr(0)
	s.Heading.Color = stringPtr(Blue.Dark)
	s.H1.Color = stringPtr(Blue.Dark)
	s.H1.BackgroundColor = nil

	return s
}

func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

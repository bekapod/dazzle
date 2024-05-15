package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	bullet = "•"
)

type Palette struct {
	Text      lipgloss.AdaptiveColor
	SubText1  lipgloss.AdaptiveColor
	SubText0  lipgloss.AdaptiveColor
	Overlay1  lipgloss.AdaptiveColor
	Surface1  lipgloss.AdaptiveColor
	Base      lipgloss.AdaptiveColor
	Rosewater lipgloss.AdaptiveColor
	Lavendar  lipgloss.AdaptiveColor
	Methods   MethodPalette
}

type MethodPalette struct {
	Get    lipgloss.AdaptiveColor
	Post   lipgloss.AdaptiveColor
	Put    lipgloss.AdaptiveColor
	Patch  lipgloss.AdaptiveColor
	Delete lipgloss.AdaptiveColor
}

func NewPalette() (p Palette) {
	return Palette{
		Text:      lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cad3f5"},
		SubText1:  lipgloss.AdaptiveColor{Light: "#5c5f77", Dark: "#b8c0e0"},
		SubText0:  lipgloss.AdaptiveColor{Light: "#6c6f85", Dark: "#a5adcb"},
		Overlay1:  lipgloss.AdaptiveColor{Light: "#8c8fa1", Dark: "#8087a2"},
		Surface1:  lipgloss.AdaptiveColor{Light: "#bcc0cc", Dark: "#494d64"},
		Base:      lipgloss.AdaptiveColor{Light: "#eff1f5", Dark: "#1a1b26"},
		Rosewater: lipgloss.AdaptiveColor{Light: "#dc8a78", Dark: "#f4dbd6"},
		Lavendar:  lipgloss.AdaptiveColor{Light: "#7287fd", Dark: "#b7bdf8"},
		Methods: MethodPalette{
			Get:    lipgloss.AdaptiveColor{Light: "#40a02b", Dark: "#a6da95"},
			Post:   lipgloss.AdaptiveColor{Light: "#8839ef", Dark: "#c6a0f6"},
			Put:    lipgloss.AdaptiveColor{Light: "#fe640b", Dark: "#f5a97f"},
			Patch:  lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#eed49f"},
			Delete: lipgloss.AdaptiveColor{Light: "#e64553", Dark: "#ed8796"},
		},
	}
}

func NewOperationListStyles() (s list.Styles) {
	palette := NewPalette()
	normalText := lipgloss.NewStyle().Foreground(palette.Text)

	s.Title = lipgloss.NewStyle().
		Foreground(palette.Base).
		Background(palette.Text).
		Bold(true).
		Padding(0, 1, 0, 1)
	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 1)

	s.StatusBar = normalText.Copy().
		Foreground(palette.Overlay1).
		Italic(true).
		Padding(0, 0, 0, 2).
		Margin(1, 0)
	s.StatusBarFilterCount = lipgloss.NewStyle().
		Padding(0, 1, 0, 1)

	s.PaginationStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(palette.Rosewater).
		SetString(bullet)
	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(palette.Surface1).
		SetString(bullet)

	s.HelpStyle = lipgloss.NewStyle().
		Padding(1, 0, 0, 2)

	s.FilterPrompt = normalText.Copy().
		Foreground(palette.SubText0).
		Padding(0, 0, 0, 1).
		Bold(true)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(palette.Rosewater)

	return s
}

func NewHelpStyles() (h help.Styles) {
	palette := NewPalette()
	normalText := lipgloss.NewStyle().Foreground(palette.Text)

	h.ShortKey = normalText.Copy().Foreground(palette.SubText1).Bold(true)
	h.ShortDesc = normalText.Copy().Foreground(palette.Overlay1)
	h.ShortSeparator = normalText.Copy().Foreground(palette.Surface1)
	h.FullKey = normalText.Copy().Foreground(palette.SubText1).Bold(true)
	h.FullDesc = normalText.Copy().Foreground(palette.Overlay1)
	h.FullSeparator = normalText.Copy().Foreground(palette.Surface1)

	return h
}

type OperationItemStyles struct {
	NormalItem    lipgloss.Style
	NormalMethod  lipgloss.Style
	NormalPath    lipgloss.Style
	NormalSummary lipgloss.Style

	SelectedItem    lipgloss.Style
	SelectedMethod  lipgloss.Style
	SelectedPath    lipgloss.Style
	SelectedSummary lipgloss.Style

	DimmedItem    lipgloss.Style
	DimmedMethod  lipgloss.Style
	DimmedPath    lipgloss.Style
	DimmedSummary lipgloss.Style

	FilterMatch lipgloss.Style
}

func NewOperationItemStyles() (s OperationItemStyles) {
	palette := NewPalette()
	normalText := lipgloss.NewStyle().Foreground(palette.Text)

	s.NormalItem = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Border(lipgloss.HiddenBorder(), false, false, false, true)

	s.NormalMethod = normalText.Copy().
		Bold(true).
		Padding(0, 1, 0, 0)

	s.NormalPath = normalText

	s.NormalSummary = lipgloss.NewStyle().
		Foreground(palette.Overlay1)

	s.SelectedItem = s.NormalItem.Copy().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(palette.Rosewater)

	s.SelectedMethod = s.NormalMethod.Copy().
		Foreground(palette.Rosewater)

	s.SelectedPath = s.NormalPath.Copy().
		Foreground(palette.Rosewater)

	s.SelectedSummary = s.NormalSummary.Copy().
		Foreground(palette.Rosewater)

	s.DimmedItem = s.NormalItem.Copy().Faint(true)

	s.DimmedMethod = s.NormalMethod.Copy().Faint(true)

	s.DimmedPath = s.NormalPath.Copy().Faint(true)

	s.DimmedSummary = s.NormalSummary.Copy().Faint(true)

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

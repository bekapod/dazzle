package main

import (
	"github.com/charmbracelet/lipgloss"
)

type Palette struct {
	Text      lipgloss.AdaptiveColor
	SubText1  lipgloss.AdaptiveColor
	SubText0  lipgloss.AdaptiveColor
	Overlay1  lipgloss.AdaptiveColor
	Rosewater lipgloss.AdaptiveColor
}

func NewPalette() (p Palette) {
	return Palette{
		Text:      lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#cad3f5"},
		SubText1:  lipgloss.AdaptiveColor{Light: "#5c5f77", Dark: "#b8c0e0"},
		SubText0:  lipgloss.AdaptiveColor{Light: "#6c6f85", Dark: "#a5adcb"},
		Overlay1:  lipgloss.AdaptiveColor{Light: "#8c8fa1", Dark: "#8087a2"},
		Rosewater: lipgloss.AdaptiveColor{Light: "#dc8a78", Dark: "#f4dbd6"},
	}
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

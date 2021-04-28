package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	classicColors = []string{"#ccc", "#00f", "#080", "#f00", "#008", "#800", "#088", "#880", "#888", "#f00"}
	BoxStyles     [10]lipgloss.Style
	baseBoxStyle  = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#ccc")).
			Foreground(lipgloss.Color("#111")).
			Width(3).
			Align(lipgloss.Center)

	MarkStyle = lipgloss.NewStyle().
			Inherit(baseBoxStyle).
			Background(lipgloss.Color("#D22")).
			Foreground(lipgloss.Color("#222"))

	SuspectStyle = lipgloss.NewStyle().
			Inherit(baseBoxStyle).
			Background(lipgloss.Color("#DD2")).
			Foreground(lipgloss.Color("#222"))

	HiddenStyle = lipgloss.NewStyle().
			Inherit(baseBoxStyle).
			Background(lipgloss.Color("#777"))
)

func InitBoxStyles() {

	for i := range BoxStyles {
		BoxStyles[i] = lipgloss.NewStyle().
			Inherit(baseBoxStyle).
			Foreground(lipgloss.Color(classicColors[i]))
	}
	BoxStyles[9] = BoxStyles[9].
		Background(lipgloss.Color("#911")).
		Foreground(lipgloss.Color("#eee"))
}

package views

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle   = lipgloss.NewStyle().Padding(1, 2)
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#205b8f")).
			Background(lipgloss.Color("#f2f2f2")).
			Padding(0, 1)

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#205b8f", Dark: "#81a2be"}).
				Render
)

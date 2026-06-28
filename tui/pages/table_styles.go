package pages

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

func gruvboxTableStyles() table.Styles {
	s := table.DefaultStyles()

	s.Header = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.BG5).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.Green).
		Padding(0, 1)

	s.Cell = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(styles.FG0)

	s.Selected = lipgloss.NewStyle().
		Padding(0, 1).
		Background(styles.BG3).
		Foreground(styles.Green).
		Bold(true)

	return s
}

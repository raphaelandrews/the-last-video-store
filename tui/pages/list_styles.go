package pages

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

func gruvboxListStyles() list.Styles {
	s := list.DefaultStyles()

	s.Title = lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Padding(0, 1)

	s.TitleBar = lipgloss.NewStyle().
		Padding(0, 0, 1, 0)

	s.Spinner = lipgloss.NewStyle().Foreground(styles.Green)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true)
	s.FilterCursor = lipgloss.NewStyle().
		Foreground(styles.Green)
	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Bold(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Padding(0, 1)

	s.StatusEmpty = styles.DimTextStyle
	s.StatusBarActiveFilter = lipgloss.NewStyle().Foreground(styles.Green)
	s.StatusBarFilterCount = styles.DimTextStyle

	s.NoItems = lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Padding(2, 4)

	s.PaginationStyle = styles.DimTextStyle
	s.ActivePaginationDot = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	s.InactivePaginationDot = styles.DimTextStyle
	s.ArabicPagination = lipgloss.NewStyle().
		Foreground(styles.Grey1)
	s.DividerDot = lipgloss.NewStyle().
		Foreground(styles.Grey0)

	s.HelpStyle = list.DefaultStyles().HelpStyle.Foreground(styles.Grey1)

	return s
}

func enableListPagination(l *list.Model) {
	l.Paginator.PerPage = 15
	l.Paginator.Type = paginator.Dots
	l.Paginator.ActiveDot = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("●")
	l.Paginator.InactiveDot = lipgloss.NewStyle().Foreground(styles.Grey0).Render("○")
	l.Paginator.ArabicFormat = "%d / %d"
	l.SetShowPagination(true)
}

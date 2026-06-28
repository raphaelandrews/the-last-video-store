package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

var cardWidth = 23
var cardHeight = 8

func MovieCardView(movie models.MovieResponse, selected bool) string {
	borderColor := styles.BG5
	if selected {
		borderColor = styles.Green
	}

	border := lipgloss.NormalBorder()
	if selected {
		border = lipgloss.ThickBorder()
	}

	card := lipgloss.NewStyle().
		Width(cardWidth).
		Height(cardHeight).
		Border(border).
		BorderForeground(borderColor).
		Padding(0, 1)

	runes := []rune(movie.Title)
	title := string(runes)
	if len(runes) > 18 {
		title = string(runes[:17]) + "…"
	}

	titleColor := styles.FG1
	if selected {
		titleColor = styles.Green
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		Width(cardWidth - 4).
		Align(lipgloss.Center)

	yearStyle := lipgloss.NewStyle().
		Foreground(styles.Grey1)

	formatBadge := styles.FormatBadge(movie.Format)

	status := "[RENT]"
	statusColor := styles.Green
	if !movie.Available {
		status = "[OUT]"
		statusColor = styles.Red
	}
	if movie.IsNewRelease {
		status = "[NEW]"
		statusColor = styles.Yellow
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		yearStyle.Render(fmt.Sprintf("(%d)", movie.Year)),
		styles.StarRating(movie.Rating),
		lipgloss.JoinHorizontal(lipgloss.Center, formatBadge, "  ", statusStyle.Render(status)),
	)

	return card.Render(content)
}

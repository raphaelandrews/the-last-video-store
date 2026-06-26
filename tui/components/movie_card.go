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

	title := truncate(movie.Title, 18)
	if len(title) < len(movie.Title) {
		title = title[:len(title)-1] + "…"
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

	stars := RatingStars(movie.Rating)

	formatBadge := FormatBadge(movie.Format)

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
		stars,
		lipgloss.JoinHorizontal(lipgloss.Center, formatBadge, "  ", statusStyle.Render(status)),
	)

	return card.Render(content)
}

func RatingStars(rating float64) string {
	full := int(rating)
	half := 0
	if rating-float64(full) >= 0.5 {
		half = 1
	}
	empty := 5 - full - half

	s := ""
	for range full {
		s += "★"
	}
	for range half {
		s += "½"
	}
	for range empty {
		s += "☆"
	}
	return lipgloss.NewStyle().Foreground(styles.Yellow).Render(s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func FormatBadge(format string) string {
	switch format {
	case "VHS":
		return lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render("📼 VHS")
	case "DVD":
		return lipgloss.NewStyle().Foreground(styles.Aqua).Bold(true).Render("📀 DVD")
	case "Blu-ray":
		return lipgloss.NewStyle().Foreground(styles.Blue).Bold(true).Render("💿 BD")
	default:
		return lipgloss.NewStyle().Foreground(styles.Grey1).Render(format)
	}
}

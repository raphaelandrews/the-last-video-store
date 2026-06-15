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
	borderColor := styles.BorderDim
	bgColor := styles.Background
	if selected {
		borderColor = styles.Magenta
		bgColor = lipgloss.Color("#1A1A4E")
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
		Background(bgColor).
		Padding(0, 1)

	title := truncate(movie.Title, 18)
	if len(title) < len(movie.Title) {
		title = title[:len(title)-1] + "…"
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Background(bgColor).
		Bold(true).
		Width(cardWidth - 4).
		Align(lipgloss.Center)

	yearStyle := lipgloss.NewStyle().
		Foreground(styles.DimTextStyle.GetForeground()).
		Background(bgColor)

	stars := RatingStars(movie.Rating)

	formatBadge := FormatBadge(movie.Format)

	status := "[RENT]"
	statusColor := styles.NeonGreen
	if !movie.Available {
		status = "[OUT]"
		statusColor = styles.Error
	}
	if movie.IsNewRelease {
		status = "[NEW]"
		statusColor = styles.Yellow
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Background(bgColor).
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
	return styles.TitleStyle.Render(s)
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
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00")).Render("📼 VHS")
	case "DVD":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00AAFF")).Render("📀 DVD")
	case "Blu-ray":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#4444FF")).Render("💿 BD")
	default:
		return lipgloss.NewStyle().Foreground(styles.DimTextStyle.GetForeground()).Render(format)
	}
}

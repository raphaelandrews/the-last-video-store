package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/styles"
)

type AdminMoviesModel struct {
	movies   []models.MovieResponse
	selected int
	errMsg   string
}

type AdminMoviesRefreshMsg struct{}

func NewAdminMoviesModel() *AdminMoviesModel {
	return &AdminMoviesModel{selected: -1}
}

func (m *AdminMoviesModel) SetMovies(movies []models.MovieResponse) {
	m.movies = movies
}

func (m *AdminMoviesModel) SelectedMovie() *models.MovieResponse {
	if m.selected >= 0 && m.selected < len(m.movies) {
		return &m.movies[m.selected]
	}
	return nil
}

func (m *AdminMoviesModel) MoveUp() {
	if m.selected > 0 {
		m.selected--
	}
}

func (m *AdminMoviesModel) MoveDown() {
	if m.selected < len(m.movies)-1 {
		m.selected++
	}
}

func (m *AdminMoviesModel) View(width, height int) string {
	title := styles.HeadingStyle.Width(width).Align(lipgloss.Center).Render("🎬 MOVIE MANAGEMENT")

	header := styles.TextStyle.Bold(true).Render(
		"  TITLE                YEAR  GENRE     FORMAT    COPIES  STAFF PICK")

	var rows []string
	rows = append(rows, header)

	for i, mv := range m.movies {
		prefix := "  "
		style := styles.TextStyle
		if i == m.selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			style = styles.HighlightStyle
		}

		pick := "—"
		if mv.IsStaffPick {
			pick = styles.HeadingStyle.Render("★")
		}

		line := fmt.Sprintf("%-22s %-6d %-10s %-10s %-8s %s",
			prefix+truncateStr(mv.Title, 20),
			mv.Year,
			mv.Genre,
			components.FormatBadge(mv.Format),
			fmt.Sprintf("%d/%d", mv.CopiesAvailable, mv.CopiesTotal),
			pick,
		)
		rows = append(rows, style.Render(line))
	}

	footer := styles.TextStyle.Render("\n[A] Add Movie  [ENTER] Edit  [D] Delete  [S] Toggle Staff Pick  [ESC] Back")

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.JoinVertical(lipgloss.Left, title, content, footer)
}

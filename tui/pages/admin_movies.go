package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/styles"
)

type AdminMoviesModel struct {
	movies     []models.MovieResponse
	selected   int
	errMsg     string
	Page       int
	TotalPages int
	total      int
	PageSize   int
}

type AdminMoviesRefreshMsg struct{}

func NewAdminMoviesModel() *AdminMoviesModel {
	return &AdminMoviesModel{selected: -1, Page: 1, PageSize: 30}
}

func (m *AdminMoviesModel) SetMovies(movies []models.MovieResponse, total, page int) {
	m.movies = movies
	m.total = total
	m.Page = page
	if m.PageSize > 0 {
		m.TotalPages = (total + m.PageSize - 1) / m.PageSize
	}
	if m.selected >= len(m.movies) {
		m.selected = len(m.movies) - 1
	}
}

func (m *AdminMoviesModel) HasNextPage() bool { return m.Page < m.TotalPages }
func (m *AdminMoviesModel) HasPrevPage() bool { return m.Page > 1 }

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

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

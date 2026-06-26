package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type MovieGridModel struct {
	movies      []models.MovieResponse
	selectedRow int
	selectedCol int
	columns     int
	rows        int
}

func NewMovieGridModel() *MovieGridModel {
	return &MovieGridModel{}
}

func (m *MovieGridModel) SetMovies(movies []models.MovieResponse, width int) {
	m.movies = movies
	m.columns = (width - 2) / (cardWidth + 1)
	if m.columns < 1 {
		m.columns = 1
	}
	m.rows = (len(movies) + m.columns - 1) / m.columns
	if m.selectedRow >= m.rows {
		m.selectedRow = m.rows - 1
	}
	if m.selectedCol >= m.columns {
		m.selectedCol = m.columns - 1
	}
	if m.selectedRow < 0 {
		m.selectedRow = 0
	}
	if m.selectedCol < 0 {
		m.selectedCol = 0
	}
}

func (m *MovieGridModel) SelectedIndex() int {
	return m.selectedRow*m.columns + m.selectedCol
}

func (m *MovieGridModel) SelectByID(id string) {
	for i, mv := range m.movies {
		if mv.ID == id {
			m.selectedRow = i / m.columns
			m.selectedCol = i % m.columns
			return
		}
	}
}

func (m *MovieGridModel) MoveUp() {
	if m.selectedRow > 0 {
		m.selectedRow--
	}
}

func (m *MovieGridModel) MoveDown() {
	if m.selectedRow < m.rows-1 {
		m.selectedRow++
	}
}

func (m *MovieGridModel) MoveLeft() {
	if m.selectedCol > 0 {
		m.selectedCol--
	}
}

func (m *MovieGridModel) MoveRight() {
	maxCol := m.columns - 1
	lastRowCols := len(m.movies) - (m.rows-1)*m.columns
	if m.selectedRow == m.rows-1 && maxCol >= lastRowCols {
		maxCol = lastRowCols - 1
	}
	if m.selectedCol < maxCol {
		m.selectedCol++
	}
}

func (m *MovieGridModel) PageDown() {
	m.selectedRow += m.rows / 2
	if m.selectedRow >= m.rows {
		m.selectedRow = m.rows - 1
	}
}

func (m *MovieGridModel) PageUp() {
	m.selectedRow -= m.rows / 2
	if m.selectedRow < 0 {
		m.selectedRow = 0
	}
}

func (m *MovieGridModel) View(width int) string {
	if len(m.movies) == 0 {
		return lipgloss.Place(width, 10, lipgloss.Center, lipgloss.Center,
			styles.DimTextStyle.Render("No titles found"))
	}

	m.columns = (width - 2) / (cardWidth + 1)
	if m.columns < 1 {
		m.columns = 1
	}

	var rows []string
	for row := range m.rows {
		var cards []string
		for col := range m.columns {
			idx := row*m.columns + col
			if idx >= len(m.movies) {
				break
			}
			selected := idx == m.SelectedIndex()
			cards = append(cards, MovieCardView(m.movies[idx], selected))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cards...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

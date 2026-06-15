package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type SearchbarModel struct {
	input    textinput.Model
	results  []models.MovieResponse
	active   bool
	selected int
}

func NewSearchbarModel() *SearchbarModel {
	ti := textinput.New()
	ti.Placeholder = "Search movies..."
	ti.Width = 40
	ti.CharLimit = 60
	ti.Prompt = "🔍 "
	return &SearchbarModel{input: ti, selected: -1}
}

func (m *SearchbarModel) Focus() {
	m.active = true
	m.input.Focus()
}

func (m *SearchbarModel) Blur() {
	m.active = false
	m.input.Blur()
	m.results = nil
	m.selected = -1
}

func (m *SearchbarModel) Value() string { return m.input.Value() }

func (m *SearchbarModel) IsActive() bool { return m.active }

func (m *SearchbarModel) Update(msg tea.Msg) {
	if !m.active {
		return
	}
	m.input, _ = m.input.Update(msg)
}

func (m *SearchbarModel) SetResults(results []models.MovieResponse) {
	m.results = results
	m.selected = -1
}

func (m *SearchbarModel) SelectedMovieID() string {
	if m.selected >= 0 && m.selected < len(m.results) {
		return m.results[m.selected].ID
	}
	return ""
}

func (m *SearchbarModel) SelectedMovie() *models.MovieResponse {
	if m.selected >= 0 && m.selected < len(m.results) {
		return &m.results[m.selected]
	}
	return nil
}

func (m *SearchbarModel) MoveSelection(delta int) {
	if len(m.results) == 0 {
		return
	}
	m.selected += delta
	if m.selected >= len(m.results) {
		m.selected = 0
	}
	if m.selected < 0 {
		m.selected = len(m.results) - 1
	}
}

func (m *SearchbarModel) View() string {
	bar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Width(42).
		Render(m.input.View())

	if len(m.results) == 0 {
		return bar
	}

	var items []string
	for i, r := range m.results {
		prefix := "  "
		if i == m.selected {
			prefix = styles.HighlightStyle.Render("▸ ")
		}
		line := prefix + r.Title + styles.DimTextStyle.Render(" ("+itoa(r.Year)+") "+r.Format)
		items = append(items, line)
	}

	dropdown := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Width(42).
		Render(lipgloss.JoinVertical(lipgloss.Left, items...))

	return lipgloss.JoinVertical(lipgloss.Left, bar, dropdown)
}

func itoa(y int) string {
	if y == 0 {
		return "----"
	}
	s := ""
	for y > 0 {
		s = string(rune('0'+y%10)) + s
		y /= 10
	}
	return s
}

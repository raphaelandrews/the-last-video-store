package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/styles"
)

type AdminMoviesRefreshMsg struct{}

type MediaType string

const (
	MediaMovies MediaType = "movie"
	MediaSeries MediaType = "series"
	MediaGames  MediaType = "game"
)

var AllMediaTypes = []MediaType{MediaMovies, MediaSeries, MediaGames}

func (m MediaType) Label() string {
	switch m {
	case MediaSeries:
		return "📺 Series"
	case MediaGames:
		return "🕹️ Games"
	default:
		return "🎬 Movies"
	}
}

type AdminMoviesModel struct {
	table     table.Model
	movies    []models.MovieResponse
	Page      int
	Total     int
	PageSize  int
	activeTab MediaType
	perTab    map[MediaType]*tabState
}

type tabState struct {
	movies  []models.MovieResponse
	page    int
	total   int
	loading bool
}

func NewAdminMoviesModel() *AdminMoviesModel {
	cols := []table.Column{
		{Title: "Title", Width: 38},
		{Title: "Year", Width: 6},
		{Title: "Genre", Width: 14},
		{Title: "Format", Width: 8},
		{Title: "Copies", Width: 10},
		{Title: "★", Width: 3},
		{Title: "Meta", Width: 14},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	t.SetStyles(gruvboxTableStyles())

	return &AdminMoviesModel{
		table:     t,
		PageSize:  50,
		activeTab: MediaMovies,
		perTab: map[MediaType]*tabState{
			MediaMovies: {page: 1, loading: true},
			MediaSeries: {page: 1, loading: true},
			MediaGames:  {page: 1, loading: true},
		},
	}
}

func (m *AdminMoviesModel) ActiveTab() MediaType { return m.activeTab }

func (m *AdminMoviesModel) SetActiveTab(t MediaType) {
	if t == m.activeTab {
		return
	}
	m.activeTab = t
	state, ok := m.perTab[t]
	if ok && !state.loading {
		m.applyState(state)
	}
}

func (m *AdminMoviesModel) applyState(s *tabState) {
	m.movies = s.movies
	m.Page = s.page
	m.Total = s.total
	m.refreshRows()
}

func (m *AdminMoviesModel) refreshRows() {
	rows := make([]table.Row, len(m.movies))
	selected := m.table.Cursor()
	for i, mv := range m.movies {
		rows[i] = buildRow(mv, i == selected)
	}
	m.table.SetRows(rows)
}

func (m *AdminMoviesModel) SetMovies(movies []models.MovieResponse, total, page int) {
	if state, ok := m.perTab[m.activeTab]; ok {
		state.movies = movies
		state.page = page
		state.total = total
		state.loading = false
	}
	m.applyState(&tabState{movies: movies, page: page, total: total})
}

func (m *AdminMoviesModel) CurrentPageFor(t MediaType) int {
	if state, ok := m.perTab[t]; ok {
		return state.page
	}
	return 1
}

func (m *AdminMoviesModel) MarkLoading(t MediaType) {
	if state, ok := m.perTab[t]; ok {
		state.loading = true
	}
	if t == m.activeTab {
		m.movies = nil
		m.table.SetRows(nil)
	}
}

func (m *AdminMoviesModel) HasNextPage() bool {
	if m.PageSize <= 0 {
		return false
	}
	return m.Page*m.PageSize < m.Total
}
func (m *AdminMoviesModel) HasPrevPage() bool { return m.Page > 1 }

func (m *AdminMoviesModel) SelectedMovie() *models.MovieResponse {
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.movies) {
		return nil
	}
	mv := m.movies[idx]
	return &mv
}

func (m *AdminMoviesModel) Update(msg tea.Msg) (*AdminMoviesModel, tea.Cmd) {
	prev := m.table.Cursor()
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	if m.table.Cursor() != prev {
		m.refreshRows()
	}
	return m, cmd
}

func (m *AdminMoviesModel) tabBarView(width int) string {
	var cells []string
	for _, t := range AllMediaTypes {
		active := t == m.activeTab
		style := lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
		if active {
			style = style.
				Bold(true).
				Foreground(styles.Green).
				BorderForeground(styles.Green)
		} else {
			style = style.
				Foreground(styles.Grey1).
				BorderForeground(styles.BG5)
		}
		label := t.Label()
		if !active {
			if state, ok := m.perTab[t]; ok && state.total > 0 {
				label = fmt.Sprintf("%s (%d)", t.Label(), state.total)
			}
		}
		cells = append(cells, style.Render(label))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)
	rule := lipgloss.NewStyle().
		Foreground(styles.Green).
		Width(width).
		Render("▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

	return row + "\n" + rule
}

func (m *AdminMoviesModel) View(w, h int) string {
	tabBar := m.tabBarView(w)

	m.table.SetWidth(w)
	tableH := h - 4
	if tableH < 5 {
		tableH = 5
	}
	m.table.SetHeight(tableH)

	state := m.perTab[m.activeTab]
	var status string
	if state.loading {
		status = styles.DimTextStyle.Render("Loading…")
	} else {
		status = styles.DimTextStyle.Render(
			fmt.Sprintf("Page %d  ·  %d total  ·  [tab] switch type  ·  [n/b] page", m.Page, m.Total),
		)
	}

	var body string
	if state.loading {
		body = styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("Loading " + m.activeTab.Label() + "…")
	} else if len(m.movies) == 0 {
		label := m.activeTab.Label()
		short := label
		for i, r := range label {
			if r == ' ' {
				short = label[i+1:]
				break
			}
		}
		body = styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No " + short + " in catalog yet — press [A] to add")
	} else {
		body = m.table.View()
	}

	return strings.Join([]string{tabBar, body, status}, "\n")
}

func buildRow(mv models.MovieResponse, selected bool) table.Row {
	title := mv.Title
	if mv.IsStaffPick {
		title = "★ " + title
	}
	if selected {
		title = "▸ " + title
	}
	format := stripAnsi(components.FormatBadge(mv.Format))
	meta := ""
	switch mv.MediaType {
	case "series":
		if mv.SeasonNumber > 0 {
			meta = fmt.Sprintf("S%02d", mv.SeasonNumber)
		}
		if mv.EpisodeCount > 0 {
			if meta != "" {
				meta += " · "
			}
			meta += fmt.Sprintf("%d eps", mv.EpisodeCount)
		}
	case "game":
		if mv.Platform != "" {
			meta = mv.Platform
		}
	}
	return table.Row{
		truncateStr(title, 38),
		fmt.Sprintf("%d", mv.Year),
		mv.Genre,
		format,
		fmt.Sprintf("%d/%d", mv.CopiesAvailable, mv.CopiesTotal),
		staffPickMark(mv.IsStaffPick),
		meta,
	}
}

func staffPickMark(b bool) string {
	if b {
		return "★"
	}
	return ""
}

func stripAnsi(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		if r == 0x1b {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

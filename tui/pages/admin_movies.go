package pages

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/styles"
)

type AdminMoviesRefreshMsg struct{}

// MediaType is the admin catalog filter. The API returns one of these
// per row in MediaType.
type MediaType string

const (
	MediaMovies MediaType = "movie"
	MediaSeries MediaType = "series"
	MediaGames  MediaType = "game"
)

// AllMediaTypes is the order shown in the tabs.
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

// ─── Item ──────────────────────────────────────────────────────────────────

type adminMovieItem struct {
	movie models.MovieResponse
}

func (a adminMovieItem) Title() string { return a.movie.Title }
func (a adminMovieItem) Description() string {
	return a.detailLine()
}
func (a adminMovieItem) FilterValue() string {
	return a.movie.Title + " " + a.movie.Genre + " " + a.movie.Format + " " + a.movie.Platform
}

func (a adminMovieItem) detailLine() string {
	format := components.FormatBadge(a.movie.Format)
	copies := fmt.Sprintf("%d/%d copies", a.movie.CopiesAvailable, a.movie.CopiesTotal)
	pick := ""
	if a.movie.IsStaffPick {
		pick = "  ★ staff pick"
	}

	// Media-type-specific badges
	meta := ""
	switch a.movie.MediaType {
	case "series":
		if a.movie.SeasonNumber > 0 {
			meta = fmt.Sprintf("S%02d", a.movie.SeasonNumber)
		}
		if a.movie.EpisodeCount > 0 {
			if meta != "" {
				meta += " · "
			}
			meta += fmt.Sprintf("%d eps", a.movie.EpisodeCount)
		}
	case "game":
		if a.movie.Platform != "" {
			meta = a.movie.Platform
		}
	}
	if meta != "" {
		meta = "  ·  " + meta
	}

	return fmt.Sprintf("%d  %s  ·  %s  ·  %s%s%s", a.movie.Year, a.movie.Genre, format, copies, pick, meta)
}

// ─── Delegate ──────────────────────────────────────────────────────────────

type adminMovieDelegate struct{}

func newAdminMovieDelegate() adminMovieDelegate { return adminMovieDelegate{} }

func (d adminMovieDelegate) Height() int                             { return 2 }
func (d adminMovieDelegate) Spacing() int                            { return 2 }
func (d adminMovieDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d adminMovieDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	mi, ok := item.(adminMovieItem)
	if !ok {
		return
	}
	mv := mi.movie

	selected := index == m.Index()

	marker := "  "
	titleStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		titleStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	pick := ""
	if mv.IsStaffPick {
		pick = lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render("  ★")
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		titleStyle.Render(truncateStr(mv.Title, 38)),
		pick,
	)

	// Detail: year, genre, format, copies
	format := components.FormatBadge(mv.Format)

	copyColor := styles.Green
	if mv.CopiesAvailable == 0 {
		copyColor = styles.Red
	} else if mv.CopiesAvailable <= 2 {
		copyColor = styles.Yellow
	}
	copies := lipgloss.NewStyle().Foreground(copyColor).Render(
		fmt.Sprintf("%d/%d copies", mv.CopiesAvailable, mv.CopiesTotal),
	)

	meta := []string{
		fmt.Sprintf("%d", mv.Year),
		mv.Genre,
		format,
		copies,
	}
	switch mv.MediaType {
	case "series":
		if mv.SeasonNumber > 0 {
			meta = append(meta, fmt.Sprintf("S%02d", mv.SeasonNumber))
		}
		if mv.EpisodeCount > 0 {
			meta = append(meta, fmt.Sprintf("%d eps", mv.EpisodeCount))
		}
	case "game":
		if mv.Platform != "" {
			meta = append(meta, mv.Platform)
		}
	}
	metaLine := styles.DimTextStyle.Render("  " + strings.Join(meta, "  ·  "))

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}

// ─── Model ─────────────────────────────────────────────────────────────────

// AdminMoviesModel manages the entire catalog (movies, series, games).
// A top tab row lets the admin switch the active media type; the list
// itself only shows the active type. Each tab maintains its own
// pagination state so switching back doesn't lose your scroll.
type AdminMoviesModel struct {
	list     list.Model
	movies   []models.MovieResponse
	Page     int
	Total    int
	PageSize int

	// Per-tab cached data, so flipping tabs is instant and remembers
	// where the admin left off in each section.
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
	delegate := newAdminMovieDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	// Server-side pagination handles the catalog; the inner list
	// paginator is disabled so 30 fetched items render as a single
	// page (the list will scroll within those 30).
	l.SetShowPagination(false)
	l.Paginator.PerPage = 0
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	return &AdminMoviesModel{
		list:      l,
		PageSize:  50,
		activeTab: MediaMovies,
		perTab: map[MediaType]*tabState{
			MediaMovies: {page: 1, loading: true},
			MediaSeries: {page: 1, loading: true},
			MediaGames:  {page: 1, loading: true},
		},
	}
}

// ActiveTab returns the currently-selected media type.
func (m *AdminMoviesModel) ActiveTab() MediaType { return m.activeTab }

// SetActiveTab switches the visible list to the given media type and
// returns true if the underlying data changed (caller can then refetch).
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
	items := make([]list.Item, len(s.movies))
	for i, mv := range s.movies {
		items[i] = adminMovieItem{movie: mv}
	}
	m.list.SetItems(items)
}

// SetMovies updates the current tab's cached data and refreshes the list.
func (m *AdminMoviesModel) SetMovies(movies []models.MovieResponse, total, page int) {
	if state, ok := m.perTab[m.activeTab]; ok {
		state.movies = movies
		state.page = page
		state.total = total
		state.loading = false
	}
	m.applyState(&tabState{movies: movies, page: page, total: total})
}

// CurrentPageFor returns the page number for the given media type tab.
func (m *AdminMoviesModel) CurrentPageFor(t MediaType) int {
	if state, ok := m.perTab[t]; ok {
		return state.page
	}
	return 1
}

// MarkLoading marks a tab as loading (used when navigating to it).
func (m *AdminMoviesModel) MarkLoading(t MediaType) {
	if state, ok := m.perTab[t]; ok {
		state.loading = true
	}
	if t == m.activeTab {
		m.movies = nil
		m.list.SetItems(nil)
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
	if mi, ok := m.list.SelectedItem().(adminMovieItem); ok {
		return &mi.movie
	}
	return nil
}


func (m *AdminMoviesModel) Update(msg tea.Msg) (*AdminMoviesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// tabBarView renders the three media-type tabs at the top of the
// management screen. The active tab is rendered with a green border
// + bold green text.
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
		// Append the row count for inactive tabs as a faint hint.
		label := t.Label()
		if !active {
			if state, ok := m.perTab[t]; ok && state.total > 0 {
				label = fmt.Sprintf("%s (%d)", t.Label(), state.total)
			}
		}
		cells = append(cells, style.Render(label))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)

	// Underline using a thin green rule that spans the whole width.
	rule := lipgloss.NewStyle().
		Foreground(styles.Green).
		Width(width).
		Render("▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")

	return row + "\n" + rule
}

func (m *AdminMoviesModel) View(w, h int) string {
	// Reserve 4 lines: tabs (2) + status bar (1) + bottom strip (1).
	tabBar := m.tabBarView(w)

	// Reset the list's title for the active tab.
	m.list.Title = m.activeTab.Label()

	listH := h - 4
	if listH < 5 {
		listH = 5
	}
	m.list.SetSize(w, listH)

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
		// Strip the leading emoji + space so the message reads naturally.
		label := m.activeTab.Label()
		// Find first space and drop everything up to and including it.
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
		body = m.list.View()
	}

	return strings.Join([]string{tabBar, body, status}, "\n")
}

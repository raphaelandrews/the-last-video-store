package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type RentRequestMsg struct{ MovieID string }

type MovieDetailModel struct {
	Movie           *models.MovieResponse
	Rental          *models.RentalResponse
	Rented          bool
	ErrMsg          string
	StatusMsg       string
	FreeRentals     int
	MaxFree         int
	Balance         float64
	Choosing        bool
	UseTicket       bool
	SequelTitle     string
	Franchise       []models.MovieResponse
	Recommendations []models.MovieResponse
	RelatedSelected int

	viewport viewport.Model
	ready    bool
	width    int
	height   int
}

func NewMovieDetailModel(m *models.MovieResponse) *MovieDetailModel {
	return &MovieDetailModel{Movie: m, RelatedSelected: -1}
}

func (m *MovieDetailModel) SetRental(r *models.RentalResponse) { m.Rental = r; m.Rented = true }
func (m *MovieDetailModel) SetError(s string)                  { m.ErrMsg = s }
func (m *MovieDetailModel) SetUserContext(freeRentals, maxFree int, balance float64) {
	m.FreeRentals = freeRentals
	m.MaxFree = maxFree
	m.Balance = balance
}

func (m *MovieDetailModel) MoveRelatedUp() {
	all := append(m.Franchise, m.Recommendations...)
	if len(all) == 0 {
		return
	}
	m.RelatedSelected--
	if m.RelatedSelected < -1 {
		m.RelatedSelected = len(all) - 1
	}
}

func (m *MovieDetailModel) MoveRelatedDown() {
	all := append(m.Franchise, m.Recommendations...)
	if len(all) == 0 {
		return
	}
	m.RelatedSelected++
	if m.RelatedSelected >= len(all) {
		m.RelatedSelected = -1
	}
}

func (m *MovieDetailModel) SelectedRelated() *models.MovieResponse {
	all := append(m.Franchise, m.Recommendations...)
	if m.RelatedSelected >= 0 && m.RelatedSelected < len(all) {
		return &all[m.RelatedSelected]
	}
	return nil
}

func (m *MovieDetailModel) SetSize(w, h int) {
	if !m.ready {
		m.width = w
		m.height = h
		return
	}
	m.resize(w, h)
}

func (m *MovieDetailModel) resize(w, h int) {
	headerH := lipgloss.Height(m.headerView(w))
	footerH := lipgloss.Height(m.footerView())
	viewportH := h - headerH - footerH
	if viewportH < 3 {
		viewportH = 3
	}
	m.viewport.Width = w
	m.viewport.Height = viewportH
	m.width = w
	m.height = h
}

func (m *MovieDetailModel) Update(msg tea.Msg) (*MovieDetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, 20)
			m.ready = true
			m.width = msg.Width
			m.height = msg.Height
		}
		m.resize(msg.Width, msg.Height)
		m.viewport.SetContent(m.bodyView(m.width))
	case tea.KeyMsg:
		if !m.ready {
			return m, nil
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *MovieDetailModel) View(w, h int) string {
	if m.Movie == nil {
		return ""
	}
	if !m.ready || m.width != w || m.height != h {
		m.resize(w, h)
		if m.width != w || m.height != h {
			m.width = w
			m.height = h
		}
		m.viewport.SetContent(m.bodyView(w))
	}

	header := m.headerView(w)
	body := m.viewport.View()
	footer := m.footerView()
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *MovieDetailModel) headerView(w int) string {
	if m.Movie == nil {
		return ""
	}
	mv := m.Movie

	title := lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Width(w).
		Align(lipgloss.Left).
		Render("─── " + mv.Title + " ───")

	meta := fmt.Sprintf("%d · %s · %s · Dir: %s", mv.Year, mv.Genre, styles.FormatBadge(mv.Format), mv.Director)
	stars := styles.StarRating(mv.Rating)
	rating := fmt.Sprintf("%s  %.1f/5 (%d ratings)", stars, mv.Rating, mv.RatingCount)

	badge := m.badgeView()

	sequelInfo := ""
	if mv.SequelTo != "" {
		title := m.SequelTitle
		if title == "" {
			title = mv.SequelTo
		}
		sequelInfo = styles.DimTextStyle.Render("📽️ Sequel to: " + title)
	}

	costInfo := m.costInfo()

	lines := []string{title, "", meta, rating, badge}
	if sequelInfo != "" {
		lines = append(lines, sequelInfo)
	}
	if !mv.Available && !m.Rented {
		lines = append(lines, styles.ErrorTextStyle.Render("🔴 No copies available — press [W] to join the waitlist"))
	}
	if costInfo != "" {
		lines = append(lines, styles.TextStyle.Render(costInfo))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *MovieDetailModel) badgeView() string {
	if m.Movie == nil {
		return ""
	}
	mv := m.Movie
	switch {
	case m.Rented:
		out := lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[RENTED ✓]")
		if m.Rental != nil {
			out += "  Due: " + styles.TextStyle.Render(fmt.Sprintf("%d", m.Rental.DueDate))
		}
		return out
	case mv.IsNewRelease:
		return lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render("[NEW RELEASE]")
	case !mv.Available:
		return lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("[RENTED OUT]")
	default:
		return lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[AVAILABLE]")
	}
}

func (m *MovieDetailModel) costInfo() string {
	mv := m.Movie
	if !mv.Available || m.Rented {
		return ""
	}
	if m.Choosing {
		return styles.HighlightStyle.Render("[T] Use ticket  [M] Pay with money  [ESC] Cancel")
	}
	if m.FreeRentals > 0 {
		return fmt.Sprintf("🎟️ Free rental (%d/%d remaining) — Press ENTER to rent", m.FreeRentals, m.MaxFree)
	}
	c := models.MovieCost(mv.RentalPrice, mv.Format)
	return fmt.Sprintf("💵 $%.2f (balance: $%.2f) — Press ENTER to rent", c, m.Balance)
}

func (m *MovieDetailModel) bodyView(w int) string {
	if m.Movie == nil {
		return ""
	}
	mv := m.Movie

	innerW := w - 4
	if innerW < 20 {
		innerW = 20
	}

	divider := lipgloss.NewStyle().Foreground(styles.BG5).Render(strings.Repeat("─", innerW))

	lines := []string{
		divider,
		"",
		styles.HeadingStyle.Render("Synopsis"),
		styles.TextStyle.Width(innerW).Render(mv.Synopsis),
		"",
		divider,
		"",
		fmt.Sprintf("📀 %d/%d copies available", mv.CopiesAvailable, mv.CopiesTotal),
	}

	if len(mv.Cast) > 0 {
		cast := "Cast: "
		for i, c := range mv.Cast {
			if i > 0 {
				cast += ", "
			}
			cast += c
		}
		lines = append(lines, styles.DimTextStyle.Render(cast))
	}

	if len(m.Franchise) > 0 {
		lines = append(lines, "", styles.TextStyle.Bold(true).Render("📽️ Franchise:"))
		for i, f := range m.Franchise {
			prefix := "  "
			if i == m.RelatedSelected {
				prefix = styles.HighlightStyle.Render("▸ ")
			}
			lines = append(lines, styles.DimTextStyle.Render(prefix+f.Title))
		}
	}

	if len(m.Recommendations) > 0 {
		fOffset := len(m.Franchise)
		lines = append(lines, "", styles.TextStyle.Bold(true).Render("Also in "+mv.Genre+":"))
		for i, r := range m.Recommendations {
			prefix := "  • "
			if fOffset+i == m.RelatedSelected {
				prefix = styles.HighlightStyle.Render("▸ ")
			}
			lines = append(lines, styles.DimTextStyle.Render(prefix+r.Title))
		}
	}

	if m.ErrMsg != "" {
		lines = append(lines, "", styles.ErrorTextStyle.Render(m.ErrMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *MovieDetailModel) footerView() string {
	status := ""
	if m.StatusMsg != "" {
		status = styles.SuccessTextStyle.Render("✓ " + m.StatusMsg)
	}

	scrollInfo := ""
	if m.ready {
		pct := m.viewport.ScrollPercent() * 100
		scrollInfo = lipgloss.NewStyle().Foreground(styles.Grey1).Render(
			fmt.Sprintf(" %.0f%% ", pct),
		)
	}

	help := styles.DimTextStyle.Render("↑↓ scroll · [ENTER] rent · [W] waitlist · [Q] back")

	row := help
	if scrollInfo != "" {
		row = lipgloss.JoinVertical(lipgloss.Left, help, scrollInfo)
	}
	if status != "" {
		row = lipgloss.JoinVertical(lipgloss.Left, row, status)
	}
	return row
}

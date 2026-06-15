package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type (
	BrowseSelectMsg struct{ MovieID string }
	BrowseReloadMsg struct{}
)

type BrowseMode int

const (
	ModeAll BrowseMode = iota
	ModeStaffPicks
	ModeLastChance
)

type BrowseModel struct {
	Movies     []models.MovieResponse
	Selected   int
	Total      int
	Loading    bool
	Status     string
	Page       int
	PageSize   int
	TotalPages int
	Mode       BrowseMode
}

func NewBrowseModel() *BrowseModel {
	return &BrowseModel{Selected: -1, Loading: true, Page: 1, PageSize: 40}
}

func (m *BrowseModel) SetMovies(movies []models.MovieResponse, total int, page int) {
	m.Movies = movies
	m.Total = total
	m.Page = page
	if m.Mode == ModeAll {
		m.TotalPages = (total + m.PageSize - 1) / m.PageSize
		m.Status = fmt.Sprintf("page %d/%d · %d movies", page, m.TotalPages, total)
	} else {
		m.TotalPages = 1
		label := "Staff Picks"
		if m.Mode == ModeLastChance {
			label = "Last Chance"
		}
		m.Status = fmt.Sprintf("%s · %d movies", label, total)
	}
	m.Loading = false
	if len(movies) > 0 && m.Selected < 0 {
		m.Selected = 0
	}
	if len(movies) == 0 {
		m.Selected = -1
	}
}

func (m *BrowseModel) HasNextPage() bool { return m.Page < m.TotalPages }
func (m *BrowseModel) HasPrevPage() bool { return m.Page > 1 }

func (m *BrowseModel) MoveUp() {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Movies) - 1
	}
}

func (m *BrowseModel) MoveDown() {
	m.Selected++
	if m.Selected >= len(m.Movies) {
		m.Selected = 0
	}
}

func (m *BrowseModel) SelectedMovie() *models.MovieResponse {
	if m.Selected >= 0 && m.Selected < len(m.Movies) {
		return &m.Movies[m.Selected]
	}
	return nil
}

func (m *BrowseModel) View(w, h int) string {
	if m.Loading {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("Loading catalog..."))
	}
	if len(m.Movies) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("No movies found"))
	}

	cardW := 28
	cols := (w - 4) / (cardW + 2)
	if cols < 1 {
		cols = 1
	}

	var cards []string
	for i, mv := range m.Movies {
		bd := styles.GlassBlue
		if i == m.Selected {
			bd = styles.Yellow
		}
		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(bd).
			Width(cardW).Height(7).Padding(0, 1)

		title := trunc(mv.Title, cardW-6)
		stars := styles.StarRating(mv.Rating)
		fb := styles.FormatBadge(mv.Format)
		status := "[RENT]"
		sc := styles.SuccessGrn
		if !mv.Available {
			status = "[OUT]"
			sc = styles.ErrorRed
		}
		if mv.IsNewRelease {
			status = "[NEW]"
			sc = styles.WarningAmb
		}
		inner := lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(styles.SkyBlue).Bold(true).Render(title),
			styles.DimTextStyle.Render(fmt.Sprintf("(%d)", mv.Year)),
			stars,
			fb+"  "+lipgloss.NewStyle().Foreground(sc).Bold(true).Render(status),
		)
		cards = append(cards, card.Render(inner))
	}

	var rows []string
	for i := 0; i < len(cards); i += cols {
		e := i + cols
		if e > len(cards) {
			e = len(cards)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cards[i:e]...))
	}

	v := lipgloss.JoinVertical(lipgloss.Left, rows...)
	if m.Status != "" {
		v = styles.DimTextStyle.Render(m.Status) + "\n" + v
	}
	return v
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "..."
}

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
	Genre      string
	MediaType  string
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
		label := "movies"
		if m.MediaType == "series" {
			label = "series"
		} else if m.MediaType == "game" {
			label = "games"
		}
		m.Status = fmt.Sprintf("page %d/%d · %d %s", page, m.TotalPages, total, label)
	} else {
		m.TotalPages = 1
		label := "Staff Picks"
		if m.Mode == ModeLastChance {
			label = "Last Chance"
		}
		m.Status = fmt.Sprintf("%s · %d titles", label, total)
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
		label := "movies"
		if m.MediaType == "series" {
			label = "series"
		} else if m.MediaType == "game" {
			label = "games"
		}
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("No "+label+" found"))
	}

	cardW := 28
	cols := (w - 4) / (cardW + 2)
	if cols < 1 {
		cols = 1
	}

	var cards []string
	for i, mv := range m.Movies {
		bd := styles.BG5
		if i == m.Selected {
			bd = styles.Green
		}
		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(bd).
			Width(cardW).Height(7).Padding(0, 1)

		title := trunc(mv.Title, cardW-6)
		stars := styles.StarRating(mv.Rating)
		status := "[RENT]"
		sc := styles.Green
		fb := styles.FormatBadge(mv.Format)
		if !mv.Available {
			status = "[OUT]"
			sc = styles.Red
		} else if mv.MediaType == "game" {
			status = "[PLAY]"
			sc = styles.Orange
			fb = lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render("[" + mv.Platform + "]")
		} else if mv.MediaType == "series" {
			status = fmt.Sprintf("S%d", mv.SeasonNumber)
			if mv.SeasonNumber <= 1 {
				status = "[TV]"
			}
			sc = styles.Purple
		} else if mv.IsNewRelease {
			status = "[NEW]"
			sc = styles.Yellow
		}
		info := fmt.Sprintf("(%d)", mv.Year)
		if mv.EpisodeCount > 0 {
			info += fmt.Sprintf(" · %d eps", mv.EpisodeCount)
		}
		titleColor := styles.FG1
		if i == m.Selected {
			titleColor = styles.Green
		}
		inner := lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(titleColor).Bold(true).Render(title),
			styles.DimTextStyle.Render(info),
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
		statusBar := lipgloss.NewStyle().
			Foreground(styles.Grey1).
			Padding(0, 1).
			Render(m.Status)
		v = statusBar + "\n" + v
	}
	return v
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "..."
}

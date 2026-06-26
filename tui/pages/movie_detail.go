package pages

import (
	"fmt"

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

func (m *MovieDetailModel) View(w, h int) string {
	if m.Movie == nil {
		return ""
	}
	mv := m.Movie

	titleBar := lipgloss.NewStyle().
		Foreground(styles.BG0).
		Background(styles.Green).
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render(" " + mv.Title + " ")

	meta := fmt.Sprintf("%d · %s · %s · Dir: %s", mv.Year, mv.Genre, styles.FormatBadge(mv.Format), mv.Director)
	stars := styles.StarRating(mv.Rating)
	rating := fmt.Sprintf("%s  %.1f/5 (%d ratings)", stars, mv.Rating, mv.RatingCount)

	badge := ""
	switch {
	case m.Rented:
		badge = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[RENTED ✓]")
		if m.Rental != nil {
			badge += "  Due: " + styles.TextStyle.Render(fmt.Sprintf("%d", m.Rental.DueDate))
		}
	case mv.IsNewRelease:
		badge = lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render("[NEW RELEASE]")
	case !mv.Available:
		badge = lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("[RENTED OUT]")
	default:
		badge = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[AVAILABLE]")
	}

	synopsis := styles.TextStyle.Width(w - 4).Render(wrap(mv.Synopsis, w-4))
	copies := fmt.Sprintf("📀 %d/%d copies available", mv.CopiesAvailable, mv.CopiesTotal)

	sequelInfo := ""
	if mv.SequelTo != "" {
		title := m.SequelTitle
		if title == "" {
			title = mv.SequelTo
		}
		sequelInfo = styles.DimTextStyle.Render("📽️ Sequel to: " + title)
	}

	costInfo := ""
	if mv.Available && !m.Rented {
		if m.Choosing {
			costInfo = styles.HighlightStyle.Render("[T] Use ticket  [M] Pay with money  [ESC] Cancel")
		} else if m.FreeRentals > 0 {
			costInfo = fmt.Sprintf("🎟️ Free rental (%d/%d remaining) — Press ENTER to rent", m.FreeRentals, m.MaxFree)
		} else {
			c := models.MovieCost(mv.RentalPrice, mv.Format)
			costInfo = fmt.Sprintf("💵 $%.2f (balance: $%.2f) — Press ENTER to rent", c, m.Balance)
		}
	}
	cast := "Cast: "
	for i, c := range mv.Cast {
		if i > 0 {
			cast += ", "
		}
		cast += c
	}

	divider := lipgloss.NewStyle().Foreground(styles.BG5).Render("────────────────────────────────────────")

	lines := []string{titleBar, "", meta, rating, badge}
	if sequelInfo != "" {
		lines = append(lines, sequelInfo)
	}
	if !mv.Available && !m.Rented {
		lines = append(lines, styles.ErrorTextStyle.Render("🔴 No copies available — press [W] to join the waitlist"))
	}
	if costInfo != "" {
		lines = append(lines, styles.TextStyle.Render(costInfo))
	}
	lines = append(lines, "", divider, "", synopsis, "", copies, styles.DimTextStyle.Render(cast))
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
		lines = append(lines, styles.ErrorTextStyle.Render(m.ErrMsg))
	}
	if m.StatusMsg != "" {
		lines = append(lines, styles.SuccessTextStyle.Render(m.StatusMsg))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func wrap(s string, w int) string {
	if len(s) <= w {
		return s
	}
	var out string
	for len(s) > w {
		brk := w
		for brk > 0 && s[brk] != ' ' {
			brk--
		}
		if brk == 0 {
			brk = w
		}
		out += s[:brk] + "\n"
		s = s[brk:]
		if len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
	return out + s
}

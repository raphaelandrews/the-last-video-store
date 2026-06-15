package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type RentRequestMsg struct{ MovieID string }

type MovieDetailModel struct {
	Movie     *models.MovieResponse
	Rental    *models.RentalResponse
	Rented    bool
	ErrMsg    string
	StatusMsg string
}

func NewMovieDetailModel(m *models.MovieResponse) *MovieDetailModel {
	return &MovieDetailModel{Movie: m}
}

func (m *MovieDetailModel) SetRental(r *models.RentalResponse) { m.Rental = r; m.Rented = true }
func (m *MovieDetailModel) SetError(s string)                  { m.ErrMsg = s }

func (m *MovieDetailModel) View(w, h int) string {
	if m.Movie == nil {
		return ""
	}
	mv := m.Movie

	title := lipgloss.NewStyle().Foreground(styles.SkyBlue).Bold(true).Width(w).Align(lipgloss.Center).Render(mv.Title)
	meta := fmt.Sprintf("%d · %s · %s · Dir: %s", mv.Year, mv.Genre, styles.FormatBadge(mv.Format), mv.Director)
	stars := styles.StarRating(mv.Rating)
	rating := fmt.Sprintf("%s  %.1f/5 (%d ratings)", stars, mv.Rating, mv.RatingCount)

	badge := ""
	switch {
	case m.Rented:
		badge = lipgloss.NewStyle().Foreground(styles.SuccessGrn).Bold(true).Render("[RENTED ✓]")
		if m.Rental != nil {
			badge += "  Due: " + styles.TextStyle.Render(fmt.Sprintf("%d", m.Rental.DueDate))
		}
	case mv.IsNewRelease:
		badge = lipgloss.NewStyle().Foreground(styles.WarningAmb).Bold(true).Render("[NEW RELEASE]")
	case !mv.Available:
		badge = lipgloss.NewStyle().Foreground(styles.ErrorRed).Bold(true).Render("[RENTED OUT]")
	default:
		badge = lipgloss.NewStyle().Foreground(styles.SuccessGrn).Bold(true).Render("[AVAILABLE]")
	}

	synopsis := styles.TextStyle.Width(w - 4).Render(wrap(mv.Synopsis, w-4))
	copies := fmt.Sprintf("📀 %d/%d copies available", mv.CopiesAvailable, mv.CopiesTotal)
	cast := "Cast: "
	for i, c := range mv.Cast {
		if i > 0 {
			cast += ", "
		}
		cast += c
	}

	lines := []string{title, "", meta, rating, badge, "", synopsis, "", copies, styles.DimTextStyle.Render(cast)}
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

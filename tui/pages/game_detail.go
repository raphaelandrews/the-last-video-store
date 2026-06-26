package pages

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type GamePlayStartMsg struct{ GameID string }
type GamePlayEndMsg struct{ SessionID string }

type GameDetailModel struct {
	Game            *models.MovieResponse
	Session         *models.GameSession
	Rental          *models.RentalResponse
	Playing         bool
	ChoosingTime    bool
	ErrMsg          string
	StatusMsg       string
	Balance         float64
	RelatedSelected int
	Recommendations []models.MovieResponse
}

func NewGameDetailModel(m *models.MovieResponse) *GameDetailModel {
	return &GameDetailModel{Game: m, RelatedSelected: -1}
}

func (m *GameDetailModel) SetSession(s *models.GameSession) {
	m.Session = s
	m.Playing = true
	m.ChoosingTime = false
}
func (m *GameDetailModel) SetRental(r *models.RentalResponse) { m.Rental = r; m.StatusMsg = "" }
func (m *GameDetailModel) SetError(s string)                  { m.ErrMsg = s }

func (m *GameDetailModel) CheckExpired() bool {
	if m.Playing && m.Session != nil && m.Session.ExpiresAt > 0 {
		if time.Now().Unix() >= m.Session.ExpiresAt {
			m.Playing = false
			m.Session = nil
			m.ErrMsg = "⌛ Play time expired"
			return true
		}
	}
	return false
}

func (m *GameDetailModel) MoveRelatedUp() {
	if len(m.Recommendations) == 0 {
		return
	}
	m.RelatedSelected--
	if m.RelatedSelected < -1 {
		m.RelatedSelected = len(m.Recommendations) - 1
	}
}

func (m *GameDetailModel) MoveRelatedDown() {
	if len(m.Recommendations) == 0 {
		return
	}
	m.RelatedSelected++
	if m.RelatedSelected >= len(m.Recommendations) {
		m.RelatedSelected = -1
	}
}

func (m *GameDetailModel) SelectedRelated() *models.MovieResponse {
	if m.RelatedSelected >= 0 && m.RelatedSelected < len(m.Recommendations) {
		return &m.Recommendations[m.RelatedSelected]
	}
	return nil
}

func (m *GameDetailModel) View(w, h int) string {
	if m.Game == nil {
		return ""
	}
	g := m.Game

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Orange).
		Bold(true).
		Width(w).
		Align(lipgloss.Center)
	title := titleStyle.Render("─── 🕹️ " + g.Title + " ───")

	meta := fmt.Sprintf("%d · %s · %s · %s", g.Year, g.Genre, g.Platform, g.Director)
	stars := styles.StarRating(g.Rating)
	rating := fmt.Sprintf("%s  %.1f/5 (%d ratings)", stars, g.Rating, g.RatingCount)

	m.CheckExpired()
	badge := lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render("[" + g.Platform + "]")
	if m.Rental != nil {
		badge = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("[RENTED ✓]")
		badge += "  Due: " + styles.TextStyle.Render(fmt.Sprintf("%d", m.Rental.DueDate))
	} else if m.Playing && m.Session != nil {
		remaining := m.Session.ExpiresAt - time.Now().Unix()
		if remaining < 0 {
			remaining = 0
		}
		mins := remaining / 60
		secs := remaining % 60
		badge = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render(fmt.Sprintf("[PLAYING · %dm%02ds]", mins, secs))
	} else if m.ChoosingTime {
		badge = lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render("[SELECT TIME]")
	} else if !g.Available {
		badge = lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("[RENTED OUT]")
	} else if g.PlayPrice > 0 {
		badge += "  " + lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(fmt.Sprintf("$%.2f/hr", g.PlayPrice))
	}

	synopsis := styles.TextStyle.Width(w - 4).Render(wrap(g.Synopsis, w-4))
	copies := fmt.Sprintf("🎮 %d/%d copies available", g.CopiesAvailable, g.CopiesTotal)

	actionLine := ""
	if m.Rental != nil {
		actionLine = "🎮 Game rented — enjoy!"
	} else if m.ChoosingTime {
		actionLine = fmt.Sprintf("⏱️  How long?  [1] 1min · [2] 2min · [3] 3min · [4] 4min · [5] 5min · [ESC] cancel")
	} else if m.Playing {
		actionLine = "🎮 Press [E] to end play session"
	} else if g.Available && g.RentalPrice > 0 {
		actionLine = fmt.Sprintf("🕹️ Play $%.2f/hr — Press [P] (balance: $%.2f)", g.PlayPrice, m.Balance)
	}

	divider := lipgloss.NewStyle().Foreground(styles.BG5).Render("────────────────────────────────────────")

	lines := []string{title, "", meta, rating, badge}
	if actionLine != "" {
		lines = append(lines, styles.TextStyle.Render(actionLine))
	}
	lines = append(lines, "", divider, "", synopsis, "", copies)
	if len(m.Recommendations) > 0 {
		lines = append(lines, "", styles.TextStyle.Bold(true).Render("Also in "+g.Genre+":"))
		for i, r := range m.Recommendations {
			prefix := "  • "
			if i == m.RelatedSelected {
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

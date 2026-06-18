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
	Playing         bool
	ErrMsg          string
	StatusMsg       string
	Balance         float64
	RelatedSelected int
	Recommendations []models.MovieResponse
}

func NewGameDetailModel(m *models.MovieResponse) *GameDetailModel {
	return &GameDetailModel{Game: m, RelatedSelected: -1}
}

func (m *GameDetailModel) SetSession(s *models.GameSession) { m.Session = s; m.Playing = true }
func (m *GameDetailModel) SetError(s string)                { m.ErrMsg = s }

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

	title := lipgloss.NewStyle().Foreground(styles.SkyBlue).Bold(true).Width(w).Align(lipgloss.Center).Render(g.Title)
	meta := fmt.Sprintf("%d · %s · %s · %s", g.Year, g.Genre, g.Platform, g.Director)
	stars := styles.StarRating(g.Rating)
	rating := fmt.Sprintf("%s  %.1f/5 (%d ratings)", stars, g.Rating, g.RatingCount)

	badge := lipgloss.NewStyle().Foreground(styles.Purple).Bold(true).Render("[" + g.Platform + "]")
	if m.Playing {
		elapsed := time.Now().Unix() - m.Session.StartedAt
		mins := elapsed / 60
		badge = lipgloss.NewStyle().Foreground(styles.SuccessGrn).Bold(true).Render(fmt.Sprintf("[PLAYING · %dm]", mins))
	} else if g.PlayPrice > 0 {
		badge += "  " + lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(fmt.Sprintf("$%.2f/hr", g.PlayPrice))
	}

	synopsis := styles.TextStyle.Width(w - 4).Render(wrap(g.Synopsis, w-4))
	copies := fmt.Sprintf("🎮 %d/%d copies available", g.CopiesAvailable, g.CopiesTotal)

	actionLine := ""
	if !m.Playing && g.Available && g.RentalPrice > 0 {
		actionLine = fmt.Sprintf("💵 Rent $%.2f — Press [R]  |  🕹️ Play $%.2f/hr — Press [P] (balance: $%.2f)", g.RentalPrice, g.PlayPrice, m.Balance)
	} else if !m.Playing && g.Available && g.RentalPrice == 0 {
		actionLine = fmt.Sprintf("🕹️ Play $%.2f/hr — Press [P] (balance: $%.2f)", g.PlayPrice, m.Balance)
	} else if m.Playing {
		actionLine = "🎮 Press [E] to end play session"
	}

	lines := []string{title, "", meta, rating, badge}
	if actionLine != "" {
		lines = append(lines, styles.TextStyle.Render(actionLine))
	}
	lines = append(lines, "", synopsis, "", copies)
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

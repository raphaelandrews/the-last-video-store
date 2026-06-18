package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type GameSessionModel struct {
	Sessions []models.GameSession
}

func NewGameSessionModel() *GameSessionModel { return &GameSessionModel{} }

func (m *GameSessionModel) SetSessions(sessions []models.GameSession) { m.Sessions = sessions }

func (m *GameSessionModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🎮 ACTIVE PLAY SESSIONS")

	if len(m.Sessions) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, title, "",
				styles.DimTextStyle.Render("No active play sessions")))
	}

	var rows []string
	for _, s := range m.Sessions {
		elapsed := "just started"
		if s.StartedAt > 0 {
			secs := int64(0)
			if s.EndedAt > 0 {
				secs = s.EndedAt - s.StartedAt
			}
			mins := secs / 60
			elapsed = fmt.Sprintf("%dm", mins)
		}
		line := fmt.Sprintf("  🎮 %-30s by %s · %s",
			s.GameTitle, s.UserID, elapsed)
		rows = append(rows, styles.TextStyle.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title, ""}, rows...)...)
	return content
}

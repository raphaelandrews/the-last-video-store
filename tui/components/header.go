package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type HeaderModel struct{}

func NewHeaderModel() *HeaderModel { return &HeaderModel{} }

func (h *HeaderModel) View(w int, loggedIn bool, username, tier string, points int) string {
	border := lipgloss.NewStyle().Foreground(styles.GlassBlue).Width(w).Render("─")
	title := lipgloss.NewStyle().Foreground(styles.SkyBlue).Background(styles.BgBlue).Bold(true).Width(w).Align(lipgloss.Center).Render("THE LAST VIDEO STORE")
	userLine := ""
	if loggedIn && username != "" {
		badge := styles.TierBadgeStyle(tier).Render(" " + tier + " ")
		userLine = styles.TextStyle.Render("🎫 " + username + "  " + badge + "  🍿 " + fmt.Sprintf("%d", points) + " pts")
	}
	return lipgloss.JoinVertical(lipgloss.Top, border, title, userLine, border)
}

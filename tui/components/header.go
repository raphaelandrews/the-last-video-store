package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type HeaderModel struct{}

func NewHeaderModel() *HeaderModel { return &HeaderModel{} }

func (h *HeaderModel) View(w int, loggedIn bool, username, tier string, points, freeRentals int, balance float64, subscription string) string {
	topBorder := lipgloss.NewStyle().
		Foreground(styles.Orange).
		Background(styles.BG0).
		Width(w).
		Align(lipgloss.Center).
		Render("═══════════════════════════════════════════════════════════════")

	title := lipgloss.NewStyle().
		Foreground(styles.BG0).
		Background(styles.Green).
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render("  🎬 THE LAST VIDEO STORE  ")

	subtitle := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Background(styles.BG0).
		Width(w).
		Align(lipgloss.Center).
		Render("─── Be Kind, Rewind ───")

	userLine := ""
	if loggedIn && username != "" {
		badge := styles.TierBadgeStyle(tier).Render(" " + tier + " ")
		info := fmt.Sprintf(" 🎫 %s  %s  🍿 %d pts ", username, badge, points)
		t := models.TierByName(subscription)
		info += fmt.Sprintf(" 🎟️ %d/%d  🏷️ %s  💵 $%.2f ", freeRentals, t.FreeRentals, t.Label, balance)
		userLine = lipgloss.NewStyle().
			Foreground(styles.FG0).
			Background(styles.BG1).
			Width(w).
			Align(lipgloss.Center).
			Padding(0, 1).
			Render(info)
	}

	bottomBorder := lipgloss.NewStyle().
		Foreground(styles.BG5).
		Background(styles.BG0).
		Width(w).
		Render("───────────────────────────────────────────────────────────────")

	return lipgloss.JoinVertical(lipgloss.Top, topBorder, title, subtitle, userLine, bottomBorder)
}

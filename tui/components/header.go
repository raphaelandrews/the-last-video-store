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
	ornament := lipgloss.NewStyle().
		Foreground(styles.Orange).
		Width(w).
		Align(lipgloss.Center).
		Render("▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")

	title := lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render("  🎬 THE LAST VIDEO STORE  ")

	subtitle := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Italic(true).
		Width(w).
		Align(lipgloss.Center).
		Render("─── Be Kind, Rewind ───")

	userLine := ""
	if loggedIn && username != "" {
		badge := styles.TierBadgeStyle(tier).Render("[" + tier + "]")
		t := models.TierByName(subscription)
		userLine = fmt.Sprintf(" 🎫 %s  %s  🍿 %d pts  🎟️ %d/%d  🏷️ %s  💵 $%.2f ",
			username, badge, points, freeRentals, t.FreeRentals, t.Label, balance)
	}

	bottomBorder := lipgloss.NewStyle().
		Foreground(styles.BG5).
		Width(w).
		Render("───────────────────────────────────────────────────────────────")

	parts := []string{ornament, title, subtitle}
	if userLine != "" {
		parts = append(parts, userLine)
	}
	parts = append(parts, bottomBorder)

	return lipgloss.JoinVertical(lipgloss.Top, parts...)
}

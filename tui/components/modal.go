package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

func ModalView(title, message string, width, height int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Yellow).
		Background(styles.BG1).
		Padding(2, 4).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Background(styles.BG1).
		Bold(true)

	msgStyle := lipgloss.NewStyle().
		Foreground(styles.FG0).
		Background(styles.BG1)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		"",
		msgStyle.Render(message),
		"",
		styles.DimTextStyle.Background(styles.BG1).Render("[ENTER] Confirm  [ESC] Cancel"),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		box.Render(content))
}

func ModalErrorView(title, message string, width, height int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Red).
		Background(styles.BG1).
		Padding(2, 4).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Red).
		Background(styles.BG1).
		Bold(true)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		"",
		styles.TextStyle.Background(styles.BG1).Render(message),
		"",
		styles.DimTextStyle.Background(styles.BG1).Render("[ESC] Dismiss"),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		box.Render(content))
}

func AccessDeniedModal(width, height int) string {
	return ModalErrorView("⛔ ACCESS DENIED", "Insufficient clearance level", width, height)
}

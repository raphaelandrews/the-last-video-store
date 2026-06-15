package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

func ModalView(title, message string, width, height int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Cyan).
		Background(styles.Surface).
		Padding(2, 4).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Background(styles.Surface).
		Bold(true)

	msgStyle := lipgloss.NewStyle().
		Foreground(styles.TextStyle.GetForeground()).
		Background(styles.Surface)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		"",
		msgStyle.Render(message),
		"",
		styles.DimTextStyle.Render("[ENTER] Confirm  [ESC] Cancel"),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		box.Render(content))
}

func ModalErrorView(title, message string, width, height int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Error).
		Background(styles.Surface).
		Padding(2, 4).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Error).
		Background(styles.Surface).
		Bold(true)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		"",
		styles.TextStyle.Render(message),
		"",
		styles.DimTextStyle.Render("[ESC] Dismiss"),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		box.Render(content))
}

func AccessDeniedModal(width, height int) string {
	return ModalErrorView("⛔ ACCESS DENIED", "Insufficient clearance level", width, height)
}

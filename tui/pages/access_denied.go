package pages

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type AccessDeniedModel struct {
	Message string
}

func NewAccessDeniedModel(message string) *AccessDeniedModel {
	return &AccessDeniedModel{Message: message}
}

func (m *AccessDeniedModel) View(w, h int) string {
	banner := lipgloss.NewStyle().
		Foreground(styles.Red).
		Bold(true).
		Width(50).
		Align(lipgloss.Center).
		Render(m.Message)

	content := lipgloss.JoinVertical(lipgloss.Center,
		banner,
		"",
		styles.DimTextStyle.Render("Your current role does not have permission to access this area."),
		"",
		styles.DimTextStyle.Render("Press [Q] to go back"),
	)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

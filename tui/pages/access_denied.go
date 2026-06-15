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
	content := lipgloss.JoinVertical(lipgloss.Center,
		styles.ErrorTextStyle.Bold(true).Render("⛔ ACCESS DENIED"),
		"",
		styles.TextStyle.Render(m.Message),
		"",
		styles.DimTextStyle.Render("Your current role does not have permission to access this area."),
	)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type AdminUsersModel struct {
	users    []models.UserResponse
	selected int
	ErrMsg   string
}

type AdminUsersRefreshMsg struct{}

func NewAdminUsersModel() *AdminUsersModel {
	return &AdminUsersModel{selected: -1}
}

func (m *AdminUsersModel) SetUsers(users []models.UserResponse) {
	m.users = users
}

func (m *AdminUsersModel) SelectedUser() *models.UserResponse {
	if m.selected >= 0 && m.selected < len(m.users) {
		return &m.users[m.selected]
	}
	return nil
}

func (m *AdminUsersModel) MoveUp() {
	if m.selected > 0 {
		m.selected--
	}
}

func (m *AdminUsersModel) MoveDown() {
	if m.selected < len(m.users)-1 {
		m.selected++
	}
}

func (m *AdminUsersModel) View(width, height int) string {
	title := styles.HeadingStyle.Width(width).Align(lipgloss.Center).Render("👥 USER MANAGEMENT")

	header := styles.TextStyle.Bold(true).Render(
		"  USERNAME           TIER          RENTALS   BANNED   TOTP")

	var rows []string
	rows = append(rows, header)

	for i, u := range m.users {
		prefix := "  "
		style := styles.TextStyle
		if i == m.selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			style = styles.HighlightStyle
		}

		banned := "—"
		if u.Banned {
			banned = styles.ErrorTextStyle.Render("BANNED")
		}

		totp := "—"
		if u.TOTPEnabled {
			totp = "🔒"
		}

		badge := styles.TierBadgeStyle(u.TierName).Render(u.TierName)

		line := fmt.Sprintf("%-20s %-14s %-10s %-9s %s",
			prefix+truncateStr(u.Username, 18),
			badge,
			fmt.Sprintf("%d/%d", u.RentalCount, u.MaxRentals),
			banned,
			totp,
		)
		rows = append(rows, style.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	result := lipgloss.JoinVertical(lipgloss.Left, title, content)
	if m.ErrMsg != "" {
		result += "\n" + styles.ErrorTextStyle.Render(m.ErrMsg)
	}
	return result
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

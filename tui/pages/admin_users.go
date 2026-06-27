package pages

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type AdminUsersRefreshMsg struct{}


type adminUserItem struct {
	user models.UserResponse
}

func (a adminUserItem) Title() string { return a.user.Username }
func (a adminUserItem) Description() string {
	return a.detailLine()
}
func (a adminUserItem) FilterValue() string {
	return a.user.Username + " " + a.user.TierName
}

func (a adminUserItem) detailLine() string {
	return fmt.Sprintf("%s  ·  %d/%d rentals", a.user.TierName, a.user.RentalCount, a.user.MaxRentals)
}


type adminUserDelegate struct{}

func newAdminUserDelegate() adminUserDelegate { return adminUserDelegate{} }

func (d adminUserDelegate) Height() int                             { return 2 }
func (d adminUserDelegate) Spacing() int                            { return 1 }
func (d adminUserDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d adminUserDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ui, ok := item.(adminUserItem)
	if !ok {
		return
	}
	u := ui.user

	selected := index == m.Index()

	marker := "  "
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}
	var flags []string
	if u.Banned {
		flags = append(flags, lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("🚫 BANNED"))
	}
	if u.TOTPEnabled {
		flags = append(flags, lipgloss.NewStyle().Foreground(styles.Blue).Bold(true).Render("🔒 TOTP"))
	}
	flagStr := ""
	if len(flags) > 0 {
		flagStr = "  " + strings.Join(flags, "  ")
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		nameStyle.Render(truncateStr(u.Username, 24)),
		flagStr,
	)
	badge := styles.TierBadgeStyle(u.TierName).Render(" " + u.TierName + " ")
	rentals := styles.DimTextStyle.Render(
		fmt.Sprintf("  %d / %d concurrent rentals", u.RentalCount, u.MaxRentals),
	)

	metaLine := lipgloss.JoinHorizontal(lipgloss.Left,
		"  ",
		badge,
		rentals,
	)

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}


type AdminUsersModel struct {
	list      list.Model
	users     []models.UserResponse
	ErrMsg    string
	StatusMsg string
}

func NewAdminUsersModel() *AdminUsersModel {
	delegate := newAdminUserDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "👥 USER MANAGEMENT"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	return &AdminUsersModel{list: l}
}

func (m *AdminUsersModel) SetUsers(users []models.UserResponse) {
	m.users = users
	items := make([]list.Item, len(users))
	for i, u := range users {
		items[i] = adminUserItem{user: u}
	}
	m.list.SetItems(items)
}

func (m *AdminUsersModel) SelectedUser() *models.UserResponse {
	if ui, ok := m.list.SelectedItem().(adminUserItem); ok {
		return &ui.user
	}
	return nil
}


func (m *AdminUsersModel) Update(msg tea.Msg) (*AdminUsersModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *AdminUsersModel) View(w, h int) string {
	m.list.SetSize(w, h-1)
	body := m.list.View()
	if m.ErrMsg != "" {
		body += "\n" + styles.ErrorTextStyle.Render(m.ErrMsg)
	}
	if m.StatusMsg != "" {
		body += "\n" + styles.SuccessTextStyle.Render(m.StatusMsg)
	}
	return body
}

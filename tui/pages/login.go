package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type LoginModel struct {
	form     *huh.Form
	username string
	password string
	errMsg   string
}

func NewLoginModel() *LoginModel {
	m := &LoginModel{}

	theme := gruvboxHuhTheme()

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("username").
				Title("Username").
				Placeholder("e.g. bronze").
				Prompt("▸ ").
				CharLimit(20).
				Validate(func(s string) error {
					if s == "" {
						return errorMsg("username is required")
					}
					return nil
				}).
				Value(&m.username),

			huh.NewInput().
				Key("password").
				Title("Password").
				Placeholder("your secret").
				Prompt("▸ ").
				CharLimit(64).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if s == "" {
						return errorMsg("password is required")
					}
					return nil
				}).
				Value(&m.password),
		),
	).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(theme)

	return m
}

func (m *LoginModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *LoginModel) SetError(s string) { m.errMsg = s }

func (m *LoginModel) Update(msg tea.Msg) (*LoginModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		username := m.username
		password := m.password
		m.username = ""
		m.password = ""
		m.form = NewLoginModel().form
		return m, func() tea.Msg { return LoginRequestMsg{Username: username, Password: password} }
	}

	return m, cmd
}

func (m *LoginModel) View(w, h int) string {
	title := styles.TitleStyle.
		Width(54).
		Align(lipgloss.Center).
		Render("─── MEMBER LOGIN ───")

	subtitle := styles.DimTextStyle.
		Width(54).
		Align(lipgloss.Center).
		Render("Welcome back. Sign in to your account.")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 3).
		Width(54)

	body := m.form.View()

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		subtitle,
		"",
		box.Render(body),
	)

	if m.errMsg != "" {
		content += "\n" + styles.ErrorTextStyle.Render("⛔ "+m.errMsg)
	}

	help := styles.DimTextStyle.
		Width(54).
		Align(lipgloss.Center).
		Render("tab/↑↓ navigate · enter submit · ctrl+r sign up · ctrl+c quit")

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, "", help))
}

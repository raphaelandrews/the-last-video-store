package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type RegisterModel struct {
	form        *huh.Form
	username    string
	password    string
	confirmPass string
	errMsg      string
}

func NewRegisterModel() *RegisterModel {
	m := &RegisterModel{}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("username").
				Title("Choose a username").
				Placeholder("3-20 characters").
				Prompt("▸ ").
				CharLimit(20).
				Validate(func(s string) error {
					if len(s) < 3 {
						return errorMsg("username must be at least 3 characters")
					}
					return nil
				}).
				Value(&m.username),

			huh.NewInput().
				Key("password").
				Title("Create a password").
				Placeholder("at least 6 characters").
				Prompt("▸ ").
				CharLimit(64).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if len(s) < 6 {
						return errorMsg("password must be at least 6 characters")
					}
					return nil
				}).
				Value(&m.password),

			huh.NewInput().
				Key("confirm").
				Title("Confirm your password").
				Placeholder("retype your password").
				Prompt("▸ ").
				CharLimit(64).
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if s != m.password {
						return errorMsg("passwords do not match")
					}
					return nil
				}).
				Value(&m.confirmPass),
		),
	).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(gruvboxHuhTheme())

	return m
}

func (m *RegisterModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *RegisterModel) SetError(s string) { m.errMsg = s }

func (m *RegisterModel) Update(msg tea.Msg) (*RegisterModel, tea.Cmd) {
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
		m.confirmPass = ""
		m.form = NewRegisterModel().form
		return m, func() tea.Msg { return RegisterRequestMsg{Username: username, Password: password} }
	}

	return m, cmd
}

func (m *RegisterModel) View(w, h int) string {
	title := styles.TitleStyle.
		Width(54).
		Align(lipgloss.Center).
		Render("─── NEW MEMBERSHIP ───")

	subtitle := styles.DimTextStyle.
		Width(54).
		Align(lipgloss.Center).
		Render("Create your account to start renting.")

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
		Render("tab/↑↓ navigate · enter submit · ctrl+l back to login · ctrl+c quit")

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, "", help))
}

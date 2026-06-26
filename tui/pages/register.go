package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type RegisterRequestMsg struct {
	Username string
	Password string
}

type RegisterModel struct {
	username        textinput.Model
	password        textinput.Model
	confirmPassword textinput.Model
	focus           int
	errMsg          string
	loading         bool
}

func NewRegisterModel() *RegisterModel {
	un := textinput.New()
	un.Placeholder = "Username (3-20 chars)"
	un.Width = 30
	un.CharLimit = 20
	un.Prompt = "▸ "
	un.Focus()

	pw := textinput.New()
	pw.Placeholder = "Password (min 6 chars)"
	pw.EchoMode = textinput.EchoPassword
	pw.Width = 30
	pw.CharLimit = 64
	pw.Prompt = "▸ "

	cp := textinput.New()
	cp.Placeholder = "Confirm password"
	cp.EchoMode = textinput.EchoPassword
	cp.Width = 30
	cp.CharLimit = 64
	cp.Prompt = "▸ "

	return &RegisterModel{
		username:        un,
		password:        pw,
		confirmPassword: cp,
	}
}

func (m *RegisterModel) Init() tea.Cmd { return textinput.Blink }

func (m *RegisterModel) SetError(s string) { m.errMsg = s; m.loading = false }

func (m *RegisterModel) Update(msg tea.Msg) (*RegisterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.focus = (m.focus + 1) % 3
			m.updateFocus()
			return m, nil
		case tea.KeyEnter:
			u := strings.TrimSpace(m.username.Value())
			p := m.password.Value()
			cp := m.confirmPassword.Value()
			if len(u) < 3 {
				m.errMsg = "Username must be at least 3 characters"
				return m, nil
			}
			if len(p) < 6 {
				m.errMsg = "Password must be at least 6 characters"
				return m, nil
			}
			if p != cp {
				m.errMsg = "Passwords do not match"
				return m, nil
			}
			m.loading = true
			m.errMsg = ""
			return m, func() tea.Msg { return RegisterRequestMsg{Username: u, Password: p} }
		}
	}

	var cmd tea.Cmd
	switch m.focus {
	case 0:
		m.username, cmd = m.username.Update(msg)
	case 1:
		m.password, cmd = m.password.Update(msg)
	case 2:
		m.confirmPassword, cmd = m.confirmPassword.Update(msg)
	}
	return m, cmd
}

func (m *RegisterModel) updateFocus() {
	m.username.Blur()
	m.password.Blur()
	m.confirmPassword.Blur()
	switch m.focus {
	case 0:
		m.username.Focus()
	case 1:
		m.password.Focus()
	case 2:
		m.confirmPassword.Focus()
	}
}

func (m *RegisterModel) View(w, h int) string {
	headerBlock := lipgloss.NewStyle().
		Foreground(styles.BG0).
		Background(styles.Green).
		Bold(true).
		Width(42).
		Align(lipgloss.Center).
		Render("🎬 NEW MEMBERSHIP")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Yellow).
		Background(styles.BG1).
		Padding(2, 4).
		Width(42)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		headerBlock,
		"",
		styles.DimTextStyle.Render("Username:"),
		m.username.View(),
		"",
		styles.DimTextStyle.Render("Password:"),
		m.password.View(),
		"",
		styles.DimTextStyle.Render("Confirm:"),
		m.confirmPassword.View(),
	)

	if m.errMsg != "" {
		inner += "\n" + styles.ErrorTextStyle.Render(m.errMsg)
	}
	if m.loading {
		inner += "\n" + styles.DimTextStyle.Render("Creating account...")
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, box.Render(inner))
}

package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type NavigateMsg struct{ Page string }
type ErrorMsg struct{ Message string }
type LoginSuccessMsg struct {
	AccessToken  string
	RefreshToken string
	User         *models.UserResponse
}
type LoginRequestMsg struct {
	Username string
	Password string
}

type LoginModel struct {
	username textinput.Model
	password textinput.Model
	focus    int
	errMsg   string
	loading  bool
}

func NewLoginModel() *LoginModel {
	un := textinput.New()
	un.Placeholder = "Username"
	un.Width = 30
	un.Prompt = "▸ "
	un.Focus()
	pw := textinput.New()
	pw.Placeholder = "Password"
	pw.EchoMode = textinput.EchoPassword
	pw.Width = 30
	pw.Prompt = "▸ "
	return &LoginModel{username: un, password: pw}
}

func (m *LoginModel) Init() tea.Cmd { return textinput.Blink }

func (m *LoginModel) Update(msg tea.Msg) (*LoginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focus = 1 - m.focus
			m.updateFocus()
			return m, nil
		case "enter":
			u := strings.TrimSpace(m.username.Value())
			p := m.password.Value()
			if u != "" && p != "" {
				m.loading = true
				m.errMsg = ""
				return m, func() tea.Msg { return LoginRequestMsg{Username: u, Password: p} }
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	if m.focus == 0 {
		m.username, cmd = m.username.Update(msg)
	} else {
		m.password, cmd = m.password.Update(msg)
	}
	return m, cmd
}

func (m *LoginModel) SetError(s string) { m.errMsg = s; m.loading = false }

func (m *LoginModel) updateFocus() {
	if m.focus == 0 {
		m.username.Focus()
		m.password.Blur()
	} else {
		m.username.Blur()
		m.password.Focus()
	}
}

func (m *LoginModel) View(w, h int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SkyBlue).
		Background(styles.BgWhite).
		Padding(2, 4).
		Width(42).
		Align(lipgloss.Center)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(styles.SkyBlue).Bold(true).Render("MEMBER LOGIN"),
		"",
		styles.TextStyle.Render("Username:"),
		m.username.View(),
		"",
		styles.TextStyle.Render("Password:"),
		m.password.View(),
	)

	if m.errMsg != "" {
		inner += "\n" + styles.ErrorTextStyle.Render(m.errMsg)
	}
	if m.loading {
		inner += "\n" + styles.DimTextStyle.Render("Authenticating...")
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, box.Render(inner))
}

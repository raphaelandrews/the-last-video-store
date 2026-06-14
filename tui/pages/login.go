package pages

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

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
	un.Focus()
	un.Width = 30
	un.CharLimit = 20
	un.Prompt = "▸ "

	pw := textinput.New()
	pw.Placeholder = "Password"
	pw.EchoMode = textinput.EchoPassword
	pw.Width = 30
	pw.CharLimit = 64
	pw.Prompt = "▸ "

	return &LoginModel{
		username: un,
		password: pw,
		focus:    0,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (*LoginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.focus = (m.focus + 1) % 2
			m.updateFocus()
			return m, nil
		case tea.KeyEnter:
			if m.username.Value() != "" && m.password.Value() != "" {
				m.loading = true
				m.errMsg = ""
				return m, m.loginCmd()
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

func (m *LoginModel) updateFocus() {
	if m.focus == 0 {
		m.username.Focus()
		m.password.Blur()
	} else {
		m.username.Blur()
		m.password.Focus()
	}
}

func (m *LoginModel) loginCmd() tea.Cmd {
	return func() tea.Msg {
		return performLoginMsg{
			username: m.username.Value(),
			password: m.password.Value(),
		}
	}
}

type performLoginMsg struct {
	username, password string
}

func (m *LoginModel) View(width, height int, api interface{}, errMsg string) string {
	title := styles.TitleStyle.Render("╔══════════════════════╗")
	title += "\n" + styles.TitleStyle.Render("║   MEMBER LOGIN       ║")
	title += "\n" + styles.TitleStyle.Render("╚══════════════════════╝")

	form := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		styles.TextStyle.Render("Username:"),
		m.username.View(),
		"",
		styles.TextStyle.Render("Password:"),
		m.password.View(),
	)

	errSection := ""
	if errMsg != "" {
		errSection = "\n\n" + styles.ErrorTextStyle.Render("⚠ "+errMsg)
	}

	loading := ""
	if m.loading {
		loading = "\n\n" + styles.TextStyle.Render("Authenticating...")
	}

	hint := styles.DimTextStyle.Render("\n\n[TAB] switch  [ENTER] login  [R] register  [ESC] quit")
	content := form + errSection + loading + hint

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func (m *LoginModel) SetError(msg string) {
	m.errMsg = msg
	m.loading = false
}

func (m *LoginModel) ClearError() {
	m.errMsg = ""
	m.loading = false
}

func RegisterView(width, height int, api interface{}, errMsg string) string {
	title := styles.TitleStyle.Render("REGISTER NEW MEMBER")
	msg := styles.TextStyle.Render("\nRegistration form coming soon.\nPress ESC to go back.")
	content := lipgloss.JoinVertical(lipgloss.Center, title, msg)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

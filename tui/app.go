package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/pages"
	"github.com/thelastvideostore/tui/styles"
)

type Page int

const (
	PageSplash Page = iota
	PageLogin
	PageRegister
	PageBrowse
	PageMovieDetail
	PageMyRentals
	PageProfile
	PageAdminUsers
	PageAdminMovies
	PageAuditLog
)

type Model struct {
	currentPage Page
	prevPage    Page
	session     *SessionState
	api         *APIClient
	width       int
	height      int
	ready       bool
	lastTick    string

	splash  *pages.SplashModel
	login   *pages.LoginModel
	browse  *pages.BrowseModel
	header  *components.HeaderModel
	message string
	errMsg  string
}

type tickMsg string

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t.Format("Mon Jan 02 2006  3:04 PM"))
	})
}

func NewModel(apiBaseURL string) *Model {
	session := NewSessionState(apiBaseURL)
	api := NewAPIClient(apiBaseURL, session)

	return &Model{
		currentPage: PageSplash,
		session:     session,
		api:         api,
		splash:      pages.NewSplashModel(),
		login:       pages.NewLoginModel(),
		header:      components.NewHeaderModel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.ClearScreen,
		tickCmd(),
		m.splash.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "esc":
			if m.currentPage == PageBrowse {
				return m, nil
			}
			if m.currentPage == PageLogin || m.currentPage == PageSplash {
				return m, tea.Quit
			}
			m.navigateTo(PageBrowse)
			return m, nil
		}

	case tickMsg:
		m.lastTick = string(msg)
		cmds = append(cmds, tickCmd())

	case pages.SplashDoneMsg:
		m.navigateTo(PageLogin)

	case pages.NavigateMsg:
		switch msg.Page {
		case "browse":
			m.navigateTo(PageBrowse)
		case "login":
			m.session.Logout()
			m.navigateTo(PageLogin)
		case "register":
			m.navigateTo(PageRegister)
		case "profile":
			m.navigateTo(PageProfile)
		case "my-rentals":
			m.navigateTo(PageMyRentals)
		case "admin-users":
			m.navigateTo(PageAdminUsers)
		case "admin-movies":
			m.navigateTo(PageAdminMovies)
		case "audit-log":
			m.navigateTo(PageAuditLog)
		}

	case pages.LoginSuccessMsg:
		m.session.Login(msg.AccessToken, msg.RefreshToken, msg.User)
		m.navigateTo(PageBrowse)

	case pages.ErrorMsg:
		m.errMsg = msg.Message

	case pages.ClearErrorMsg:
		m.errMsg = ""
		m.message = ""
	}

	m.header.Update(msg)

	switch m.currentPage {
	case PageSplash:
		_, cmd = m.splash.Update(msg)
		cmds = append(cmds, cmd)
	case PageLogin:
		_, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	headerHeight := 7
	contentHeight := m.height - headerHeight - 3
	if contentHeight < 5 {
		contentHeight = 5
	}

	headerView := m.header.View(m.width,
		m.session.IsLoggedIn,
		func() string {
			if m.session.User != nil {
				return m.session.User.Username
			}
			return ""
		}(),
		func() string {
			if m.session.User != nil {
				return m.session.User.TierName
			}
			return ""
		}(),
		func() int {
			if m.session.User != nil {
				return m.session.User.PopcornPoints
			}
			return 0
		}(),
		m.lastTick,
	)

	var pageView string
	switch m.currentPage {
	case PageSplash:
		pageView = m.splash.View(m.width, contentHeight)
	case PageLogin:
		pageView = m.login.View(m.width, contentHeight, m.api, m.errMsg)
	case PageRegister:
		pageView = pages.RegisterView(m.width, contentHeight, m.api, m.errMsg)
	case PageBrowse:
		if m.browse == nil {
			m.browse = pages.NewBrowseModel()
		}
		contentHeight = m.height - headerHeight - 2
		pageView = m.browse.View(m.width, contentHeight)
	default:
		pageView = centeredText(m.width, contentHeight, "Coming soon...")
	}

	footer := footerView(m.width, m.currentPage, m.session)

	return lipgloss.JoinVertical(lipgloss.Top, headerView, pageView, footer)
}

func (m *Model) navigateTo(page Page) {
	m.prevPage = m.currentPage
	m.currentPage = page
	m.errMsg = ""
	m.message = ""

	switch page {
	case PageBrowse:
		if m.browse == nil {
			m.browse = pages.NewBrowseModel()
		}
	}
}

func footerView(width int, page Page, session *SessionState) string {
	var hints string
	switch page {
	case PageSplash:
		hints = "[ENTER] continue  [ESC] quit"
	case PageLogin:
		hints = "[TAB] switch field  [ENTER] login  [R] register  [ESC] quit"
	case PageRegister:
		hints = "[TAB] switch field  [ENTER] register  [ESC] back"
	case PageBrowse:
		hints = "[↑↓] navigate  [ENTER] details  [/] search  [R] my rentals  [W] wishlist  [P] profile"
		if session.CanAccessAdmin() {
			hints += "  [U] users  [A] audit"
		}
		if session.CanAdminMovies() {
			hints += "  [M] movies"
		}
		hints += "  [Q] quit"
	default:
		hints = "[ESC] back  [Q] quit"
	}

	return styles.FooterStyle.
		Width(width).
		Render(hints)
}

func centeredText(width, height int, text string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		styles.TextStyle.Render(text))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

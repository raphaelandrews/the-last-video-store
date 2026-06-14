package pages

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type SplashDoneMsg struct{}

type NavigateMsg struct {
	Page string
}

type LoginSuccessMsg struct {
	AccessToken  string
	RefreshToken string
	User         *models.UserResponse
}

type ErrorMsg struct {
	Message string
}

type ClearErrorMsg struct{}

type SplashModel struct {
	frame int
	done  bool
}

func NewSplashModel() *SplashModel {
	return &SplashModel{}
}

func (m *SplashModel) Init() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickIncrementMsg{}
	})
}

type tickIncrementMsg struct{}

func (m *SplashModel) Update(msg tea.Msg) (*SplashModel, tea.Cmd) {
	switch msg.(type) {
	case tickIncrementMsg:
		m.frame++
		if m.frame >= 12 {
			m.done = true
			return m, func() tea.Msg { return SplashDoneMsg{} }
		}
		return m, tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
			return tickIncrementMsg{}
		})
	}
	return m, nil
}

func (m *SplashModel) View(width, height int) string {
	if m.done {
		return ""
	}

	logo := []string{
		"██╗  ██╗████████╗██╗     ██╗   ██╗███████╗",
		"╚██╗██╔╝╚══██╔══╝██║     ██║   ██║██╔════╝",
		" ╚███╔╝    ██║   ██║     ██║   ██║███████╗",
		" ██╔██╗    ██║   ██║     ╚██╗ ██╔╝╚════██║",
		"██╔╝ ██╗   ██║   ███████╗ ╚████╔╝ ███████║",
		"╚═╝  ╚═╝   ╚═╝   ╚══════╝  ╚═══╝  ╚══════╝",
	}

	progress := m.frame / 2
	if progress > 5 {
		progress = 5
	}

	var view string
	for i, line := range logo {
		style := lipgloss.NewStyle().
			Foreground(styles.Cyan).
			Background(styles.Background).
			Bold(true)

		if i > progress {
			style = style.Foreground(lipgloss.Color("#0A1A3E"))
		}

		view += style.Render(line) + "\n"
	}

	subtitle := styles.DimTextStyle.Render("THE LAST VIDEO STORE · EST. 2002")
	view += "\n" + subtitle + "\n\n"

	blink := "░"
	if m.frame%4 < 2 {
		blink = "█"
	}
	view += styles.TextStyle.Render("INSERT MEMBERSHIP CARD " + blink)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, view)
}

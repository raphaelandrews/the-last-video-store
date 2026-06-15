package pages

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type SplashDoneMsg struct{}
type tickMsg struct{}

type SplashModel struct {
	frame int
	done  bool
}

func NewSplashModel() *SplashModel { return &SplashModel{} }

func (m *SplashModel) Init() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg{} })
}

func (m *SplashModel) Update(msg tea.Msg) (*SplashModel, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.done = true
		return m, func() tea.Msg { return SplashDoneMsg{} }
	case tickMsg:
		m.frame++
		if m.frame >= 15 {
			m.done = true
			return m, func() tea.Msg { return SplashDoneMsg{} }
		}
		return m, tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg{} })
	}
	return m, nil
}

func (m *SplashModel) View(w, h int) string {
	if m.done {
		return ""
	}
	banner := []string{
		"▀████ █   █▀▄ ██  █ █▀▄ █ ▄▀▀ ▀█▀   █ ▄▀▄ █▀▀ ▀▀█ ▀▀█",
		" █▄▄  ▀▄▀ █ █ █▄▄ █ █▀  █ ▄█▄  █    █ █▀█ █▄▄ ██▄ ██▄",
	}
	prog := m.frame / 3
	if prog > 3 {
		prog = 3
	}
	var lines []string
	for i, l := range banner {
		c := styles.SkyBlue
		if i > prog {
			c = styles.TextLight
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(c).Bold(true).Render(l))
	}
	lines = append(lines, "", styles.DimTextStyle.Render("THE LAST VIDEO STORE"), "")
	blink := "░"
	if m.frame%4 < 2 {
		blink = "█"
	}
	lines = append(lines, styles.TextStyle.Render("INSERT MEMBERSHIP CARD "+blink))
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, lines...))
}

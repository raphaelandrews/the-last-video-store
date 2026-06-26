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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.done = true
			return m, func() tea.Msg { return SplashDoneMsg{} }
		}
	case tickMsg:
		m.frame++
		if m.frame > 24 {
			m.frame = 0
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
		`##::::::::::'###:::::'######::'########:           `,
		` ##:::::::::'## ##:::'##... ##:... ##..::           `,
		` ##::::::::'##:. ##:: ##:::..::::: ##::::           `,
		` ##:::::::'##:::. ##:. ######::::: ##::::           `,
		` ##::::::: #########::..... ##:::: ##::::           `,
		` ##::::::: ##.... ##:'##::: ##:::: ##::::           `,
		` ########: ##:::: ##:. ######::::: ##::::           `,
		`........::..:::::..:::......::::::..:::::           `,
		`'##::::'##:'####:'########::'########::'#######::   `,
		` ##:::: ##:. ##:: ##.... ##: ##.....::'##.... ##:   `,
		` ##:::: ##:: ##:: ##:::: ##: ##::::::: ##:::: ##:   `,
		` ##:::: ##:: ##:: ##:::: ##: ######::: ##:::: ##:   `,
		`. ##:: ##::: ##:: ##:::: ##: ##...:::: ##:::: ##:   `,
		`:. ## ##:::: ##:: ##:::: ##: ##::::::: ##:::: ##:   `,
		`::. ###::::'####: ########:: ########:. #######::   `,
		`:::...:::::....::........:::........:::.......:::   `,
		`:'######::'########::'#######::'########::'########:`,
		`'##... ##:... ##..::'##.... ##: ##.... ##: ##.....: `,
		` ##:::..::::: ##:::: ##:::: ##: ##:::: ##: ##:::::::`,
		`. ######::::: ##:::: ##:::: ##: ########:: ######:::`,
		`:..... ##:::: ##:::: ##:::: ##: ##.. ##::: ##...::::`,
		`'##::: ##:::: ##:::: ##:::: ##: ##::. ##:: ##:::::::`,
		`. ######::::: ##::::. #######:: ##:::. ##: ########:`,
		`:......::::::..::::::.......:::..:::::..::........::`,
	}
	prog := m.frame
	if prog > len(banner) {
		prog = len(banner)
	}
	var lines []string
	for i, l := range banner {
		c := styles.Green
		if i > prog {
			c = styles.Grey0
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(c).Bold(true).Render(l))
	}
	lines = append(lines, "", "")

	tagline := lipgloss.NewStyle().
		Foreground(styles.Grey2).
		Italic(true).
		Render("Friday night. Pick your movie. Grab some snacks. Enjoy the show.")
	lines = append(lines, tagline)
	lines = append(lines, "", "")

	blink := "░"
	if m.frame%4 < 2 {
		blink = "█"
	}
	insertLine := lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Render("  INSERT MEMBERSHIP CARD  " + blink + "  ")
	lines = append(lines, insertLine)

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, lines...))
}

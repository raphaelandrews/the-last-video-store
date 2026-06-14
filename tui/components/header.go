package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type HeaderModel struct {
	nowShowing string
}

func NewHeaderModel() *HeaderModel {
	return &HeaderModel{}
}

func (h *HeaderModel) Init() tea.Cmd {
	return nil
}

func (h *HeaderModel) Update(msg tea.Msg) {}

const banner = `
▀████ █   █▀▄ ██  █ █▀▄ █ ▄▀▀ ▀█▀   █ ▄▀▄ █▀▀ ▀▀█ ▀▀█
 █▄▄  ▀▄▀ █ █ █▄▄ █ █▀  █ ▄█▄  █    █ █▀█ █▄▄ ██▄ ██▄
`

func (h *HeaderModel) View(width int, loggedIn bool, username, tierName string, popcornPoints int, clock string) string {
	bannerStyle := lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Background(styles.Background).
		Bold(true).
		Width(width).
		Align(lipgloss.Center)

	topBorder := lipgloss.NewStyle().
		Foreground(styles.BorderDim).
		Width(width).
		Render("╔" + repeat("═", width-2) + "╗")

	bannerView := bannerStyle.Render(banner)

	infoLine := styles.TextStyle.Width(width).Render(clock)
	if loggedIn && username != "" {
		badge := styles.TierBadgeStyle(tierName).Render(" " + tierName + " ")
		userLine := styles.TextStyle.Render("🎫 " + username + "  " + badge)
		right := styles.TextStyle.Render("🍿 " + fmt.Sprintf("%d", popcornPoints) + " pts")
		spacer := width - lipgloss.Width(userLine) - lipgloss.Width(right) - 4
		if spacer < 1 {
			spacer = 1
		}
		infoLine = lipgloss.JoinHorizontal(lipgloss.Left,
			userLine,
			lipgloss.NewStyle().Width(spacer).Render(""),
			right,
		)
	}

	bottomBorder := lipgloss.NewStyle().
		Foreground(styles.BorderDim).
		Width(width).
		Render("╠" + repeat("═", width-2) + "╣")

	return lipgloss.JoinVertical(lipgloss.Top, topBorder, bannerView, infoLine, bottomBorder)
}

func repeat(s string, n int) string {
	result := ""
	for range n {
		result += s
	}
	return result
}

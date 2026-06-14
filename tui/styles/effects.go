package styles

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Scanlines(width, height int) string {
	var b strings.Builder
	for i := range height {
		if i%3 == 0 {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0F0F38")).
				Background(lipgloss.Color("#0A0A2E")).
				Render(strings.Repeat("░", width)))
		} else {
			b.WriteString(strings.Repeat(" ", width))
		}
		if i < height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func GlitchFrame() string {
	chars := []rune("▓▒░█▄▀▌▐")
	n := rand.Intn(3) + 1
	out := make([]rune, n)
	for i := range out {
		out[i] = chars[rand.Intn(len(chars))]
	}
	return lipgloss.NewStyle().
		Foreground(Magenta).
		Background(Background).
		Render(string(out))
}

var VHSSpinnerFrames = []string{
	"▌", "▌ ", " ▌", " ▌", "▌ ", "▌", " ▌", " ▌ ",
}

func RewindAnimation(tapeName string) string {
	text := "◄◄ REWINDING: " + tapeName + " ... ▌"
	return lipgloss.NewStyle().
		Foreground(Yellow).
		Background(Background).
		Bold(true).
		Render(text)
}

func AccessDeniedOverlay(width, height int) string {
	msg := `
╔══════════════════════════════╗
║  ⛔  ACCESS DENIED           ║
║                              ║
║  Insufficient clearance      ║
║                              ║
║  Press ESC to dismiss        ║
╚══════════════════════════════╝`
	style := lipgloss.NewStyle().
		Foreground(Error).
		Background(Surface).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(Error).
		Padding(1, 2).
		Align(lipgloss.Center)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, style.Render(msg))
}

func LoadingText() string {
	return lipgloss.NewStyle().
		Foreground(NeonGreen).
		Background(Background).
		Bold(true).
		Render("LOADING...")
}

func BeKindRewind() string {
	text := `
╔══════════════════════════════════════════╗
║                                          ║
║   ██▄ ██▀     █▄▀ █ █ █▄▄ █▄▄            ║
║   █▄█ █▄▄     █ █ ▀▄█ █▄█ █▄█            ║
║                                          ║
║      ██▄ ██▀ █ █ █ █ █▄▀ █▄▄             ║
║      █▄█ █▄▄ ▀▄▀ ▀▄▀ █ █ █▄█             ║
║                                          ║
╚══════════════════════════════════════════╝`
	return lipgloss.NewStyle().
		Foreground(Cyan).
		Background(Background).
		Bold(true).
		Render(text)
}

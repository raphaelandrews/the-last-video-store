package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

var spinnerFrames = []string{"▌", "▌ ", " ▌", " ▌", "▌ ", "▌", " ▌", " ▌ "}

func SpinnerView(frame int) string {
	idx := frame % len(spinnerFrames)
	return lipgloss.NewStyle().
		Foreground(styles.Cyan).
		Render(spinnerFrames[idx])
}

func SpinnerWithText(frame int, text string) string {
	return lipgloss.NewStyle().
		Foreground(styles.NeonGreen).
		Render(SpinnerView(frame) + " " + text)
}

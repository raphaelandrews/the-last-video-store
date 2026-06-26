package styles

import "github.com/charmbracelet/lipgloss"

func ModalView(title, msg string, w, h int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(Yellow).
		Background(BG1).
		Padding(2, 4).
		Width(50).
		Align(lipgloss.Center)

	inner := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Foreground(Yellow).Bold(true).Background(BG1).Render(title),
		"",
		TextStyle.Background(BG1).Render(msg),
		"",
		DimTextStyle.Background(BG1).Render("[ENTER] confirm  [ESC] cancel"),
	)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, box.Render(inner))
}

func LoadingView(w, h int) string {
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, TitleStyle.Render("Loading..."))
}

func StarRating(r float64) string {
	s := ""
	for i := range 5 {
		if float64(i) < r-0.5 {
			s += "★"
		} else if float64(i) < r {
			s += "½"
		} else {
			s += "☆"
		}
	}
	return lipgloss.NewStyle().Foreground(Yellow).Render(s)
}

package styles

import "github.com/charmbracelet/lipgloss"

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

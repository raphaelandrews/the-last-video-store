package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type ReturnRequestMsg struct{ RentalID string }
type RentalsReloadMsg struct{}

type MyRentalsModel struct {
	Rentals  []models.RentalResponse
	Selected int
	Status   string
}

func NewMyRentalsModel() *MyRentalsModel { return &MyRentalsModel{Selected: -1} }

func (m *MyRentalsModel) SetRentals(rs []models.RentalResponse) {
	m.Rentals = rs
	if len(rs) > 0 && m.Selected < 0 {
		m.Selected = 0
	}
}

func (m *MyRentalsModel) MoveUp() {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Rentals) - 1
	}
}

func (m *MyRentalsModel) MoveDown() {
	m.Selected++
	if m.Selected >= len(m.Rentals) {
		m.Selected = 0
	}
}

func (m *MyRentalsModel) SelectedRental() *models.RentalResponse {
	if m.Selected >= 0 && m.Selected < len(m.Rentals) {
		return &m.Rentals[m.Selected]
	}
	return nil
}

func (m *MyRentalsModel) View(w, h int) string {
	if len(m.Rentals) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("No rental history"))
	}

	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("MY RENTALS")
	var rows []string
	for i, r := range m.Rentals {
		prefix := "  "
		st := styles.TextStyle
		if i == m.Selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			st = styles.HighlightStyle
		}
		status := "🟢 active"
		if r.Status == "overdue" {
			status = "🔴 overdue"
		} else if r.Status == "returned" {
			status = "✓ returned"
		}
		line := fmt.Sprintf("%s%-30s %s  %s", prefix, truncStr(r.MovieTitle, 28),
			styles.FormatBadge(r.MovieFormat), status)
		fee := r.LateFee + r.RewindFee
		if fee > 0 {
			line += fmt.Sprintf("  $%.2f", fee)
		}
		rows = append(rows, st.Render(line))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, rows...)...)
	if m.Status != "" {
		content += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	return content
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "..."
}

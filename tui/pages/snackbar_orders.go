package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type SnackBarOrdersModel struct {
	Orders []models.SnackBarOrder
}

func NewSnackBarOrdersModel() *SnackBarOrdersModel { return &SnackBarOrdersModel{} }

func (m *SnackBarOrdersModel) SetOrders(orders []models.SnackBarOrder) { m.Orders = orders }

func (m *SnackBarOrdersModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🍿 MY SNACK BAR ORDERS")

	if len(m.Orders) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, title, "",
				styles.DimTextStyle.Render("No snack bar orders yet — press [B] to visit the snack bar!")))
	}

	var rows []string
	var grandTotal float64
	for _, o := range m.Orders {
		line := fmt.Sprintf("  %s %-30s x%d  $%.2f  %s",
			o.Emoji, o.ItemName, o.Quantity, o.Total, styles.DimTextStyle.Render(o.Status))
		rows = append(rows, styles.TextStyle.Render(line))
		grandTotal += o.Total
	}

	rows = append(rows, "",
		styles.TextStyle.Render(fmt.Sprintf("Total spent: $%.2f", grandTotal)))

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title, ""}, rows...)...)
	return content
}

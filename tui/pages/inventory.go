package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type InventoryItem struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	MerchID    string `json:"merch_id"`
	Name       string `json:"name"`
	RedeemedAt int64  `json:"redeemed_at"`
}

type InventoryModel struct {
	Items []InventoryItem
}

func NewInventoryModel() *InventoryModel { return &InventoryModel{} }

func (m *InventoryModel) SetItems(items []InventoryItem) { m.Items = items }

func (m *InventoryModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🎒 MY INVENTORY")

	if len(m.Items) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, title, "",
				styles.DimTextStyle.Render("No collectibles yet — visit the Rewards shop!")))
	}

	var rows []string
	for _, item := range m.Items {
		rows = append(rows, styles.TextStyle.Render(fmt.Sprintf("  • %s", item.Name)))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title, ""}, rows...)...)
	return content
}

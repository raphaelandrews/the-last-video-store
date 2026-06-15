package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type MerchRedeemMsg struct{ ItemID string }
type MerchReloadMsg struct{}

type MerchModel struct {
	Items    []models.MerchItem
	Selected int
	Points   int
	Status   string
	Error    string
}

func NewMerchModel(points int) *MerchModel {
	return &MerchModel{Selected: -1, Points: points}
}

func (m *MerchModel) SetItems(items []models.MerchItem) {
	m.Items = items
	if len(items) > 0 && m.Selected < 0 {
		m.Selected = 0
	}
	if len(items) == 0 {
		m.Selected = -1
	}
}

func (m *MerchModel) MoveUp() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Items) - 1
	}
}

func (m *MerchModel) MoveDown() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected++
	if m.Selected >= len(m.Items) {
		m.Selected = 0
	}
}

func (m *MerchModel) SelectedItem() *models.MerchItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return &m.Items[m.Selected]
	}
	return nil
}

func (m *MerchModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🍿 POPCORN REWARDS")
	balance := fmt.Sprintf("Balance: %d 🍿", m.Points)

	var rows []string
	rows = append(rows, styles.TextStyle.Render(balance), "")

	if len(m.Items) == 0 {
		rows = append(rows, styles.DimTextStyle.Render("Loading rewards catalog..."))
	}

	for i, item := range m.Items {
		prefix := "  "
		st := styles.TextStyle
		if i == m.Selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			st = styles.HighlightStyle
		}

		stock := styles.SuccessTextStyle.Render(fmt.Sprintf("%d in stock", item.Stock))
		if item.Stock <= 0 {
			stock = styles.ErrorTextStyle.Render("out of stock")
		}

		affordable := ""
		if m.Points < item.PointsCost {
			affordable = styles.ErrorTextStyle.Render(" (need " + fmt.Sprintf("%d", item.PointsCost-m.Points) + " more)")
		}

		line := fmt.Sprintf("%s%-30s %4d 🍿  %s%s", prefix, item.Name, item.PointsCost, stock, affordable)
		rows = append(rows, st.Render(line))
		rows = append(rows, styles.DimTextStyle.Render("    "+item.Description))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, rows...)...)
	if m.Status != "" {
		content += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	if m.Error != "" {
		content += "\n" + styles.ErrorTextStyle.Render(m.Error)
	}

	return content
}

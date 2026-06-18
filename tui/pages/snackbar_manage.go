package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type SnackBarManageModel struct {
	Items    []models.SnackBarItem
	Selected int
	Status   string
	Error    string
	IsOwner  bool
}

func NewSnackBarManageModel(tier bitmask.Permission) *SnackBarManageModel {
	return &SnackBarManageModel{
		Selected: -1,
		IsOwner:  bitmask.IsOwner(tier),
	}
}

func (m *SnackBarManageModel) SetItems(items []models.SnackBarItem) {
	m.Items = items
	if len(items) > 0 && m.Selected < 0 {
		m.Selected = 0
	}
	if len(items) == 0 {
		m.Selected = -1
	}
}

func (m *SnackBarManageModel) MoveUp() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Items) - 1
	}
}

func (m *SnackBarManageModel) MoveDown() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected++
	if m.Selected >= len(m.Items) {
		m.Selected = 0
	}
}

func (m *SnackBarManageModel) SelectedItem() *models.SnackBarItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return &m.Items[m.Selected]
	}
	return nil
}

func (m *SnackBarManageModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🍿 SNACK BAR MANAGEMENT")

	var rows []string

	if len(m.Items) == 0 {
		rows = append(rows, styles.DimTextStyle.Render("Loading inventory..."))
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

		line := fmt.Sprintf("%s %s %-28s $%5.2f  %s",
			prefix, item.Emoji, item.Name, item.Price, stock)
		rows = append(rows, st.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title, ""}, rows...)...)
	if m.Status != "" {
		content += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	if m.Error != "" {
		content += "\n" + styles.ErrorTextStyle.Render(m.Error)
	}

	return content
}

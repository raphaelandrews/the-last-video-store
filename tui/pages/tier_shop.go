package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type TierShopModel struct {
	Tiers    []models.TierInfo
	Selected int
	Current  string
	Balance  float64
	Status   string
	Error    string
}

func NewTierShopModel(currentTier string, balance float64) *TierShopModel {
	m := &TierShopModel{
		Tiers:    models.Tiers,
		Selected: -1,
		Current:  currentTier,
		Balance:  balance,
	}
	for i, t := range m.Tiers {
		if t.Name == currentTier {
			m.Selected = i
			break
		}
	}
	if m.Selected < 0 {
		m.Selected = 0
	}
	return m
}

func (m *TierShopModel) MoveUp() {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Tiers) - 1
	}
}

func (m *TierShopModel) MoveDown() {
	m.Selected++
	if m.Selected >= len(m.Tiers) {
		m.Selected = 0
	}
}

func (m *TierShopModel) SelectedTier() *models.TierInfo {
	if m.Selected >= 0 && m.Selected < len(m.Tiers) {
		return &m.Tiers[m.Selected]
	}
	return nil
}

func (m *TierShopModel) View(w, h int) string {
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("🏷️ PREMIUM TIERS")
	balance := fmt.Sprintf("💵 Balance: $%.2f", m.Balance)
	subtitle := styles.DimTextStyle.Render("Each tier grants a monthly free rental allowance and perks")

	var rows []string
	rows = append(rows, styles.TextStyle.Render(balance), subtitle, "")

	for i, t := range m.Tiers {
		prefix := "  "
		st := styles.TextStyle
		if i == m.Selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			st = styles.HighlightStyle
		}

		current := ""
		if t.Name == m.Current {
			current = styles.SuccessTextStyle.Render(" ← current")
		}

		price := "FREE"
		if t.Price > 0 {
			price = fmt.Sprintf("$%.2f", t.Price)
			if m.Balance < t.Price && t.Name != m.Current {
				price = styles.ErrorTextStyle.Render(price + " (insufficient)")
			}
		}

		line := fmt.Sprintf("%s%-12s %10s  %d free/mo  max %d rentals%s",
			prefix, t.Label, price, t.FreeRentals, t.MaxConcurrent, current)
		rows = append(rows, st.Render(line))

		perks := ""
		if t.NewReleasesOK {
			perks += "✓ new releases "
		}
		if t.NoLateFees {
			perks += "✓ no late fees "
		}
		if perks != "" {
			rows = append(rows, styles.DimTextStyle.Render("      "+perks))
		}
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

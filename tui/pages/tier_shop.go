package pages

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

// ─── Item ──────────────────────────────────────────────────────────────────

type tierItem struct {
	tier models.TierInfo
}

func (t tierItem) Title() string       { return t.tier.Label }
func (t tierItem) Description() string { return t.summary() }
func (t tierItem) FilterValue() string { return t.tier.Label + " " + t.tier.Name }

func (t tierItem) summary() string {
	perks := ""
	if t.tier.NewReleasesOK {
		perks += "✓ new releases "
	}
	if t.tier.NoLateFees {
		perks += "✓ no late fees"
	}
	if perks == "" {
		perks = "—"
	}
	return perks
}

// ─── Delegate ──────────────────────────────────────────────────────────────

type tierDelegate struct {
	balance float64
	current string
}

func newTierDelegate(balance float64, current string) tierDelegate {
	return tierDelegate{balance: balance, current: current}
}

func (d tierDelegate) Height() int                             { return 2 }
func (d tierDelegate) Spacing() int                            { return 2 }
func (d tierDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d tierDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ti, ok := item.(tierItem)
	if !ok {
		return
	}
	t := ti.tier

	selected := index == m.Index()

	marker := "  "
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	// Line 1: name + price
	priceStr := "FREE"
	priceColor := styles.Green
	if t.Price > 0 {
		priceColor = styles.Yellow
		priceStr = fmt.Sprintf("$%.2f/mo", t.Price)
		if d.balance < t.Price && t.Name != d.current {
			priceColor = styles.Red
		}
	}
	priceRender := lipgloss.NewStyle().Foreground(priceColor).Bold(true).Render(priceStr)

	currentMarker := ""
	if t.Name == d.current {
		currentMarker = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("  ← current")
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		nameStyle.Render(truncateStr(t.Label, 16)),
		"   ",
		priceRender,
		currentMarker,
	)

	// Line 2: free rentals, max concurrent, perks
	meta := []string{
		fmt.Sprintf("%d free/mo", t.FreeRentals),
		fmt.Sprintf("max %d concurrent", t.MaxConcurrent),
		ti.summary(),
	}
	metaLine := styles.DimTextStyle.Render("  " + strings.Join(meta, "  ·  "))

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}

// ─── Model ─────────────────────────────────────────────────────────────────

type TierShopModel struct {
	list    list.Model
	tiers   []models.TierInfo
	Balance float64
	Current string
	Status  string
	Error   string
}

func NewTierShopModel(currentTier string, balance float64) *TierShopModel {
	tiers := models.Tiers
	delegate := newTierDelegate(balance, currentTier)
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🏷️ PREMIUM TIERS"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	items := make([]list.Item, len(tiers))
	for i, t := range tiers {
		items[i] = tierItem{tier: t}
	}
	l.SetItems(items)

	// Pre-select current tier.
	m := &TierShopModel{
		list:    l,
		tiers:   tiers,
		Balance: balance,
		Current: currentTier,
	}
	for i, t := range tiers {
		if t.Name == currentTier {
			m.list.Select(i)
			break
		}
	}
	return m
}

func (m *TierShopModel) MoveUp()   { m.list.CursorUp() }
func (m *TierShopModel) MoveDown() { m.list.CursorDown() }

func (m *TierShopModel) SelectedTier() *models.TierInfo {
	if ti, ok := m.list.SelectedItem().(tierItem); ok {
		return &ti.tier
	}
	return nil
}

func (m *TierShopModel) Update(msg tea.Msg) (*TierShopModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *TierShopModel) View(w, h int) string {
	m.list.SetSize(w, h-2)
	balanceStr := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(
		fmt.Sprintf("💵 Balance: $%.2f", m.Balance),
	)
	body := balanceStr + "\n" + m.list.View()

	if m.Status != "" {
		body += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	if m.Error != "" {
		body += "\n" + styles.ErrorTextStyle.Render(m.Error)
	}
	return body
}

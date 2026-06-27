package pages

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

// ─── Item ──────────────────────────────────────────────────────────────────

type merchItem struct {
	item models.MerchItem
}

func (m merchItem) Title() string       { return m.item.Name }
func (m merchItem) Description() string { return m.item.Description }
func (m merchItem) FilterValue() string { return m.item.Name + " " + m.item.Description }

// ─── Delegate ──────────────────────────────────────────────────────────────

type merchDelegate struct{}

func newMerchDelegate() merchDelegate { return merchDelegate{} }

func (d merchDelegate) Height() int                             { return 2 }
func (d merchDelegate) Spacing() int                            { return 2 }
func (d merchDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d merchDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	mi, ok := item.(merchItem)
	if !ok {
		return
	}
	it := mi.item

	selected := index == m.Index()

	marker := "  "
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	// Line 1: name + popcorn cost
	costStr := lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render(
		fmt.Sprintf("%d 🍿", it.PointsCost),
	)

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		nameStyle.Render(truncateStr(it.Name, 40)),
		"   ",
		costStr,
	)

	// Line 2: stock + description preview
	var stockGlyph, stockColor lipgloss.Color
	var stockText string
	switch {
	case it.Stock <= 0:
		stockGlyph = "🔴"
		stockColor = styles.Red
		stockText = "out of stock"
	case it.Stock < 3:
		stockGlyph = "🟡"
		stockColor = styles.Yellow
		stockText = fmt.Sprintf("%d left", it.Stock)
	default:
		stockGlyph = "🟢"
		stockColor = styles.Green
		stockText = fmt.Sprintf("%d in stock", it.Stock)
	}
	stockStr := lipgloss.NewStyle().Foreground(stockColor).Render(
		fmt.Sprintf("%s %s", stockGlyph, stockText),
	)

	desc := styles.DimTextStyle.Render("  " + truncateStr(it.Description, 60))

	meta := stockStr + "    " + desc
	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, meta))
}

// ─── Model ─────────────────────────────────────────────────────────────────

type MerchRedeemMsg struct{ ItemID string }
type MerchReloadMsg struct{}

type MerchModel struct {
	list   list.Model
	items  []models.MerchItem
	Points int
	Status string
	Error  string
}

func NewMerchModel(points int) *MerchModel {
	delegate := newMerchDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🍿 POPCORN REWARDS"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	return &MerchModel{list: l, Points: points}
}

func (m *MerchModel) SetItems(items []models.MerchItem) {
	m.items = items
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = merchItem{item: it}
	}
	m.list.SetItems(listItems)
	if len(items) > 0 {
		m.list.Select(0)
	}
}


func (m *MerchModel) SelectedItem() *models.MerchItem {
	if mi, ok := m.list.SelectedItem().(merchItem); ok {
		return &mi.item
	}
	return nil
}

func (m *MerchModel) Update(msg tea.Msg) (*MerchModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *MerchModel) View(w, h int) string {
	if len(m.items) == 0 {
		balanceStr := lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render(
			fmt.Sprintf("🍿 Balance: %d points", m.Points),
		)
		empty := lipgloss.JoinVertical(lipgloss.Center,
			balanceStr,
			"",
			styles.DimTextStyle.Render("Loading rewards catalog..."),
		)
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, empty)
	}

	m.list.SetSize(w, h-2)
	balanceStr := lipgloss.NewStyle().Foreground(styles.Orange).Bold(true).Render(
		fmt.Sprintf("🍿 Balance: %d points", m.Points),
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

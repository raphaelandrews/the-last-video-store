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

type snackItem struct {
	item models.SnackBarItem
}

func (s snackItem) Title() string { return s.item.Name }
func (s snackItem) Description() string {
	return s.item.Description
}
func (s snackItem) FilterValue() string {
	return s.item.Name + " " + s.item.Category + " " + s.item.Description
}

// ─── Delegate ──────────────────────────────────────────────────────────────

type snackDelegate struct{}

func newSnackDelegate() snackDelegate { return snackDelegate{} }

func (d snackDelegate) Height() int                             { return 2 }
func (d snackDelegate) Spacing() int                            { return 2 }
func (d snackDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d snackDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	si, ok := item.(snackItem)
	if !ok {
		return
	}
	it := si.item

	selected := index == m.Index()

	marker := "  "
	emojiStyle := lipgloss.NewStyle().Foreground(styles.FG1)
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		emojiStyle = lipgloss.NewStyle().Foreground(styles.Green)
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	// Line 1: emoji + name + price
	priceColor := styles.Yellow
	if it.Stock <= 0 {
		priceColor = styles.Grey1
	}
	priceStr := lipgloss.NewStyle().Foreground(priceColor).Bold(true).Render(
		fmt.Sprintf("$%.2f", it.Price),
	)

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		emojiStyle.Render(it.Emoji+"  "),
		nameStyle.Render(truncateStr(it.Name, 32)),
		"   ",
		priceStr,
	)

	// Line 2: stock + category + affordability
	var stockGlyph, stockColor lipgloss.Color
	var stockText string
	switch {
	case it.Stock <= 0:
		stockGlyph = "🔴"
		stockColor = styles.Red
		stockText = "out of stock"
	case it.Stock < 5:
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

	meta := []string{
		styles.DimTextStyle.Render("  " + it.Category),
		stockStr,
	}
	metaLine := strings.Join(meta, "  ·  ")

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}

// ─── Model ─────────────────────────────────────────────────────────────────

type SnackBarOrderMsg struct{ ItemID string }

type SnackBarMenuModel struct {
	list     list.Model
	items    []models.SnackBarItem
	selected int
	Balance  float64
	Status   string
	Error    string
}

func NewSnackBarMenuModel(balance float64) *SnackBarMenuModel {
	delegate := newSnackDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🍿 SNACK BAR"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()
	return &SnackBarMenuModel{list: l, Balance: balance}
}

func (m *SnackBarMenuModel) SetItems(items []models.SnackBarItem) {
	m.items = items
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = snackItem{item: it}
	}
	m.list.SetItems(listItems)
	if len(items) > 0 {
		m.list.Select(0)
		m.selected = 0
	}
}

func (m *SnackBarMenuModel) MoveUp()   { m.list.CursorUp() }
func (m *SnackBarMenuModel) MoveDown() { m.list.CursorDown() }

func (m *SnackBarMenuModel) SelectedItem() *models.SnackBarItem {
	if si, ok := m.list.SelectedItem().(snackItem); ok {
		return &si.item
	}
	return nil
}

func (m *SnackBarMenuModel) Update(msg tea.Msg) (*SnackBarMenuModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if i, ok := m.list.SelectedItem().(snackItem); ok {
		for idx, it := range m.items {
			if it.ID == i.item.ID {
				m.selected = idx
				break
			}
		}
	}
	return m, cmd
}

func (m *SnackBarMenuModel) View(w, h int) string {
	if len(m.items) == 0 {
		balanceStr := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(
			fmt.Sprintf("💵 Balance: $%.2f", m.Balance),
		)
		empty := lipgloss.JoinVertical(lipgloss.Center,
			balanceStr,
			"",
			styles.DimTextStyle.Render("Loading snack bar menu..."),
		)
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, empty)
	}

	// Reserve one line for the balance + status at the top.
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

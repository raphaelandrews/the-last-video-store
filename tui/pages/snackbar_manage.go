package pages

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

// ─── Item ──────────────────────────────────────────────────────────────────

type snackManageItem struct {
	item models.SnackBarItem
}

func (s snackManageItem) Title() string { return s.item.Name }
func (s snackManageItem) Description() string {
	return fmt.Sprintf("$%.2f  ·  %d in stock", s.item.Price, s.item.Stock)
}
func (s snackManageItem) FilterValue() string {
	return s.item.Name + " " + s.item.Category
}

// ─── Delegate ──────────────────────────────────────────────────────────────

type snackManageDelegate struct{}

func newSnackManageDelegate() snackManageDelegate { return snackManageDelegate{} }

func (d snackManageDelegate) Height() int                             { return 2 }
func (d snackManageDelegate) Spacing() int                            { return 2 }
func (d snackManageDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d snackManageDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	si, ok := item.(snackManageItem)
	if !ok {
		return
	}
	it := si.item

	selected := index == m.Index()

	marker := "  "
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	priceStr := lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(
		fmt.Sprintf("$%.2f", it.Price),
	)

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		it.Emoji+"  ",
		nameStyle.Render(truncateStr(it.Name, 32)),
		"   ",
		priceStr,
	)

	var stockGlyph, stockColor lipgloss.Color
	var stockText string
	switch {
	case it.Stock <= 0:
		stockGlyph = "🔴"
		stockColor = styles.Red
		stockText = "OUT OF STOCK"
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

	restockHint := styles.DimTextStyle.Render("  ·  press [r] to restock +5")
	meta := stockStr + restockHint

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, meta))
}

// ─── Model ─────────────────────────────────────────────────────────────────

type SnackBarManageModel struct {
	list    list.Model
	items   []models.SnackBarItem
	IsOwner bool
	Status  string
	Error   string
}

func NewSnackBarManageModel(tier bitmask.Permission) *SnackBarManageModel {
	delegate := newSnackManageDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🍿 SNACK BAR MANAGEMENT"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	return &SnackBarManageModel{
		list:    l,
		IsOwner: bitmask.IsOwner(tier),
	}
}

func (m *SnackBarManageModel) SetItems(items []models.SnackBarItem) {
	m.items = items
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = snackManageItem{item: it}
	}
	m.list.SetItems(listItems)
}

func (m *SnackBarManageModel) MoveUp()   { m.list.CursorUp() }
func (m *SnackBarManageModel) MoveDown() { m.list.CursorDown() }

func (m *SnackBarManageModel) SelectedItem() *models.SnackBarItem {
	if si, ok := m.list.SelectedItem().(snackManageItem); ok {
		return &si.item
	}
	return nil
}

func (m *SnackBarManageModel) Update(msg tea.Msg) (*SnackBarManageModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *SnackBarManageModel) View(w, h int) string {
	m.list.SetSize(w, h-1)
	body := m.list.View()
	if m.Status != "" {
		body += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	if m.Error != "" {
		body += "\n" + styles.ErrorTextStyle.Render(m.Error)
	}
	return body
}

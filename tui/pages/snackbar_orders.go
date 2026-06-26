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

type snackOrderItem struct {
	order models.SnackBarOrder
}

func (s snackOrderItem) Title() string { return s.order.ItemName }
func (s snackOrderItem) Description() string {
	return fmt.Sprintf("qty %d  ·  $%.2f  ·  %s", s.order.Quantity, s.order.Total, s.order.Status)
}
func (s snackOrderItem) FilterValue() string { return s.order.ItemName }

// ─── Delegate ──────────────────────────────────────────────────────────────

type snackOrderDelegate struct{}

func newSnackOrderDelegate() snackOrderDelegate { return snackOrderDelegate{} }

func (d snackOrderDelegate) Height() int                             { return 2 }
func (d snackOrderDelegate) Spacing() int                            { return 2 }
func (d snackOrderDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d snackOrderDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	oi, ok := item.(snackOrderItem)
	if !ok {
		return
	}
	o := oi.order

	selected := index == m.Index()

	marker := "  "
	nameStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		o.Emoji+"  ",
		nameStyle.Render(truncateStr(o.ItemName, 32)),
		"   ",
		lipgloss.NewStyle().Foreground(styles.Yellow).Bold(true).Render(
			fmt.Sprintf("$%.2f", o.Total),
		),
	)

	metaLine := styles.DimTextStyle.Render(
		fmt.Sprintf("  qty %d  ·  $%.2f unit  ·  %s", o.Quantity, o.UnitPrice, o.Status),
	)

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}

// ─── Model ─────────────────────────────────────────────────────────────────

type SnackBarOrdersModel struct {
	list   list.Model
	orders []models.SnackBarOrder
}

func NewSnackBarOrdersModel() *SnackBarOrdersModel {
	delegate := newSnackOrderDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🍿 MY SNACK BAR ORDERS"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	return &SnackBarOrdersModel{list: l}
}

func (m *SnackBarOrdersModel) SetOrders(orders []models.SnackBarOrder) {
	m.orders = orders
	items := make([]list.Item, len(orders))
	for i, o := range orders {
		items[i] = snackOrderItem{order: o}
	}
	m.list.SetItems(items)
}

func (m *SnackBarOrdersModel) Update(msg tea.Msg) (*SnackBarOrdersModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *SnackBarOrdersModel) View(w, h int) string {
	if len(m.orders) == 0 {
		empty := styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No snack bar orders yet — press [B] to visit the snack bar!")
		return empty
	}
	m.list.SetSize(w, h-1)
	var total float64
	for _, o := range m.orders {
		total += o.Total
	}
	status := lipgloss.NewStyle().
		Foreground(styles.Yellow).
		Bold(true).
		Render(fmt.Sprintf("Total spent: $%.2f", total))
	return m.list.View() + "\n" + status
}

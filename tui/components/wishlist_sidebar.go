package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type WishlistSidebarModel struct {
	items   []WishlistEntry
	visible bool
}

type WishlistEntry struct {
	MovieID   string
	Title     string
	Format    string
	Available bool
}

func NewWishlistSidebarModel() *WishlistSidebarModel {
	return &WishlistSidebarModel{}
}

func (m *WishlistSidebarModel) SetItems(items []WishlistEntry) {
	m.items = items
}

func (m *WishlistSidebarModel) Toggle() {
	m.visible = !m.visible
}

func (m *WishlistSidebarModel) IsVisible() bool { return m.visible }

func (m *WishlistSidebarModel) View(width int) string {
	if !m.visible {
		return ""
	}

	sidebarW := 26
	if sidebarW > width/3 {
		sidebarW = width / 3
	}

	title := lipgloss.NewStyle().
		Foreground(styles.BG0).
		Background(styles.Green).
		Bold(true).
		Width(sidebarW).
		Align(lipgloss.Center).
		Render("📋 WISHLIST")

	var lines []string
	lines = append(lines, title)

	if len(m.items) == 0 {
		empty := styles.DimTextStyle.
			Width(sidebarW).
			Align(lipgloss.Center).
			Render("Empty — press [W] to add")
		lines = append(lines, empty)
	} else {
		for _, item := range m.items {
			availability := "🔴"
			if item.Available {
				availability = "🟢"
			}
			entry := availability + " " + truncate(item.Title, sidebarW-6) + "\n" +
				"  " + FormatBadge(item.Format)
			lines = append(lines, styles.TextStyle.Width(sidebarW).Render(entry))
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BG5).
		Background(styles.BG1).
		Width(sidebarW).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

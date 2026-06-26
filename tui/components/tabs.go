package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type TabsModel struct {
	tabs      []string
	activeTab int
}

func NewTabsModel(tabs []string) *TabsModel {
	return &TabsModel{tabs: tabs}
}

func (m *TabsModel) SetActive(index int) {
	if index >= 0 && index < len(m.tabs) {
		m.activeTab = index
	}
}

func (m *TabsModel) ActiveIndex() int {
	return m.activeTab
}

func (m *TabsModel) ActiveTab() string {
	if m.activeTab < len(m.tabs) {
		return m.tabs[m.activeTab]
	}
	return ""
}

func (m *TabsModel) Next() {
	m.activeTab = (m.activeTab + 1) % len(m.tabs)
}

func (m *TabsModel) Prev() {
	m.activeTab--
	if m.activeTab < 0 {
		m.activeTab = len(m.tabs) - 1
	}
}

func (m *TabsModel) View(width int) string {
	var items []string
	for i, tab := range m.tabs {
		style := lipgloss.NewStyle().
			Foreground(styles.Grey1).
			Background(styles.BG1).
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(styles.BG5)

		if i == m.activeTab {
			style = style.
				Foreground(styles.BG0).
				Background(styles.Green).
				Bold(true).
				BorderForeground(styles.Green)
		}

		items = append(items, style.Render(tab))
	}

	bar := lipgloss.NewStyle().
		Background(styles.BG1).
		Width(width).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, items...))

	accent := lipgloss.NewStyle().
		Foreground(styles.Green).
		Background(styles.BG0).
		Width(width).
		Render("▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀")

	return lipgloss.JoinVertical(lipgloss.Top, accent, bar)
}

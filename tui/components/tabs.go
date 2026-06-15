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
			Foreground(styles.DimTextStyle.GetForeground()).
			Background(styles.Background).
			Padding(0, 1)

		if i == m.activeTab {
			style = style.
				Foreground(styles.Background).
				Background(styles.Cyan).
				Bold(true)
		}

		items = append(items, style.Render(tab))
	}

	bar := lipgloss.NewStyle().
		Background(styles.Background).
		Width(width).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, items...))

	bottom := lipgloss.NewStyle().
		Foreground(styles.BorderDim).
		Background(styles.Background).
		Width(width).
		Render("┌" + lipgloss.NewStyle().Width(width-2).Render("") + "┐")

	return lipgloss.JoinVertical(lipgloss.Top, bottom, bar)
}

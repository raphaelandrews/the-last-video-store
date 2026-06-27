package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screenKeyMap = help.KeyMap

var helpBinding = key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
)

var quitBinding = key.NewBinding(
	key.WithKeys("ctrl+c"),
	key.WithHelp("ctrl+c", "quit"),
)

var backBinding = key.NewBinding(
	key.WithKeys("q", "esc"),
	key.WithHelp("q/esc", "back"),
)

var pageBinding = key.NewBinding(
	key.WithKeys("left", "right", "h", "l"),
	key.WithHelp("←/→", "page"),
)

type helpKeyModel struct {
	inner screenKeyMap
}

func (h helpKeyModel) ShortHelp() []key.Binding {
	out := []key.Binding{helpBinding}
	if h.inner != nil {
		out = append(out, h.inner.ShortHelp()...)
	}
	return out
}

func (h helpKeyModel) FullHelp() [][]key.Binding {
	cols := [][]key.Binding{}
	if h.inner != nil {
		cols = h.inner.FullHelp()
	}
	cols = append(cols, []key.Binding{helpBinding, quitBinding})
	return cols
}

func helpKeyMatches(msg tea.Msg) bool {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}
	return key.Matches(km, helpBinding)
}

func helpWith(inner screenKeyMap) screenKeyMap {
	return helpKeyModel{inner: inner}
}

func newHelpModel() help.Model {
	h := help.New()
	h.ShortSeparator = " · "
	h.FullSeparator = "    "
	h.Ellipsis = "…"

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#A9B665"))
	grey := lipgloss.NewStyle().Foreground(lipgloss.Color("#928374"))
	orange := lipgloss.NewStyle().Foreground(lipgloss.Color("#E78A4E"))

	h.Styles.ShortKey = green.Bold(true)
	h.Styles.ShortDesc = grey
	h.Styles.ShortSeparator = grey

	h.Styles.FullKey = green.Bold(true)
	h.Styles.FullDesc = grey
	h.Styles.FullSeparator = orange

	h.Styles.Ellipsis = orange

	return h
}

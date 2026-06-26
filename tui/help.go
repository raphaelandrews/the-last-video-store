package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// screenKeyMap aliases the help.KeyMap interface. Every per-screen keymap
// type implements it so it can be rendered by help.Model.
type screenKeyMap = help.KeyMap

// helpBinding is the global `?` toggle (present on every screen except
// splash, login and register which have their own footer).
var helpBinding = key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
)

// quitBinding is the always-on quit key.
var quitBinding = key.NewBinding(
	key.WithKeys("ctrl+c"),
	key.WithHelp("ctrl+c", "quit"),
)

// backBinding is the standard "go back" key shown on detail/sub-screens.
var backBinding = key.NewBinding(
	key.WithKeys("q", "esc"),
	key.WithHelp("q/esc", "back"),
)

// pageBinding is shown on list-based screens that have pagination.
var pageBinding = key.NewBinding(
	key.WithKeys("left", "right", "h", "l"),
	key.WithHelp("←/→", "page"),
)

// helpKeyModel wraps a per-screen keymap so the global `?` toggle and
// quit binding are always available in the help view. It itself satisfies
// help.KeyMap so the composed result can be passed to help.Model.View.
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
	// Always-present meta column: ? toggle + quit.
	cols = append(cols, []key.Binding{helpBinding, quitBinding})
	return cols
}

// helpKeyMatches reports whether the given keymsg matches the global `?`
// help-toggle binding. Returns false if the msg is not a key msg.
func helpKeyMatches(msg tea.Msg) bool {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}
	return key.Matches(km, helpBinding)
}

// helpWith wraps a per-screen keymap so the `?` toggle and quit binding are
// always available in the help view.
func helpWith(inner screenKeyMap) screenKeyMap {
	return helpKeyModel{inner: inner}
}

// newHelpModel returns a fresh help.Model with sensible separators and a
// Gruvbox Material Dark Hard color scheme. The full view appears when the
// user presses `?`; bubbles/help truncates gracefully with an ellipsis if
// the terminal is too narrow.
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

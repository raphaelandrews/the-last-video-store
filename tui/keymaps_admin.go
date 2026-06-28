package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type snackBarMenuKeys struct {
	canManage bool
}

func (s snackBarMenuKeys) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "order")),
		key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "orders")),
	}
	if s.canManage {
		bindings = append(bindings, key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "manage")))
	}
	bindings = append(bindings, pageBinding, backBinding)
	return bindings
}
func (s snackBarMenuKeys) FullHelp() [][]key.Binding {
	cols := [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select item"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "place order"))},
		{key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "order history"))},
	}
	if s.canManage {
		cols = append(cols, []key.Binding{key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "manage inventory"))})
	}
	cols = append(cols, []key.Binding{backBinding})
	return cols
}

type snackBarOrdersKeys struct{}

func (snackBarOrdersKeys) ShortHelp() []key.Binding  { return []key.Binding{backBinding} }
func (snackBarOrdersKeys) FullHelp() [][]key.Binding { return [][]key.Binding{{backBinding}} }

type snackBarManageKeys struct{}

func (snackBarManageKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restock")),
		backBinding,
	}
}
func (snackBarManageKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restock +5"))},
		{backBinding},
	}
}

type myPlaySessionsKeys struct{}

func (myPlaySessionsKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		backBinding,
	}
}
func (myPlaySessionsKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh sessions"))},
		{backBinding},
	}
}

type totpKeys struct{}

func (totpKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		quitBinding,
	}
}
func (totpKeys) FullHelp() [][]key.Binding { return nil }

type accessDeniedKeys struct{}

func (accessDeniedKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
		quitBinding,
	}
}
func (accessDeniedKeys) FullHelp() [][]key.Binding { return nil }

type adminMoviesKeys struct {
	page int
}

func (a adminMoviesKeys) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit")),
		key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "staff pick")),
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch type")),
		key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n/b", "next/prev page")),
		backBinding,
	}
	return bindings
}
func (a adminMoviesKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab/shift+tab", "switch type"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev server page"))},
		{key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit"))},
		{key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "toggle staff pick"))},
		{key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter by title"))},
		{backBinding},
	}
}

type adminUsersKeys struct{}

func (adminUsersKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "promote")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "demote")),
		key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "ban")),
		pageBinding,
		backBinding,
	}
}
func (adminUsersKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "promote")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "demote"))},
		{key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "toggle ban")),
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "toggle TOTP"))},
		{backBinding},
	}
}

type auditLogKeys struct{}

func (auditLogKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "verify")),
		key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n/b", "page")),
		key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "jump to break")),
		backBinding,
	}
}
func (auditLogKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "navigate")),
			key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/dn", "jump page"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev page")),
			key.NewBinding(key.WithKeys("home", "end"), key.WithHelp("home/end", "first/last page"))},
		{key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "verify hash chain")),
			key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "jump to broken entry"))},
		{backBinding},
	}
}

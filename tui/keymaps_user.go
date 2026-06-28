package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type authKeys struct {
	isRegister bool
}

func (a authKeys) ShortHelp() []key.Binding {
	submit := key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	tab := key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch field"))
	switchLabel := "ctrl+r sign up"
	if a.isRegister {
		switchLabel = "ctrl+l to login"
	}
	switcher := key.NewBinding(key.WithKeys("ctrl+r", "ctrl+l"), key.WithHelp(switchLabel, "switch screen"))
	return []key.Binding{submit, tab, switcher, quitBinding}
}
func (a authKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit form"))},
		{key.NewBinding(key.WithKeys("tab", "shift+tab"), key.WithHelp("tab", "next/prev field"))},
		{key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "go to register")),
			key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "go to login"))},
		{quitBinding},
	}
}

type profileKeys struct{}

func (profileKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logout")),
		key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tiers")),
		key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "snack bar")),
		key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "rewards")),
	}
}
func (profileKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logout"))},
		{key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tier shop")),
			key.NewBinding(key.WithKeys("$"), key.WithHelp("$", "top up")),
			key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "toggle TOTP"))},
		{key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "snack bar")),
			key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "rewards shop")),
			key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "inventory"))},
		{backBinding},
	}
}

type wishlistKeys struct{}

func (wishlistKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "remove")),
		backBinding,
	}
}
func (wishlistKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "info"))},
		{key.NewBinding(key.WithKeys("d", "delete"), key.WithHelp("d", "remove"))},
		{backBinding},
	}
}

type inventoryKeys struct{}

func (inventoryKeys) ShortHelp() []key.Binding {
	return []key.Binding{backBinding}
}
func (inventoryKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{backBinding}}
}

type merchKeys struct{}

func (merchKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "redeem")),
		pageBinding,
		backBinding,
	}
}
func (merchKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "redeem"))},
		{backBinding},
	}
}

type tierShopKeys struct{}

func (tierShopKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "purchase")),
		backBinding,
	}
}
func (tierShopKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select tier"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "purchase"))},
		{backBinding},
	}
}

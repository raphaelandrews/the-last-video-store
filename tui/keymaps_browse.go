package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type splashKeys struct{}

func (splashKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
		quitBinding,
	}
}
func (splashKeys) FullHelp() [][]key.Binding { return nil }

type browseKeys struct{}

func (browseKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "profile")),
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rentals")),
	}
}
func (browseKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "navigate"))},
		{key.NewBinding(key.WithKeys("enter", "d"), key.WithHelp("enter", "open details"))},
		{key.NewBinding(key.WithKeys("[", "]"), key.WithHelp("[ ]", "switch media")),
			key.NewBinding(key.WithKeys(",", "."), key.WithHelp(", .", "switch genre"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev page")),
			key.NewBinding(key.WithKeys("s", "l", "a"), key.WithHelp("s l a", "staff/last/all"))},
		{key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
			key.NewBinding(key.WithKeys("f5"), key.WithHelp("f5", "refresh"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "profile")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rentals")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "snack bar")),
			key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "wishlist"))},
	}
}

type movieDetailKeys struct {
	rented bool
}

func (d movieDetailKeys) ShortHelp() []key.Binding {
	if d.rented {
		return []key.Binding{backBinding, quitBinding}
	}
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rent")),
		key.NewBinding(key.WithKeys("↑↓"), key.WithHelp("↑↓", "scroll")),
		key.NewBinding(key.WithKeys("n/p"), key.WithHelp("n/p", "next/prev related")),
		key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "waitlist")),
		backBinding,
	}
}
func (d movieDetailKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "scroll body")),
			key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/dn", "page")),
			key.NewBinding(key.WithKeys("home", "end"), key.WithHelp("home/end", "top/bottom"))},
		{key.NewBinding(key.WithKeys("n", "p"), key.WithHelp("n p", "next/prev related"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rent")),
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "use ticket")),
			key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "pay money"))},
		{key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "waitlist")),
			key.NewBinding(key.WithKeys("f5"), key.WithHelp("f5", "refresh"))},
		{backBinding},
	}
}

type gameDetailKeys struct{}

func (gameDetailKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rent")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "end play")),
		backBinding,
	}
}
func (gameDetailKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rent game"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play in-store"))},
		{key.NewBinding(key.WithKeys("1-5"), key.WithHelp("1-5", "play duration"))},
		{key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "end play session"))},
		{backBinding},
	}
}

type rentalsKeys struct{}

func (rentalsKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "return")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "extend (30🍿)")),
		key.NewBinding(key.WithKeys("←/→"), key.WithHelp("←/→", "page")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play sessions")),
		backBinding,
	}
}
func (rentalsKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "prev/next page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "return rental"))},
		{key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "extend due date"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "view play sessions"))},
		{backBinding},
	}
}

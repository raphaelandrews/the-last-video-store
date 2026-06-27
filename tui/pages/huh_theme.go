package pages

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

func errorMsg(s string) error {
	return errors.New(s)
}

func gruvboxKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Input.Next = key.NewBinding(
		key.WithKeys("enter", "tab", "down"),
		key.WithHelp("enter/tab/↓", "next"),
	)
	km.Input.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "prev"),
	)
	km.Input.Submit = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	)
	return km
}

func gruvboxHuhTheme() *huh.Theme {
	t := huh.ThemeBase()

	greenBold := lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
	greenAccent := lipgloss.NewStyle().Foreground(styles.Green)
	neutral := lipgloss.NewStyle().Foreground(styles.FG0)
	muted := lipgloss.NewStyle().Foreground(styles.Grey1)

	t.Focused.Base = t.Focused.Base.
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(styles.Green)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = greenBold
	t.Focused.Description = muted
	t.Focused.TextInput.Prompt = greenAccent
	t.Focused.TextInput.Cursor = greenAccent
	t.Focused.TextInput.Placeholder = muted
	t.Focused.TextInput.Text = neutral
	t.Focused.NextIndicator = lipgloss.NewStyle().
		SetString("→").
		Foreground(styles.Green).
		MarginLeft(1)
	t.Focused.PrevIndicator = lipgloss.NewStyle().
		SetString("←").
		Foreground(styles.Grey1).
		MarginRight(1)
	t.Focused.SelectSelector = lipgloss.NewStyle().
		SetString("▸ ").
		Foreground(styles.Green)
	t.Focused.Option = neutral
	t.Focused.SelectedOption = greenAccent
	t.Focused.UnselectedOption = neutral
	t.Focused.SelectedPrefix = lipgloss.NewStyle().
		SetString("[•] ").
		Foreground(styles.Green)
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().
		SetString("[ ] ").
		Foreground(styles.Grey1)
	t.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(styles.BG0).
		Background(styles.Green).
		Bold(true)
	t.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(styles.Green).
		Background(styles.BG0)

	t.Focused.ErrorIndicator = lipgloss.NewStyle().
		SetString("⚠").
		Foreground(styles.Red)
	t.Focused.ErrorMessage = lipgloss.NewStyle().
		Foreground(styles.Red)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.
		BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Title = muted
	t.Blurred.TextInput.Text = muted

	return t
}

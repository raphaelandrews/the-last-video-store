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

	t.Focused.Base = t.Focused.Base.
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(styles.Green)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = styles.TitleStyle.Foreground(styles.Green)
	t.Focused.Description = styles.DimTextStyle
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(styles.Green)
	t.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(styles.Green)
	t.Focused.TextInput.Placeholder = styles.DimTextStyle
	t.Focused.TextInput.Text = styles.TextStyle
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
	t.Focused.Option = styles.TextStyle
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(styles.Green)
	t.Focused.UnselectedOption = styles.TextStyle
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
	t.Blurred.Title = styles.DimTextStyle
	t.Blurred.TextInput.Text = styles.DimTextStyle

	return t
}

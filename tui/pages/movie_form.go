package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type MovieFormMode int

const (
	FormAdd MovieFormMode = iota
	FormEdit
)

type MovieFormModel struct {
	mode       MovieFormMode
	movieID    string
	inputs     []textinput.Model
	focused    int
	ErrMsg     string
	justOpened bool
}

type MovieFormSubmitMsg struct {
	Mode     MovieFormMode
	MovieID  string
	Title    string
	Year     int
	Genre    string
	Format   string
	Director string
	Cast     string
	Synopsis string
	Copies   int
	Price    float64
}

func NewMovieFormModel() *MovieFormModel {
	m := &MovieFormModel{mode: FormAdd, justOpened: true}
	m.inputs = make([]textinput.Model, 9)

	labels := []string{"Title", "Year", "Genre", "Format", "Director", "Cast", "Synopsis", "Copies", "Price"}
	for i, label := range labels {
		ti := textinput.New()
		ti.Placeholder = label
		ti.CharLimit = 200
		if i == 7 || i == 8 {
			ti.CharLimit = 10
		}
		m.inputs[i] = ti
	}
	m.inputs[0].Focus()

	return m
}

func NewMovieEditFormModel(mv *models.MovieResponse) *MovieFormModel {
	m := &MovieFormModel{mode: FormEdit, movieID: mv.ID, justOpened: true}
	m.inputs = make([]textinput.Model, 9)

	m.inputs[0] = textinput.New()
	m.inputs[0].SetValue(mv.Title)

	m.inputs[1] = textinput.New()
	m.inputs[1].SetValue(fmt.Sprintf("%d", mv.Year))

	m.inputs[2] = textinput.New()
	m.inputs[2].SetValue(mv.Genre)

	m.inputs[3] = textinput.New()
	m.inputs[3].SetValue(mv.Format)

	m.inputs[4] = textinput.New()
	m.inputs[4].SetValue(mv.Director)

	m.inputs[5] = textinput.New()
	cast := ""
	for i, c := range mv.Cast {
		if i > 0 {
			cast += ", "
		}
		cast += c
	}
	m.inputs[5].SetValue(cast)

	m.inputs[6] = textinput.New()
	m.inputs[6].SetValue(mv.Synopsis)

	m.inputs[7] = textinput.New()
	m.inputs[7].SetValue(fmt.Sprintf("%d", mv.CopiesTotal))

	m.inputs[8] = textinput.New()
	m.inputs[8].SetValue(fmt.Sprintf("%.2f", mv.RentalPrice))

	m.inputs[0].Focus()
	return m
}

func (m *MovieFormModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.justOpened {
			m.justOpened = false
			return nil, nil
		}
		switch msg.String() {
		case "tab":
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % len(m.inputs)
			m.inputs[m.focused].Focus()
			return nil, nil
		case "enter":
			return m.submitCmd(), nil
		}
	}
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return cmd, nil
}

func (m *MovieFormModel) submitCmd() tea.Cmd {
	var year, copies int
	var price float64
	fmt.Sscanf(m.inputs[1].Value(), "%d", &year)
	fmt.Sscanf(m.inputs[7].Value(), "%d", &copies)
	fmt.Sscanf(m.inputs[8].Value(), "%f", &price)

	return func() tea.Msg {
		return MovieFormSubmitMsg{
			Mode:     m.mode,
			MovieID:  m.movieID,
			Title:    m.inputs[0].Value(),
			Year:     year,
			Genre:    m.inputs[2].Value(),
			Format:   m.inputs[3].Value(),
			Director: m.inputs[4].Value(),
			Cast:     m.inputs[5].Value(),
			Synopsis: m.inputs[6].Value(),
			Copies:   copies,
			Price:    price,
		}
	}
}

func (m *MovieFormModel) View(w, h int) string {
	formTitle := "🎬 ADD MOVIE"
	if m.mode == FormEdit {
		formTitle = "🎬 EDIT MOVIE"
	}
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render(formTitle)

	labels := []string{"Title", "Year", "Genre", "Format", "Director", "Cast (comma-separated)", "Synopsis", "Copies", "Price"}
	var rows []string
	for i := range m.inputs {
		label := styles.DimTextStyle.Render(labels[i])
		input := m.inputs[i].View()
		if i == m.focused {
			label = styles.TextStyle.Bold(true).Render(labels[i])
		}
		rows = append(rows, label, input, "")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title, ""}, rows...)...)
	if m.ErrMsg != "" {
		content += "\n" + styles.ErrorTextStyle.Render(m.ErrMsg)
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

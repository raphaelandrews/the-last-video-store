package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type MovieFormMode int

const (
	FormAdd MovieFormMode = iota
	FormEdit
)

const (
	catalogMovie  = "movie"
	catalogSeries = "series"
	catalogGame   = "game"
)

type MovieFormSubmitMsg struct {
	Mode         MovieFormMode
	MovieID      string
	MediaType    string
	Title        string
	Year         int
	Genre        string
	Format       string
	Platform     string
	SeasonNumber int
	EpisodeCount int
	Director     string
	Cast         string
	Synopsis     string
	Copies       int
	Price        float64
}

type MovieFormModel struct {
	form *huh.Form
	mode MovieFormMode

	movieID   string
	mediaType string

	title    string
	year     string
	genre    string
	format   string
	director string
	cast     string
	synopsis string
	copies   string
	price    string
	season   string
	episodes string
	platform string

	genOptions []string
	fmtOptions []string

	errMsg string

	justOpened bool
}

func NewMovieFormModel(tab MediaType) *MovieFormModel {
	mt := string(tab)
	if mt != catalogMovie && mt != catalogSeries && mt != catalogGame {
		mt = catalogMovie
	}
	m := &MovieFormModel{mode: FormAdd, mediaType: mt}
	m.buildForm()
	return m
}

func NewMovieEditFormModel(mv *models.MovieResponse) *MovieFormModel {
	mt := mv.MediaType
	if mt == "" {
		mt = catalogMovie
	}
	m := &MovieFormModel{
		mode:      FormEdit,
		movieID:   mv.ID,
		mediaType: mt,
	}
	m.SeedValues(mv)
	m.buildForm()
	return m
}

func (m *MovieFormModel) SetOptions(genres, formats []string) {
	m.genOptions = genres
	m.fmtOptions = formats
	m.buildForm()
}

func (m *MovieFormModel) Init() tea.Cmd {
	m.justOpened = true
	return m.form.Init()
}

func (m *MovieFormModel) Update(msg tea.Msg) (*MovieFormModel, tea.Cmd) {
	if m.justOpened {
		m.justOpened = false
		if _, isKey := msg.(tea.KeyMsg); isKey {
			return m, nil
		}
	}

	if km, ok := msg.(tea.KeyMsg); ok {
		if km.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		msg := m.ToMessage()
		m.buildForm()
		return m, tea.Batch(m.form.Init(), func() tea.Msg { return msg })
	}

	return m, cmd
}

func (m *MovieFormModel) SetError(s string) {
	m.errMsg = s
}

func (m *MovieFormModel) View(w, h int) string {
	verb := "ADD"
	if m.mode == FormEdit {
		verb = "EDIT"
	}
	typeName := "🎬 Movie"
	if m.mediaType == catalogSeries {
		typeName = "📺 Series"
	} else if m.mediaType == catalogGame {
		typeName = "🕹️ Game"
	}
	formTitle := fmt.Sprintf("─── %s %s ───", verb, typeName)

	header := lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Width(64).
		Align(lipgloss.Center).
		Render(formTitle)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 3).
		Width(64)

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		box.Render(m.form.View()),
	)

	if m.errMsg != "" {
		content += "\n" + styles.ErrorTextStyle.Render("⛔ "+m.errMsg)
	}

	help := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Width(64).
		Align(lipgloss.Center).
		Render("tab/↓ next · shift+tab/↑ prev · enter select/submit · esc cancel")

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, "", help))
}

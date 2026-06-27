package pages

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
		title:     mv.Title,
		year:      itoa(mv.Year),
		genre:     mv.Genre,
		format:    mv.Format,
		director:  mv.Director,
		cast:      joinCast(mv.Cast),
		synopsis:  mv.Synopsis,
		copies:    itoa(mv.CopiesTotal),
		price:     formatFloat(mv.RentalPrice),
		season:    itoa(mv.SeasonNumber),
		episodes:  itoa(mv.EpisodeCount),
		platform:  mv.Platform,
	}
	m.buildForm()
	return m
}

func (m *MovieFormModel) SetOptions(genres, formats []string) {
	m.genOptions = genres
	m.fmtOptions = formats
	m.buildForm()
}

func (m *MovieFormModel) buildForm() {
	genres := m.genOptions
	if len(genres) == 0 {
		genres = []string{"Action", "Comedy", "Drama", "Horror", "SciFi", "Thriller", "Documentary", "Animation"}
	}
	formats := m.fmtOptions
	if len(formats) == 0 {
		formats = []string{"VHS", "DVD", "Blu-ray"}
	}

	titlePH := "The Matrix"
	if m.mediaType == catalogSeries {
		titlePH = "Breaking Bad"
	} else if m.mediaType == catalogGame {
		titlePH = "Halo"
	}

	thisYear := time.Now().Year()

	genreOptions := make([]huh.Option[string], len(genres))
	for i, g := range genres {
		genreOptions[i] = huh.NewOption(g, g)
	}
	formatOptions := make([]huh.Option[string], len(formats))
	for i, f := range formats {
		formatOptions[i] = huh.NewOption(f, f)
	}

	fields := []huh.Field{
		huh.NewInput().
			Key("title").
			Title("Title").
			Placeholder(titlePH).
			Prompt("▸ ").
			CharLimit(200).
			Validate(nonEmptyString("title is required")).
			Value(&m.title),

		huh.NewInput().
			Key("year").
			Title("Year").
			Placeholder("1999").
			Prompt("▸ ").
			CharLimit(4).
			Validate(yearValidator(thisYear)).
			Value(&m.year),

		huh.NewSelect[string]().
			Key("genre").
			Title("Genre").
			Options(genreOptions...).
			Height(minInt(10, len(genreOptions))).
			Validate(func(v string) error {
				if strings.TrimSpace(v) == "" {
					return errorMsg("genre is required")
				}
				return nil
			}).
			Value(&m.genre),

		huh.NewSelect[string]().
			Key("format").
			Title("Format").
			Options(formatOptions...).
			Height(minInt(6, len(formatOptions))).
			Validate(func(v string) error {
				switch v {
				case "VHS", "DVD", "Blu-ray":
					return nil
				}
				return errorMsg("format must be VHS, DVD, or Blu-ray")
			}).
			Value(&m.format),

		huh.NewInput().
			Key("director").
			Title("Director").
			Placeholder("Wachowski").
			Prompt("▸ ").
			CharLimit(120).
			Validate(nonEmptyString("director is required")).
			Value(&m.director),

		huh.NewInput().
			Key("cast").
			Title("Cast (comma-separated)").
			Placeholder("Keanu Reeves, Laurence Fishburne").
			Prompt("▸ ").
			CharLimit(500).
			Validate(nonEmptyString("at least one cast member is required")).
			Value(&m.cast),

		huh.NewInput().
			Key("synopsis").
			Title("Synopsis").
			Placeholder("A computer hacker learns the true nature of reality").
			Prompt("▸ ").
			CharLimit(1000).
			Value(&m.synopsis),

		huh.NewInput().
			Key("copies").
			Title("Total copies").
			Placeholder("5").
			Prompt("▸ ").
			CharLimit(4).
			Validate(positiveIntValidator("copies", 1, 9999)).
			Value(&m.copies),

		huh.NewInput().
			Key("price").
			Title("Rental price (USD)").
			Placeholder("4.99").
			Prompt("▸ ").
			CharLimit(7).
			Validate(priceValidator()).
			Value(&m.price),
	}

	if m.mediaType == catalogSeries {
		fields = append(fields,
			huh.NewInput().
				Key("season").
				Title("Season").
				Placeholder("1").
				Prompt("▸ ").
				CharLimit(2).
				Validate(optionalIntValidator("season", 1, 99)).
				Value(&m.season),

			huh.NewInput().
				Key("episodes").
				Title("Episodes").
				Placeholder("13").
				Prompt("▸ ").
				CharLimit(4).
				Validate(optionalIntValidator("episodes", 1, 9999)).
				Value(&m.episodes),
		)
	}
	if m.mediaType == catalogGame {
		fields = append(fields,
			huh.NewInput().
				Key("platform").
				Title("Platform").
				Placeholder("PS5").
				Prompt("▸ ").
				CharLimit(40).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errorMsg("platform is required for games")
					}
					return nil
				}).
				Value(&m.platform),
		)
	}

	m.form = huh.NewForm(huh.NewGroup(fields...)).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(gruvboxHuhTheme()).
		WithKeyMap(gruvboxKeyMap())
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
		msg := m.buildSubmitMsg()
		m.buildForm()
		return m, tea.Batch(m.form.Init(), func() tea.Msg { return msg })
	}

	return m, cmd
}

func (m *MovieFormModel) buildSubmitMsg() MovieFormSubmitMsg {
	year, _ := strconv.Atoi(strings.TrimSpace(m.year))
	copies, _ := strconv.Atoi(strings.TrimSpace(m.copies))
	price, _ := strconv.ParseFloat(strings.TrimSpace(m.price), 64)

	msg := MovieFormSubmitMsg{
		Mode:      m.mode,
		MovieID:   m.movieID,
		MediaType: m.mediaType,
		Title:     strings.TrimSpace(m.title),
		Year:      year,
		Genre:     strings.TrimSpace(m.genre),
		Format:    strings.TrimSpace(m.format),
		Director:  strings.TrimSpace(m.director),
		Cast:      strings.TrimSpace(m.cast),
		Synopsis:  strings.TrimSpace(m.synopsis),
		Copies:    copies,
		Price:     price,
	}

	if m.mediaType == catalogSeries {
		if s, err := strconv.Atoi(strings.TrimSpace(m.season)); err == nil && s > 0 {
			msg.SeasonNumber = s
		}
		if e, err := strconv.Atoi(strings.TrimSpace(m.episodes)); err == nil && e > 0 {
			msg.EpisodeCount = e
		}
	}
	if m.mediaType == catalogGame {
		msg.Platform = strings.TrimSpace(m.platform)
	}
	return msg
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

func nonEmptyString(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errorMsg(label)
		}
		return nil
	}
}

func yearValidator(thisYear int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errorMsg("year is required")
		}
		if len(s) != 4 {
			return errorMsg("year must be 4 digits (e.g. 1999)")
		}
		for _, r := range s {
			if r < '0' || r > '9' {
				return errorMsg("year must be digits only")
			}
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid year")
		}
		if n < 1880 || n > thisYear+5 {
			return errorMsg(fmt.Sprintf("year must be between 1880 and %d", thisYear+5))
		}
		return nil
	}
}

func positiveIntValidator(label string, min, max int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errorMsg(label + " is required")
		}
		if !digitsOnly(s) {
			return errorMsg(label + " must be digits only")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid number")
		}
		if n < min || n > max {
			return errorMsg(fmt.Sprintf("%s must be between %d and %d", label, min, max))
		}
		return nil
	}
}

func optionalIntValidator(label string, min, max int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		if !digitsOnly(s) {
			return errorMsg(label + " must be digits only")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid number")
		}
		if n < min || n > max {
			return errorMsg(fmt.Sprintf("%s must be between %d and %d", label, min, max))
		}
		return nil
	}
}

func priceValidator() func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		for _, r := range s {
			if (r < '0' || r > '9') && r != '.' {
				return errorMsg("price must be a number (e.g. 4.99)")
			}
		}
		if strings.Count(s, ".") > 1 {
			return errorMsg("price can have at most one decimal point")
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return errorMsg("not a valid amount")
		}
		if n < 0 || n > 999.99 {
			return errorMsg("price must be between 0.00 and 999.99")
		}
		return nil
	}
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func itoa(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

func formatFloat(f float64) string {
	if f == 0 {
		return ""
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func joinCast(cast []string) string {
	out := ""
	for i, c := range cast {
		if i > 0 {
			out += ", "
		}
		out += c
	}
	return out
}

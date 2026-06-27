package pages

import (
	"fmt"
	"strconv"

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

var FormMediaTypes = []struct {
	Label string
	Value string
}{
	{"🎬 Movie", "movie"},
	{"📺 Series", "series"},
	{"🕹️ Game", "game"},
}

type MovieFormModel struct {
	mode      MovieFormMode
	movieID   string
	form      *huh.Form
	mediaType string
	title     string
	year      string
	genre     string
	format    string
	platform  string
	season    string
	episodes  string
	director  string
	cast      string
	synopsis  string
	copies    string
	price     string
	ErrMsg    string
	showType  bool
}

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

func NewMovieFormModel() *MovieFormModel {
	m := &MovieFormModel{mode: FormAdd, mediaType: "movie", showType: true}
	m.form = m.buildForm()
	return m
}

func NewMovieEditFormModel(mv *models.MovieResponse) *MovieFormModel {
	mt := mv.MediaType
	if mt == "" {
		mt = "movie"
	}
	m := &MovieFormModel{
		mode:      FormEdit,
		movieID:   mv.ID,
		mediaType: mt,
		title:     mv.Title,
		year:      fmt.Sprintf("%d", mv.Year),
		genre:     mv.Genre,
		format:    mv.Format,
		platform:  mv.Platform,
		season:    fmt.Sprintf("%d", mv.SeasonNumber),
		episodes:  fmt.Sprintf("%d", mv.EpisodeCount),
		director:  mv.Director,
		cast:      joinCast(mv.Cast),
		synopsis:  mv.Synopsis,
		copies:    fmt.Sprintf("%d", mv.CopiesTotal),
		price:     fmt.Sprintf("%.2f", mv.RentalPrice),
		showType:  true,
	}
	m.form = m.buildForm()
	return m
}

func (m *MovieFormModel) buildForm() *huh.Form {
	yearValidate := func(s string) error {
		y, err := strconv.Atoi(s)
		if err != nil || y < 1880 || y > 2100 {
			return errorMsg("enter a valid year (1880-2100)")
		}
		return nil
	}
	copiesValidate := func(s string) error {
		c, err := strconv.Atoi(s)
		if err != nil || c < 0 {
			return errorMsg("copies must be a non-negative number")
		}
		return nil
	}
	priceValidate := func(s string) error {
		p, err := strconv.ParseFloat(s, 64)
		if err != nil || p < 0 {
			return errorMsg("price must be a non-negative number")
		}
		return nil
	}
	seasonValidate := func(s string) error {
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 {
			return errorMsg("season must be a non-negative number")
		}
		return nil
	}
	episodesValidate := func(s string) error {
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 {
			return errorMsg("episodes must be a non-negative number")
		}
		return nil
	}

	genreOptions := make([]huh.Option[string], 0)
	switch m.mediaType {
	case "series":
		for _, g := range models.SeriesGenreList {
			genreOptions = append(genreOptions, huh.NewOption(g, g))
		}
	case "game":
		for _, g := range models.GameGenreList {
			genreOptions = append(genreOptions, huh.NewOption(g, g))
		}
	default:
		for _, g := range models.GenreList {
			genreOptions = append(genreOptions, huh.NewOption(g, g))
		}
	}

	groups := []*huh.Group{}

	if m.showType {
		mediaTypeOptions := make([]huh.Option[string], len(FormMediaTypes))
		for i, mt := range FormMediaTypes {
			mediaTypeOptions[i] = huh.NewOption(mt.Label, mt.Value)
		}
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Key("mediaType").
				Title("Type").
				Description("movie, series, or game").
				Options(mediaTypeOptions...).
				Value(&m.mediaType),
		))
	}

	common := func(g *huh.Group) *huh.Group { return g }

	groups = append(groups, common(huh.NewGroup(
		huh.NewInput().
			Key("title").
			Title("Title").
			Prompt("▸ ").
			CharLimit(200).
			Validate(func(s string) error {
				if s == "" {
					return errorMsg("title is required")
				}
				return nil
			}).
			Value(&m.title),

		huh.NewInput().
			Key("year").
			Title("Year").
			Prompt("▸ ").
			CharLimit(4).
			Validate(yearValidate).
			Value(&m.year),

		huh.NewSelect[string]().
			Key("genre").
			Title("Genre").
			Options(genreOptions...).
			Value(&m.genre),

		huh.NewSelect[string]().
			Key("format").
			Title("Format").
			Options(
				huh.NewOption("VHS", models.FormatVHS),
				huh.NewOption("DVD", models.FormatDVD),
				huh.NewOption("Blu-ray", models.FormatBluRay),
			).
			Value(&m.format),
	)))

	switch m.mediaType {
	case "series":
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Key("season").
				Title("Season").
				Prompt("▸ ").
				CharLimit(3).
				Validate(seasonValidate).
				Value(&m.season),

			huh.NewInput().
				Key("episodes").
				Title("Episodes").
				Prompt("▸ ").
				CharLimit(4).
				Validate(episodesValidate).
				Value(&m.episodes),
		))
	case "game":
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Key("platform").
				Title("Platform").
				Description("PS5 · Xbox · Switch · PC · Arcade").
				Prompt("▸ ").
				CharLimit(30).
				Value(&m.platform),
		))
	}

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Key("director").
			Title("Director").
			Prompt("▸ ").
			CharLimit(100).
			Value(&m.director),

		huh.NewInput().
			Key("cast").
			Title("Cast").
			Description("comma-separated").
			Prompt("▸ ").
			CharLimit(500).
			Value(&m.cast),

		huh.NewText().
			Key("synopsis").
			Title("Synopsis").
			CharLimit(1000).
			Lines(3).
			Value(&m.synopsis),
	))

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Key("copies").
			Title("Total copies").
			Prompt("▸ ").
			CharLimit(6).
			Validate(copiesValidate).
			Value(&m.copies),

		huh.NewInput().
			Key("price").
			Title("Rental price").
			Description("USD").
			Prompt("▸ ").
			CharLimit(10).
			Validate(priceValidate).
			Value(&m.price),
	))

	return huh.NewForm(groups...).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(gruvboxHuhTheme()).
		WithKeyMap(gruvboxKeyMap())
}

func (m *MovieFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *MovieFormModel) Update(msg tea.Msg) (tea.Cmd, error) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
		return tea.Quit, nil
	}
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		return nil, nil
	}

	prevType := m.mediaType
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.mediaType != prevType {
		preserved := map[string]string{
			"title":    m.title,
			"year":     m.year,
			"director": m.director,
			"cast":     m.cast,
			"synopsis": m.synopsis,
			"copies":   m.copies,
			"price":    m.price,
		}
		m.showType = false
		m.form = m.buildForm()
		m.title = preserved["title"]
		m.year = preserved["year"]
		m.director = preserved["director"]
		m.cast = preserved["cast"]
		m.synopsis = preserved["synopsis"]
		m.copies = preserved["copies"]
		m.price = preserved["price"]
	}

	if m.form.State == huh.StateCompleted {
		year, _ := strconv.Atoi(m.year)
		copies, _ := strconv.Atoi(m.copies)
		price, _ := strconv.ParseFloat(m.price, 64)
		season, _ := strconv.Atoi(m.season)
		episodes, _ := strconv.Atoi(m.episodes)
		mode := m.mode
		movieID := m.movieID
		mediaType := m.mediaType
		title := m.title
		genre := m.genre
		format := m.format
		platform := m.platform
		director := m.director
		cast := m.cast
		synopsis := m.synopsis
		m.title = ""
		m.year = ""
		m.genre = ""
		m.format = ""
		m.platform = ""
		m.season = ""
		m.episodes = ""
		m.director = ""
		m.cast = ""
		m.synopsis = ""
		m.copies = ""
		m.price = ""
		m.form = NewMovieFormModel().form
		return func() tea.Msg {
			return MovieFormSubmitMsg{
				Mode:         mode,
				MovieID:      movieID,
				MediaType:    mediaType,
				Title:        title,
				Year:         year,
				Genre:        genre,
				Format:       format,
				Platform:     platform,
				SeasonNumber: season,
				EpisodeCount: episodes,
				Director:     director,
				Cast:         cast,
				Synopsis:     synopsis,
				Copies:       copies,
				Price:        price,
			}
		}, nil
	}

	return cmd, nil
}

func (m *MovieFormModel) View(w, h int) string {
	formTitle := "─── ADD CATALOG ITEM ───"
	if m.mode == FormEdit {
		formTitle = "─── EDIT CATALOG ITEM ───"
	}

	typeStrip := ""
	if m.showType {
		current := "🎬 Movie"
		switch m.mediaType {
		case "series":
			current = "📺 Series"
		case "game":
			current = "🕹️ Game"
		}
		typeStrip = lipgloss.NewStyle().
			Foreground(styles.Grey1).
			Width(64).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("STEP 1: pick the type · currently: %s", current))
	}

	header := styles.TitleStyle.
		Width(64).
		Align(lipgloss.Center).
		Render(formTitle)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 3).
		Width(64)

	parts := []string{header}
	if typeStrip != "" {
		parts = append(parts, typeStrip)
	}
	parts = append(parts, "", box.Render(m.form.View()))

	content := lipgloss.JoinVertical(lipgloss.Center, parts...)

	if m.ErrMsg != "" {
		content += "\n" + styles.ErrorTextStyle.Render("⛔ "+m.ErrMsg)
	}

	help := styles.DimTextStyle.
		Width(64).
		Align(lipgloss.Center).
		Render("tab/↓ next · shift+tab/↑ prev · enter submit · esc cancel")

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, content, "", help))
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

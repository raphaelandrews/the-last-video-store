package pages

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

const (
	catalogMovie  = "movie"
	catalogSeries = "series"
	catalogGame   = "game"
)

type catalogField struct {
	label       string
	placeholder string
	value       string
	required    bool
	validate    func(string) error
}

type MovieFormModel struct {
	mode    MovieFormMode
	movieID string

	mediaType string

	fields []catalogField
	focus  int
	inputs []textinput.Model
	width  int

	ErrMsg string
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
	m := &MovieFormModel{mode: FormAdd, mediaType: catalogMovie}
	m.buildInputs()
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
	m.buildInputs()
	for i := range m.inputs {
		if i < len(m.fields) {
			m.inputs[i].SetValue(m.fields[i].value)
		}
	}
	return m
}

func (m *MovieFormModel) buildInputs() {
	m.fields = nil
	m.inputs = nil

	add := func(label, placeholder, def string, required bool, validate func(string) error) {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.CharLimit = 200
		ti.Width = 30
		ti.Prompt = "▸ "
		ti.SetValue(def)
		m.inputs = append(m.inputs, ti)
		m.fields = append(m.fields, catalogField{
			label:       label,
			placeholder: placeholder,
			value:       def,
			required:    required,
			validate:    validate,
		})
	}

	intValidator := func(s string, min, max int) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return fmt.Errorf("must be a number")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("not a valid number")
		}
		if n < min || n > max {
			return fmt.Errorf("must be between %d and %d", min, max)
		}
		return nil
	}

	floatValidator := func(s string, min, max float64) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("not a valid amount")
		}
		if n < min || n > max {
			return fmt.Errorf("must be between %.2f and %.2f", min, max)
		}
		return nil
	}

	nonEmpty := func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("required")
		}
		return nil
	}

	thisYear := time.Now().Year()
	yearValidator := func(s string) error {
		return intValidator(s, 1880, thisYear+5)
	}

	priceValidator := func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return floatValidator(s, 0, 999.99)
	}

	platformValidator := func(s string) error {
		if m.mediaType != catalogGame {
			return nil
		}
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("required for games")
		}
		return nil
	}
	seasonValidator := func(s string) error {
		if m.mediaType != catalogSeries {
			return nil
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return intValidator(s, 1, 99)
	}
	episodesValidator := func(s string) error {
		if m.mediaType != catalogSeries {
			return nil
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return intValidator(s, 1, 9999)
	}

	formatValidator := func(s string) error {
		switch strings.TrimSpace(strings.ToUpper(s)) {
		case "VHS", "DVD", "BLU-RAY", "BLURAY":
			return nil
		}
		return fmt.Errorf("must be VHS, DVD, or Blu-ray")
	}

	genreValidator := func(s string) error {
		return nonEmpty(s)
	}

	if m.mode == FormEdit {
		switch m.mediaType {
		case catalogSeries:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "Breaking Bad", required: true, validate: nonEmpty,
			})
		case catalogGame:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "Halo", required: true, validate: nonEmpty,
			})
		default:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "The Matrix", required: true, validate: nonEmpty,
			})
		}
	} else {
		switch m.mediaType {
		case catalogSeries:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "Breaking Bad", required: true, validate: nonEmpty,
			})
		case catalogGame:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "Halo", required: true, validate: nonEmpty,
			})
		default:
			m.fields = append(m.fields, catalogField{
				label: "Title", placeholder: "The Matrix", required: true, validate: nonEmpty,
			})
		}
	}
	add("Year", "1999", "", true, yearValidator)
	add("Genre", "SciFi", "", true, genreValidator)
	add("Format", "VHS / DVD / Blu-ray", "", true, formatValidator)
	add("Director", "Wachowski", "", true, nonEmpty)
	add("Cast (comma-separated)", "Keanu Reeves, Laurence Fishburne", "", true, nonEmpty)
	add("Synopsis", "A computer hacker learns the true nature of reality", "", false, func(s string) error { return nil })
	add("Total copies", "5", "5", true, func(s string) error {
		return intValidator(s, 1, 9999)
	})
	add("Rental price (USD)", "4.99", "0", false, priceValidator)

	if m.mediaType == catalogSeries {
		add("Season", "1", "", false, seasonValidator)
		add("Episodes", "13", "", false, episodesValidator)
	}
	if m.mediaType == catalogGame {
		add("Platform", "PS5", "", true, platformValidator)
	}

	m.inputs = nil
	for _, f := range m.fields {
		ti := textinput.New()
		ti.Placeholder = f.placeholder
		ti.CharLimit = 300
		ti.Prompt = "▸ "
		ti.SetValue(f.value)
		m.inputs = append(m.inputs, ti)
	}
}

func (m *MovieFormModel) SetMediaType(mediaType string) {
	if mediaType == "" {
		mediaType = catalogMovie
	}
	if m.mediaType == mediaType {
		return
	}
	previous := make([]string, len(m.fields))
	for i, f := range m.fields {
		previous[i] = m.inputs[i].Value()
		f.value = m.inputs[i].Value()
	}
	m.mediaType = mediaType
	m.buildInputs()
	common := len(m.fields)
	if len(previous) < common {
		common = len(previous)
	}
	for i := 0; i < common; i++ {
		m.inputs[i].SetValue(previous[i])
	}
}

func (m *MovieFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *MovieFormModel) Update(msg tea.Msg) (*MovieFormModel, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, nil
		case "tab", "down":
			m.nextField()
			return m, nil
		case "shift+tab", "up":
			m.prevField()
			return m, nil
		case "1":
			if m.focus == 0 {
				m.SetMediaType(catalogMovie)
				return m, nil
			}
		case "2":
			if m.focus == 0 {
				m.SetMediaType(catalogSeries)
				return m, nil
			}
		case "3":
			if m.focus == 0 {
				m.SetMediaType(catalogGame)
				return m, nil
			}
		case "ctrl+s", "ctrl+enter":
			return m, m.trySubmit()
		}
	}

	if m.focus >= 0 && m.focus < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		m.fields[m.focus].value = m.inputs[m.focus].Value()
		return m, cmd
	}
	return m, nil
}

func (m *MovieFormModel) nextField() {
	if len(m.inputs) == 0 {
		return
	}
	if m.focus >= 0 {
		m.inputs[m.focus].Blur()
	}
	m.focus++
	if m.focus >= len(m.inputs) {
		m.focus = 0
	}
	m.inputs[m.focus].Focus()
}

func (m *MovieFormModel) prevField() {
	if len(m.inputs) == 0 {
		return
	}
	if m.focus >= 0 {
		m.inputs[m.focus].Blur()
	}
	m.focus--
	if m.focus < 0 {
		m.focus = len(m.inputs) - 1
	}
	m.inputs[m.focus].Focus()
}

func (m *MovieFormModel) trySubmit() tea.Cmd {
	for i, f := range m.fields {
		v := strings.TrimSpace(m.inputs[i].Value())
		if f.required && v == "" {
			m.ErrMsg = fmt.Sprintf("⛔ %s is required", f.label)
			m.focus = i
			m.inputs[i].Focus()
			for j := range m.inputs {
				if j != i {
					m.inputs[j].Blur()
				}
			}
			return nil
		}
		if err := f.validate(v); err != nil {
			m.ErrMsg = fmt.Sprintf("⛔ %s: %s", f.label, err)
			m.focus = i
			m.inputs[i].Focus()
			for j := range m.inputs {
				if j != i {
					m.inputs[j].Blur()
				}
			}
			return nil
		}
	}

	thisYear := time.Now().Year()
	year, _ := strconv.Atoi(strings.TrimSpace(m.inputs[1].Value()))
	copies, _ := strconv.Atoi(strings.TrimSpace(m.fieldValue("Total copies")))
	price, _ := strconv.ParseFloat(strings.TrimSpace(m.fieldValue("Rental price (USD)")), 64)

	msg := MovieFormSubmitMsg{
		Mode:      m.mode,
		MovieID:   m.movieID,
		MediaType: m.mediaType,
		Title:     strings.TrimSpace(m.fieldValue("Title")),
		Year:      year,
		Genre:     strings.TrimSpace(m.fieldValue("Genre")),
		Format:    m.normalizeFormat(m.fieldValue("Format")),
		Director:  strings.TrimSpace(m.fieldValue("Director")),
		Cast:      strings.TrimSpace(m.fieldValue("Cast (comma-separated)")),
		Synopsis:  strings.TrimSpace(m.fieldValue("Synopsis")),
		Copies:    copies,
		Price:     price,
	}

	if m.mediaType == catalogSeries {
		season, _ := strconv.Atoi(strings.TrimSpace(m.fieldValue("Season")))
		episodes, _ := strconv.Atoi(strings.TrimSpace(m.fieldValue("Episodes")))
		if season > 0 {
			msg.SeasonNumber = season
		}
		if episodes > 0 {
			msg.EpisodeCount = episodes
		}
	}
	if m.mediaType == catalogGame {
		msg.Platform = strings.TrimSpace(m.fieldValue("Platform"))
	}

	if year < 1900 || year > thisYear+5 {
		m.ErrMsg = fmt.Sprintf("⛔ Year must be between 1900 and %d", thisYear+5)
		return nil
	}
	if msg.Copies < 1 {
		m.ErrMsg = "⛔ Total copies must be at least 1"
		return nil
	}

	m.ErrMsg = ""
	return func() tea.Msg { return msg }
}

func (m *MovieFormModel) fieldValue(label string) string {
	for i, f := range m.fields {
		if f.label == label {
			return m.inputs[i].Value()
		}
	}
	return ""
}

func (m *MovieFormModel) normalizeFormat(s string) string {
	s = strings.TrimSpace(s)
	switch strings.ToUpper(s) {
	case "BLURAY":
		return "Blu-ray"
	}
	return s
}

func (m *MovieFormModel) View(w, h int) string {
	formTitle := "─── ADD CATALOG ITEM ───"
	if m.mode == FormEdit {
		formTitle = "─── EDIT CATALOG ITEM ───"
	}

	header := styles.TitleStyle.
		Width(64).
		Align(lipgloss.Center).
		Render(formTitle)

	typeLabel := map[string]string{
		catalogMovie:  "🎬 Movie",
		catalogSeries: "📺 Series",
		catalogGame:   "🕹️ Game",
	}
	current := typeLabel[m.mediaType]
	if current == "" {
		current = "🎬 Movie"
	}
	typeBar := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Width(64).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("[1] 🎬  [2] 📺  [3] 🕹️      currently: %s", current))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Padding(1, 3).
		Width(64)

	fieldLines := []string{}
	for i, f := range m.fields {
		if i >= len(m.inputs) {
			break
		}
		label := styles.DimTextStyle.Render(f.label)
		fieldLines = append(fieldLines, lipgloss.JoinHorizontal(lipgloss.Left, "  ", label))
		fieldLines = append(fieldLines, "  "+m.inputs[i].View())
		fieldLines = append(fieldLines, "")
	}
	body := lipgloss.JoinVertical(lipgloss.Left, fieldLines...)
	boxed := box.Render(body)

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		typeBar,
		"",
		boxed,
	)

	if m.ErrMsg != "" {
		content += "\n" + styles.ErrorTextStyle.Render(m.ErrMsg)
	}

	help := styles.DimTextStyle.
		Width(64).
		Align(lipgloss.Center).
		Render("tab/↓ next · shift+tab/↑ prev · [1][2][3] type · ctrl+s submit · esc cancel")

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

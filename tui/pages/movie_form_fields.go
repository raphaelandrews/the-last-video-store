package pages

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/thelastvideostore/internal/models"
)

func (m *MovieFormModel) buildFormFields(thisYear int) []huh.Field {
	titlePH := "The Matrix"
	if m.mediaType == catalogSeries {
		titlePH = "Breaking Bad"
	} else if m.mediaType == catalogGame {
		titlePH = "Halo"
	}

	genreOptions := make([]huh.Option[string], len(m.genOptions))
	for i, g := range m.genOptions {
		genreOptions[i] = huh.NewOption(g, g)
	}
	formatOptions := make([]huh.Option[string], len(m.fmtOptions))
	for i, f := range m.fmtOptions {
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

	return fields
}

func (m *MovieFormModel) buildForm() {
	m.form = huh.NewForm(huh.NewGroup(m.buildFormFields(time.Now().Year())...)).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(gruvboxHuhTheme()).
		WithKeyMap(gruvboxKeyMap())
}

func (m *MovieFormModel) SeedValues(mv *models.MovieResponse) {
	m.title = mv.Title
	m.year = itoa(mv.Year)
	m.genre = mv.Genre
	m.format = mv.Format
	m.director = mv.Director
	m.cast = joinCast(mv.Cast)
	m.synopsis = mv.Synopsis
	m.copies = itoa(mv.CopiesTotal)
	m.price = formatFloat(mv.RentalPrice)
	m.season = itoa(mv.SeasonNumber)
	m.episodes = itoa(mv.EpisodeCount)
	m.platform = mv.Platform
}

func (m *MovieFormModel) ToMessage() MovieFormSubmitMsg {
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

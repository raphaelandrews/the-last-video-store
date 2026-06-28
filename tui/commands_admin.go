package tui

import (
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadAdminMovies(page int) tea.Cmd {
	ps := m.adminMovies.PageSize
	if ps <= 0 {
		ps = 30
	}
	mediaType := m.adminMovies.ActiveTab()
	return func() tea.Msg {
		if m.userResp == nil || !bitmask.CanAdmin(m.userResp.Tier) {
			return pages.ErrorMsg{Message: "⛔ ACCESS DENIED — Manager or Owner required"}
		}
		url := fmt.Sprintf("%s/api/v1/movies?page_size=%d&page=%d", m.baseURL, ps, page)
		if mediaType != "" {
			url += "&media_type=" + string(mediaType)
		}
		resp, err := m.apiGetURL(url)
		if err != nil {
			return loadAdminMoviesMsg{page: page, errMsg: err.Error()}
		}
		if resp == nil {
			return loadAdminMoviesMsg{page: page, errMsg: "no response from server"}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return loadAdminMoviesMsg{page: page, errMsg: fmt.Sprintf("server returned %d", resp.StatusCode)}
		}
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return loadAdminMoviesMsg{page: page, errMsg: "failed to decode response"}
		}
		return loadAdminMoviesMsg{movies: r.Movies, total: r.Total, page: page}
	}
}

func (m *Model) loadCatalogOptions() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/movies/options")
		if resp == nil {
			return loadCatalogOptionsMsg{}
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var r struct {
			Genres  []string `json:"genres"`
			Formats []string `json:"formats"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadCatalogOptionsMsg{genres: r.Genres, formats: r.Formats}
	}
}

func (m *Model) moveToAdminMovies() {
	m.screen = scrAdminMovies
	m.movieForm = nil
}

func (m *Model) doCreateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"media_type":   msg.MediaType,
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"platform":     msg.Platform,
			"season":       msg.SeasonNumber,
			"episodes":     msg.EpisodeCount,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		resp, err := m.apiPost("/api/v1/movies", string(body))
		if err != nil {
			m.movieForm.SetError(err.Error())
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.movieForm.SetError(e.Error)
			return nil
		}
		m.moveToAdminMovies()
		m.adminMovies.CurrentPageFor(m.adminMovies.ActiveTab())
		return m.loadAdminMovies(1)()
	}
}

func (m *Model) doUpdateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"media_type":   msg.MediaType,
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"platform":     msg.Platform,
			"season":       msg.SeasonNumber,
			"episodes":     msg.EpisodeCount,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		resp, err := m.apiPut("/api/v1/movies/"+msg.MovieID, string(body))
		if err != nil {
			m.movieForm.SetError(err.Error())
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.movieForm.SetError(e.Error)
			return nil
		}
		m.moveToAdminMovies()
		m.adminMovies.CurrentPageFor(m.adminMovies.ActiveTab())
		return m.loadAdminMovies(1)()
	}
}

func (m *Model) doToggleStaffPick(movieID string, current bool) tea.Cmd {
	return func() tea.Msg {
		var resp *http.Response
		if current {
			resp, _ = m.apiDelete("/api/v1/movies/" + movieID + "/staff-pick")
		} else {
			resp, _ = m.apiPostEmpty("/api/v1/movies/" + movieID + "/staff-pick")
		}
		if resp != nil {
			defer resp.Body.Close()
			if errMsg := handleAPIErr(resp); errMsg != nil {
				return errMsg
			}
			if resp.StatusCode >= 400 {
				var e struct {
					Error string `json:"error"`
				}
				json.NewDecoder(resp.Body).Decode(&e)
				if e.Error != "" {
					return pages.ErrorMsg{Message: e.Error}
				}
			}
		}
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doDeleteMovie(movieID string) tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiDelete("/api/v1/movies/" + movieID)
		if resp != nil {
			defer resp.Body.Close()
			if errMsg := handleAPIErr(resp); errMsg != nil {
				return errMsg
			}
			if resp.StatusCode >= 400 {
				var e struct {
					Error string `json:"error"`
				}
				json.NewDecoder(resp.Body).Decode(&e)
				if e.Error != "" {
					return pages.ErrorMsg{Message: e.Error}
				}
			}
		}
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

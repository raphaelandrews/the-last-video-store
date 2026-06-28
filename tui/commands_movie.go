package tui

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doSearch(query string) tea.Cmd {
	if query == "" {
		return nil
	}
	return func() tea.Msg {
		resp, _ := m.apiGetURL(m.baseURL + "/api/v1/movies/search?q=" + query)
		if resp == nil {
			return searchResultsMsg{}
		}
		defer resp.Body.Close()
		var results []models.MovieResponse
		json.NewDecoder(resp.Body).Decode(&results)
		return searchResultsMsg{results: results}
	}
}

func (m *Model) loadMovies(page int, genre string) tea.Cmd {
	ps := m.browse.PageSize
	if ps <= 0 {
		ps = 20
	}
	return func() tea.Msg {
		m.browseReqID++
		rid := m.browseReqID
		url := fmt.Sprintf("%s/api/v1/movies?page_size=%d&page=%d", m.baseURL, ps, page)
		if genre != "" {
			url += "&genre=" + genre
		}
		if m.browse.MediaType != "" {
			url += "&media_type=" + m.browse.MediaType
		}
		resp, _ := m.apiGetURL(url)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMoviesMsg{movies: r.Movies, total: r.Total, page: page, reqID: rid}
	}
}

func (m *Model) loadStaffPicks() tea.Cmd {
	return func() tea.Msg {
		m.browseReqID++
		rid := m.browseReqID
		url := m.baseURL + "/api/v1/movies/staff-picks"
		if m.browse.MediaType != "" {
			url += "?media_type=" + m.browse.MediaType
		}
		resp, _ := m.apiGetURL(url)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var movies []models.MovieResponse
		json.NewDecoder(resp.Body).Decode(&movies)
		return loadMoviesMsg{movies: movies, total: len(movies), page: 1, reqID: rid}
	}
}

func (m *Model) loadLastChance() tea.Cmd {
	return func() tea.Msg {
		m.browseReqID++
		rid := m.browseReqID
		url := m.baseURL + "/api/v1/movies/last-chance"
		if m.browse.MediaType != "" {
			url += "?media_type=" + m.browse.MediaType
		}
		resp, _ := m.apiGetURL(url)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var movies []models.MovieResponse
		json.NewDecoder(resp.Body).Decode(&movies)
		return loadMoviesMsg{movies: movies, total: len(movies), page: 1, reqID: rid}
	}
}

func (m *Model) doRent(movieID string) tea.Cmd {
	return func() tea.Msg {
		useTicket := m.detail != nil && m.detail.UseTicket
		body, _ := json.Marshal(map[string]interface{}{
			"movie_id":   movieID,
			"use_ticket": useTicket,
		})
		resp, err := m.apiPost("/api/v1/rentals/rent", string(body))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var rental models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rental)
		if m.detail != nil {
			m.detail.SetRental(&rental)
			m.detail.Movie.CopiesAvailable--
			if m.detail.Movie.CopiesAvailable <= 0 {
				m.detail.Movie.Available = false
			}
		}
		if m.gameDetail != nil && m.gameDetail.Game != nil {
			m.gameDetail.SetRental(&rental)
			m.gameDetail.Game.CopiesAvailable--
			if m.gameDetail.Game.CopiesAvailable <= 0 {
				m.gameDetail.Game.Available = false
			}
		}
		if m.userResp != nil {
			if m.detail != nil && m.detail.UseTicket {
				m.userResp.FreeRentals--
				m.browse.Status = "Rented! " + rental.MovieTitle + " (🎟️ free rental)"
			} else {
				cost := models.MovieCost(0, rental.MovieFormat)
				m.userResp.Balance -= cost
				m.userResp.RentalCount++
				m.browse.Status = fmt.Sprintf("Rented! %s (💵 $%.2f)", rental.MovieTitle, cost)
			}
			if m.detail != nil {
				m.detail.UseTicket = false
			}
		}
		return m.loadMovies(m.browse.Page, m.browse.Genre)()
	}
}

func (m *Model) doAddToWishlist(movieID string, fromDetail bool) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"movie_id": movieID})
		resp, err := m.apiPost("/api/v1/wishlist", string(body))
		if err != nil {
			if fromDetail {
				m.detail.ErrMsg = err.Error()
			} else {
				m.browse.Status = err.Error()
			}
			return wishlistResultMsg{}
		}
		defer resp.Body.Close()
		if resp.StatusCode == 409 {
			if fromDetail {
				m.detail.ErrMsg = "Already in waitlist"
			} else {
				m.browse.Status = "Already in waitlist"
			}
			return wishlistResultMsg{}
		}
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			if fromDetail {
				m.detail.ErrMsg = e.Error
			} else {
				m.browse.Status = e.Error
			}
			return wishlistResultMsg{}
		}
		if fromDetail {
			m.detail.StatusMsg = "Added to waitlist ✓"
		} else {
			m.browse.Status = "Added to waitlist ✓"
		}
		return wishlistResultMsg{}
	}
}

func (m *Model) doRefreshDetail(movieID string) tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/movies/" + movieID)
		if resp == nil {
			return nil
		}
		defer resp.Body.Close()
		var movie models.MovieResponse
		json.NewDecoder(resp.Body).Decode(&movie)
		return refreshDetailMsg{movie: &movie}
	}
}

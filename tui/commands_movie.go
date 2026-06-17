package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doSearch(query string) tea.Cmd {
	if query == "" {
		return nil
	}
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/search?q="+query, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
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
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/staff-picks", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/last-chance", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		body := fmt.Sprintf(`{"movie_id":"%s","use_ticket":%v}`, movieID, useTicket)
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/rentals/rent", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			return pages.ErrorMsg{Message: e.Error}
		}
		var rental models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rental)
		if m.detail != nil {
			m.detail.SetRental(&rental)
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
			m.detail.UseTicket = false
		}
		return m.loadMovies(m.browse.Page, m.browse.Genre)()
	}
}

func (m *Model) doAddToWishlist(movieID string, fromDetail bool) tea.Cmd {
	return func() tea.Msg {
		body := `{"movie_id":"` + movieID + `"}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/wishlist", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if fromDetail {
				m.detail.ErrMsg = err.Error()
			} else {
				m.browse.Status = err.Error()
			}
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode == 409 {
			if fromDetail {
				m.detail.ErrMsg = "Already in waitlist"
			} else {
				m.browse.Status = "Already in waitlist"
			}
			return nil
		}
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			if fromDetail {
				m.detail.ErrMsg = e.Error
			} else {
				m.browse.Status = e.Error
			}
			return nil
		}
		if fromDetail {
			m.detail.StatusMsg = "Added to waitlist ✓"
		} else {
			m.browse.Status = "Added to waitlist ✓"
		}
		return nil
	}
}

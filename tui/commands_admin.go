package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	return func() tea.Msg {
		if m.userResp == nil || !bitmask.CanAdmin(m.userResp.Tier) {
			return pages.ErrorMsg{Message: "⛔ ACCESS DENIED — Manager or Owner required"}
		}
		url := fmt.Sprintf("%s/api/v1/movies?page_size=%d&page=%d", m.baseURL, ps, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAdminMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadAdminMoviesMsg{movies: r.Movies, total: r.Total, page: page}
	}
}

func (m *Model) loadAdminUsers() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAdminUsersMsg{}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if e.Error == "" {
				e.Error = "ACCESS DENIED — Insufficient clearance"
			}
			return pages.ErrorMsg{Message: e.Error}
		}
		var users []models.UserResponse
		json.NewDecoder(resp.Body).Decode(&users)
		return loadAdminUsersMsg{users: users}
	}
}

func (m *Model) doVerifyAuditChain() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/audit/verify", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.auditLog.VerifyMsg = "⚠️ Verification failed: " + err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			ChainIntact bool   `json:"chain_intact"`
			Message     string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		m.auditLog.VerifyMsg = r.Message
		return nil
	}
}

func (m *Model) loadAuditLog() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/audit", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAuditLogMsg{}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if e.Error == "" {
				e.Error = "ACCESS DENIED — Insufficient clearance"
			}
			return pages.ErrorMsg{Message: e.Error}
		}
		var entries []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&entries)
		return loadAuditLogMsg{entries: entries}
	}
}

func (m *Model) doCreateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/movies", strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
				return pages.ErrorMsg{Message: e.Error}
			}
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doUpdateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		req, _ := http.NewRequest("PUT", m.baseURL+"/api/v1/movies/"+msg.MovieID, strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
				return pages.ErrorMsg{Message: e.Error}
			}
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doUpdateUser(userID, action string) tea.Cmd {
	return func() tea.Msg {
		var body string
		u := m.adminUsers.SelectedUser()
		if action == "promote" {
			next := nextTier(u.TierName)
			body = `{"tier":"` + next + `"}`
		} else if action == "demote" {
			prev := prevTier(u.TierName)
			body = `{"tier":"` + prev + `"}`
		} else if action == "ban" {
			if u.Banned {
				body = `{"banned":false}`
			} else {
				body = `{"banned":true}`
			}
		}
		req, _ := http.NewRequest("PUT", m.baseURL+"/api/v1/users/"+userID, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.adminUsers.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
				return pages.ErrorMsg{Message: e.Error}
			}
			m.adminUsers.ErrMsg = e.Error
			return nil
		}
		m.adminUsers.ErrMsg = ""
		return m.loadAdminUsers()()
	}
}

func (m *Model) doToggleStaffPick(movieID string, current bool) tea.Cmd {
	method := "POST"
	if current {
		method = "DELETE"
	}
	return func() tea.Msg {
		req, _ := http.NewRequest(method, m.baseURL+"/api/v1/movies/"+movieID+"/staff-pick", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
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
		req, _ := http.NewRequest("DELETE", m.baseURL+"/api/v1/movies/"+movieID, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
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

func (m *Model) moveToAdminMovies() {
	m.screen = scrAdminMovies
	m.movieForm = nil
}

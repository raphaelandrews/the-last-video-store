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
		resp, _ := m.apiGetURL(url)
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
		resp, _ := m.apiGet("/api/v1/users")
		if resp == nil {
			return loadAdminUsersMsg{}
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var users []models.UserResponse
		json.NewDecoder(resp.Body).Decode(&users)
		return loadAdminUsersMsg{users: users}
	}
}

func (m *Model) doVerifyAuditChain() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.apiGet("/api/v1/audit/verify")
		if err != nil {
			m.auditLog.VerifyMsg = "⚠️ Verification failed: " + err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			ChainIntact bool   `json:"chain_intact"`
			Message     string `json:"message"`
			BrokenAt    int    `json:"broken_at"`
			Reason      string `json:"reason"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.ChainIntact {
			m.auditLog.BrokenAt = -1
		} else {
			m.auditLog.BrokenAt = r.BrokenAt
		}
		m.auditLog.VerifyMsg = r.Message
		return nil
	}
}

func (m *Model) loadAuditLog() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/audit")
		if resp == nil {
			return loadAuditLogMsg{}
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
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
			"media_type":   msg.MediaType,
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"platform":     msg.Platform,
			"season":       msg.Season,
			"episodes":     msg.Episodes,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		resp, err := m.apiPost("/api/v1/movies", string(body))
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
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
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.CurrentPageFor(m.adminMovies.ActiveTab()))()
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
			"season":       msg.Season,
			"episodes":     msg.Episodes,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		resp, err := m.apiPut("/api/v1/movies/"+msg.MovieID, string(body))
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
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
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.CurrentPageFor(m.adminMovies.ActiveTab()))()
	}
}

func (m *Model) doUpdateUser(userID, action string) tea.Cmd {
	return func() tea.Msg {
		u := m.adminUsers.SelectedUser()
		if u == nil {
			return nil
		}
		var body string
		switch action {
		case "promote":
			next, ok := canPromote(u.TierName)
			if !ok {
				m.adminUsers.StatusMsg = fmt.Sprintf("⛔ %s is already at the highest tier", u.TierName)
				return nil
			}
			body = `{"tier":"` + next + `"}`
		case "demote":
			prev, ok := canDemote(u.TierName)
			if !ok {
				m.adminUsers.StatusMsg = fmt.Sprintf("⛔ %s is already at the lowest tier", u.TierName)
				return nil
			}
			body = `{"tier":"` + prev + `"}`
		case "ban":
			if u.Banned {
				body = `{"banned":false}`
			} else {
				body = `{"banned":true}`
			}
		default:
			return nil
		}
		resp, err := m.apiPut("/api/v1/users/"+userID, body)
		if err != nil {
			m.adminUsers.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if apiErr, ok := decodeAPIErr(resp); ok {
			m.adminUsers.ErrMsg = apiErr
			return nil
		}
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if e.Error != "" {
				m.adminUsers.ErrMsg = e.Error
			} else {
				m.adminUsers.ErrMsg = fmt.Sprintf("server returned %d", resp.StatusCode)
			}
			return nil
		}
		m.adminUsers.ErrMsg = ""
		switch action {
		case "promote":
			m.adminUsers.StatusMsg = fmt.Sprintf("✅ Promoted %s", u.TierName)
		case "demote":
			m.adminUsers.StatusMsg = fmt.Sprintf("✅ Demoted %s", u.TierName)
		case "ban":
			if u.Banned {
				m.adminUsers.StatusMsg = fmt.Sprintf("✅ Unbanned %s", u.Username)
			} else {
				m.adminUsers.StatusMsg = fmt.Sprintf("🚫 Banned %s", u.Username)
			}
		}
		return m.loadAdminUsers()()
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

func (m *Model) moveToAdminMovies() {
	m.screen = scrAdminMovies
	m.movieForm = nil
}

func (m *Model) doTOTPToggle(userID string) tea.Cmd {
	return func() tea.Msg {
		u := m.adminUsers.SelectedUser()
		if u == nil {
			return nil
		}
		enabled := !u.TOTPEnabled
		body := fmt.Sprintf(`{"enabled":%v}`, enabled)
		resp, err := m.apiPost("/api/v1/users/"+userID+"/totp", body)
		if err != nil {
			m.adminUsers.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if enabled {
			var r struct {
				Secret string `json:"secret"`
				URL    string `json:"url"`
			}
			json.NewDecoder(resp.Body).Decode(&r)
			m.adminUsers.StatusMsg = fmt.Sprintf("🔒 TOTP enabled — secret: %s", r.Secret)
		} else {
			m.adminUsers.StatusMsg = "🔓 TOTP disabled"
		}
		return m.loadAdminUsers()()
	}
}

func (m *Model) doProfileTOTP() tea.Cmd {
	return func() tea.Msg {
		if m.userResp == nil {
			return nil
		}
		enabled := !m.userResp.TOTPEnabled
		body := fmt.Sprintf(`{"enabled":%v}`, enabled)
		resp, err := m.apiPost("/api/v1/users/"+m.userResp.ID+"/totp", body)
		if err != nil {
			m.profile.StatusMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if enabled {
			var r struct {
				Secret string `json:"secret"`
				URL    string `json:"url"`
			}
			json.NewDecoder(resp.Body).Decode(&r)
			m.userResp.TOTPEnabled = true
			m.profile.StatusMsg = fmt.Sprintf("🔒 TOTP enabled\nSecret: %s\nScan this into your authenticator app", r.Secret)
		} else {
			m.userResp.TOTPEnabled = false
			m.profile.StatusMsg = "🔓 TOTP disabled"
		}
		return tea.Sequence(
			func() tea.Msg { return wishlistResultMsg{} },
			m.doRefreshMe(),
		)
	}
}

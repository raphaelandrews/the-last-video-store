package tui

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
)

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

func (m *Model) doUpdateUser(userID, action string) tea.Cmd {
	return func() tea.Msg {
		u := m.adminUsers.SelectedUser()
		if u == nil {
			return nil
		}
		var body []byte
		switch action {
		case "promote":
			next, ok := canPromote(u.TierName)
			if !ok {
				m.adminUsers.StatusMsg = fmt.Sprintf("⛔ %s is already at the highest tier", u.TierName)
				return nil
			}
			body, _ = json.Marshal(map[string]string{"tier": next})
		case "demote":
			prev, ok := canDemote(u.TierName)
			if !ok {
				m.adminUsers.StatusMsg = fmt.Sprintf("⛔ %s is already at the lowest tier", u.TierName)
				return nil
			}
			body, _ = json.Marshal(map[string]string{"tier": prev})
		case "ban":
			body, _ = json.Marshal(map[string]bool{"banned": !u.Banned})
		default:
			return nil
		}
		resp, err := m.apiPut("/api/v1/users/"+userID, string(body))
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

func (m *Model) doTOTPToggle(userID string) tea.Cmd {
	return func() tea.Msg {
		u := m.adminUsers.SelectedUser()
		if u == nil {
			return nil
		}
		enabled := !u.TOTPEnabled
		body, _ := json.Marshal(map[string]bool{"enabled": enabled})
		resp, err := m.apiPost("/api/v1/users/"+userID+"/totp", string(body))
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

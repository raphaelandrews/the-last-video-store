package tui

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) doProfileTOTP() tea.Cmd {
	return func() tea.Msg {
		if m.userResp == nil {
			return nil
		}
		enabled := !m.userResp.TOTPEnabled
		body, _ := json.Marshal(map[string]bool{"enabled": enabled})
		resp, err := m.apiPost("/api/v1/users/"+m.userResp.ID+"/totp", string(body))
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

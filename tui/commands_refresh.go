package tui

import (
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
)

func (m *Model) doRefreshMe() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return nil
		}
		var user models.UserResponse
		json.NewDecoder(resp.Body).Decode(&user)
		return refreshMeMsg{user: &user}
	}
}

package tui

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
)

func (m *Model) doRefreshMe() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/auth/me")
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

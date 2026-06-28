package tui

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doGamePlayStart(gameID string, durationMinutes int) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]interface{}{
			"game_id":          gameID,
			"duration_minutes": durationMinutes,
		})
		resp, err := m.apiPost("/api/v1/games/play/start", string(body))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			Error   string             `json:"error"`
			Message string             `json:"message"`
			Session models.GameSession `json:"session"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if resp.StatusCode >= 400 || r.Error != "" {
			return pages.ErrorMsg{Message: r.Error}
		}
		m.gameDetail.SetSession(&r.Session)
		return tea.Sequence(
			func() tea.Msg { return gameRefreshMsg{} },
			m.doRefreshMe(),
		)
	}
}

func (m *Model) doGamePlayEnd(sessionID string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"session_id": sessionID})
		resp, err := m.apiPost("/api/v1/games/play/end", string(body))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			Error   string             `json:"error"`
			Message string             `json:"message"`
			Session models.GameSession `json:"session"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if resp.StatusCode >= 400 || r.Error != "" {
			return pages.ErrorMsg{Message: r.Error}
		}
		m.gameDetail.Playing = false
		m.gameDetail.Session = nil
		return tea.Batch(
			func() tea.Msg { return gameRefreshMsg{} },
			m.doRefreshMe(),
		)
	}
}

package tui

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doGamePlayStart(gameID string, gameTitle string, durationMinutes int) tea.Cmd {
	return func() tea.Msg {
		body := fmt.Sprintf(`{"game_id":"%s","duration_minutes":%d}`, gameID, durationMinutes)
		resp, err := m.apiPost("/api/v1/games/play/start", body)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			Error   string             `json:"error"`
			Message string             `json:"message"`
			Session models.GameSession `json:"session"`
			Rate    float64            `json:"rate"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if resp.StatusCode >= 400 || r.Error != "" {
			return pages.ErrorMsg{Message: r.Error}
		}
		m.gameDetail.SetSession(&r.Session)
		if m.userResp != nil {
			m.userResp.Balance -= r.Session.Cost
		}
		return tea.Sequence(
			func() tea.Msg { return gameRefreshMsg{} },
			m.doRefreshMe(),
		)
	}
}

func (m *Model) doGamePlayEnd(sessionID string) tea.Cmd {
	return func() tea.Msg {
		body := fmt.Sprintf(`{"session_id":"%s"}`, sessionID)
		resp, err := m.apiPost("/api/v1/games/play/end", body)
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

func (m *Model) loadGameSessions() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/games/play/active")
		if resp == nil {
			return loadGameSessionsMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Sessions []models.GameSession `json:"sessions"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		for i, j := 0, len(r.Sessions)-1; i < j; i, j = i+1, j-1 {
			r.Sessions[i], r.Sessions[j] = r.Sessions[j], r.Sessions[i]
		}
		return loadGameSessionsMsg{sessions: r.Sessions}
	}
}

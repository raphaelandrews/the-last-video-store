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

func (m *Model) doGamePlayStart(gameID string, gameTitle string) tea.Cmd {
	return func() tea.Msg {
		body := fmt.Sprintf(`{"game_id":"%s"}`, gameID)
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/games/play/start", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
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
		if m.userResp != nil && r.Rate > 0 {
			m.userResp.Balance -= r.Rate
		}
		return gameRefreshMsg{}
	}
}

func (m *Model) doGamePlayEnd(sessionID string) tea.Cmd {
	return func() tea.Msg {
		body := fmt.Sprintf(`{"session_id":"%s"}`, sessionID)
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/games/play/end", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/games/play/active", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadGameSessionsMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Sessions []models.GameSession `json:"sessions"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadGameSessionsMsg{sessions: r.Sessions}
	}
}

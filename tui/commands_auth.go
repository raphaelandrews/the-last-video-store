package tui

import (
	"encoding/json"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doSubmitTOTP(tempToken, code string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"code": code})
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/auth/login/totp", strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tempToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			AccessToken  string              `json:"access_token"`
			RefreshToken string              `json:"refresh_token"`
			Error        string              `json:"error"`
			User         models.UserResponse `json:"user"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "invalid code"
			}
			m.totpCode = ""
			return pages.ErrorMsg{Message: r.Error}
		}
		return pages.LoginSuccessMsg{
			AccessToken:  r.AccessToken,
			RefreshToken: r.RefreshToken,
			User:         &r.User,
		}
	}
}

func (m *Model) doLogin(u, p string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"username": u, "password": p})
		resp, err := http.Post(m.baseURL+"/api/v1/auth/login", "application/json", strings.NewReader(string(body)))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			AccessToken  string              `json:"access_token"`
			RefreshToken string              `json:"refresh_token"`
			Error        string              `json:"error"`
			User         models.UserResponse `json:"user"`
			TOTPRequired bool                `json:"totp_required"`
			TempToken    string              `json:"temp_token"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.TOTPRequired && r.TempToken != "" {
			m.tempToken = r.TempToken
			m.totpCode = ""
			m.screen = scrTOTP
			return nil
		}
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "invalid credentials"
			}
			return pages.ErrorMsg{Message: r.Error}
		}
		return pages.LoginSuccessMsg{AccessToken: r.AccessToken, RefreshToken: r.RefreshToken, User: &r.User}
	}
}

func (m *Model) doRegister(u, p string) tea.Cmd {
	return func() tea.Msg {
		if len(u) < 3 {
			return pages.ErrorMsg{Message: "Username must be at least 3 characters"}
		}
		if len(p) < 6 {
			return pages.ErrorMsg{Message: "Password must be at least 6 characters"}
		}
		body, _ := json.Marshal(map[string]string{"username": u, "password": p})
		resp, err := http.Post(m.baseURL+"/api/v1/auth/register", "application/json", strings.NewReader(string(body)))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "registration failed"
			}
			return pages.ErrorMsg{Message: r.Error}
		}
		m.register = nil
		return pages.NavigateMsg{Page: "login"}
	}
}

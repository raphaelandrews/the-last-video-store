package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) apiGet(path string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", m.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+m.token)
	return http.DefaultClient.Do(req)
}

func (m *Model) apiGetURL(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+m.token)
	return http.DefaultClient.Do(req)
}

func (m *Model) apiPost(path, body string) (*http.Response, error) {
	return m.apiPostWithToken(path, body, m.token)
}

func (m *Model) apiPostWithToken(path, body, token string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", m.baseURL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

func (m *Model) apiPostEmpty(path string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", m.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+m.token)
	return http.DefaultClient.Do(req)
}

func (m *Model) apiPut(path, body string) (*http.Response, error) {
	req, _ := http.NewRequest("PUT", m.baseURL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)
	return http.DefaultClient.Do(req)
}

func (m *Model) apiDelete(path string) (*http.Response, error) {
	req, _ := http.NewRequest("DELETE", m.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+m.token)
	return http.DefaultClient.Do(req)
}

func decodeAPIErr(resp *http.Response) (string, bool) {
	if resp == nil || resp.StatusCode < 400 {
		return "", false
	}
	var e struct {
		Error string `json:"error"`
	}
	if resp.Body != nil {
		_ = json.NewDecoder(resp.Body).Decode(&e)
	}
	if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
		return e.Error, true
	}
	return e.Error, false
}

func handleAPIErr(resp *http.Response) tea.Msg {
	if resp == nil || resp.StatusCode < 400 {
		return nil
	}
	var e struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&e)
	if e.Error == "" {
		e.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return pages.ErrorMsg{Message: e.Error}
}

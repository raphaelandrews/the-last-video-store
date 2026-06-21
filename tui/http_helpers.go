package tui

import (
	"encoding/json"
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
	req, _ := http.NewRequest("POST", m.baseURL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)
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
	json.NewDecoder(resp.Body).Decode(&e)
	if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
		return e.Error, true
	}
	return e.Error, false
}

func handleAPIErr(resp *http.Response) tea.Msg {
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		return nil
	}
	var e struct {
		Error string `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&e)
	if strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") || e.Error != "" {
		return pages.ErrorMsg{Message: e.Error}
	}
	return nil
}

func handleAPIErrInline(resp *http.Response, fallbackErr string) string {
	msg, _ := decodeAPIErr(resp)
	if msg == "" {
		msg = fallbackErr
	}
	return msg
}

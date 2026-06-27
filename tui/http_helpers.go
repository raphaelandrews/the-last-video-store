package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

// logAPIError writes a one-line diagnostic to stderr. The X-Request-ID
// header is included so the line can be cross-referenced with server logs.
func logAPIError(tag, method, path string, resp *http.Response, body string) {
	if resp == nil {
		fmt.Fprintf(os.Stderr, "[tui %s] %s %s: no response (network error)\n", tag, method, path)
		return
	}
	reqID := resp.Header.Get("X-Request-ID")
	fmt.Fprintf(os.Stderr, "[tui %s] %s %s → %d  req_id=%s  body=%s\n",
		tag, method, path, resp.StatusCode, reqID, strings.TrimSpace(body))
}

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

// decodeAPIErr returns the {error} body for any non-2xx response.
// The bool is true when the response is a 403 (access denied).
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

// handleAPIErr turns a 4xx/5xx response into a pages.ErrorMsg. Returns
// nil for 2xx.
func handleAPIErr(resp *http.Response) tea.Msg {
	if resp == nil {
		return nil
	}
	if resp.StatusCode < 400 {
		return nil
	}
	var e struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&e)
	if e.Error == "" {
		e.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	if strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
		return pages.ErrorMsg{Message: e.Error}
	}
	return pages.ErrorMsg{Message: e.Error}
}

func handleAPIErrInline(resp *http.Response, fallbackErr string) string {
	msg, _ := decodeAPIErr(resp)
	if msg == "" {
		msg = fallbackErr
	}
	return msg
}

func logFailedRequest(tag, method, path string, resp *http.Response) {
	if resp == nil {
		return
	}
	logAPIError(tag, method, path, resp, "")
}

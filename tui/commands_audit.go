package tui

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) loadAuditLog() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/audit")
		if resp == nil {
			return loadAuditLogMsg{}
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var entries []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&entries)
		return loadAuditLogMsg{entries: entries}
	}
}

func (m *Model) doVerifyAuditChain() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.apiGet("/api/v1/audit/verify")
		if err != nil {
			m.auditLog.VerifyMsg = "⚠️ Verification failed: " + err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			ChainIntact bool   `json:"chain_intact"`
			Message     string `json:"message"`
			BrokenAt    int    `json:"broken_at"`
			BrokenID    string `json:"broken_id"`
			Reason      string `json:"reason"`
			Total       int    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.ChainIntact {
			m.auditLog.MarkIntact(r.Total)
		} else {
			m.auditLog.MarkBroken(r.BrokenAt, r.BrokenID, r.Reason)
		}
		m.auditLog.VerifyMsg = r.Message
		return nil
	}
}

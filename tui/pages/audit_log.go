package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type AuditLogModel struct {
	entries   []map[string]interface{}
	selected  int
	scroll    int
	errMsg    string
	VerifyMsg string
}

type AuditLogRefreshMsg struct{}

func NewAuditLogModel() *AuditLogModel {
	return &AuditLogModel{selected: -1}
}

func (m *AuditLogModel) SetEntries(entries []map[string]interface{}) {
	m.entries = entries
}

func (m *AuditLogModel) MoveUp() {
	if m.selected > 0 {
		m.selected--
	}
	if m.selected < m.scroll {
		m.scroll = m.selected
	}
}

func (m *AuditLogModel) MoveDown() {
	if m.selected < len(m.entries)-1 {
		m.selected++
	}
	if m.selected >= m.scroll+20 {
		m.scroll = m.selected - 20 + 1
	}
}

func (m *AuditLogModel) PageUp() {
	m.scroll -= 10
	if m.scroll < 0 {
		m.scroll = 0
	}
	m.selected = m.scroll
}

func (m *AuditLogModel) PageDown() {
	m.scroll += 10
	maxScroll := len(m.entries) - 20
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
	m.selected = m.scroll
}

func (m *AuditLogModel) View(width, height int) string {
	title := styles.HeadingStyle.Width(width).Align(lipgloss.Center).Render("🔗 AUDIT LOG — Hash Chain Viewer")

	status := styles.DimTextStyle.Render("Press [V] to verify chain integrity")
	if m.VerifyMsg != "" {
		if strings.HasPrefix(m.VerifyMsg, "✅") {
			status = styles.SuccessTextStyle.Render(m.VerifyMsg)
		} else {
			status = styles.ErrorTextStyle.Render(m.VerifyMsg)
		}
	}
	status += styles.DimTextStyle.Render("(" + itoaStr(len(m.entries)) + " entries)")

	var rows []string
	rows = append(rows, status)

	end := m.scroll + 20
	if end > len(m.entries) {
		end = len(m.entries)
	}

	for i := m.scroll; i < end; i++ {
		e := m.entries[i]
		prefix := "  "
		style := styles.TextStyle
		if i == m.selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			style = styles.HighlightStyle
		}

		action, _ := e["action"].(string)
		actor, _ := e["actor_id"].(string)
		target, _ := e["target_id"].(string)
		ts, _ := e["timestamp"].(float64)

		line := prefix + formatAction(action) + " | " + truncateStr(actor, 12) + " → " + truncateStr(target, 12)
		if ts > 0 {
			line += styles.DimTextStyle.Render("  " + formatTimestamp(int64(ts)))
		}
		rows = append(rows, style.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.JoinVertical(lipgloss.Left, title, content)
}

func formatAction(a string) string {
	switch a {
	case "login":
		return "🔑 LOGIN"
	case "logout":
		return "🚪 LOGOUT"
	case "rent":
		return "📼 RENT"
	case "return":
		return "📀 RETURN"
	case "register":
		return "👤 REGISTER"
	case "promote":
		return "⬆ PROMOTE"
	case "demote":
		return "⬇ DEMOTE"
	case "ban":
		return "🚫 BAN"
	default:
		return a
	}
}

func formatTimestamp(ts int64) string {
	return fmt.Sprintf("%d", ts)
}

func itoaStr(n int) string {
	return fmt.Sprintf("%d", n)
}

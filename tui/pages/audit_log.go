package pages

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type AuditLogRefreshMsg struct{}

const auditPageSize = 20

type AuditLogModel struct {
	table     table.Model
	entries   []map[string]interface{}
	errMsg    string
	VerifyMsg string

	brokenIDX int
	brokenID  string

	paginator paginator.Model
	page      int
	pageSize  int
}

func NewAuditLogModel() *AuditLogModel {
	cols := []table.Column{
		{Title: "TIME", Width: 8},
		{Title: "ACTION", Width: 14},
		{Title: "ACTOR", Width: 14},
		{Title: "TARGET", Width: 14},
		{Title: "HASH", Width: 12},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(auditPageSize),
	)
	t.SetStyles(gruvboxTableStyles())

	p := paginator.New()
	p.PerPage = auditPageSize
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(styles.Green).Bold(true).Render("●")
	p.InactiveDot = lipgloss.NewStyle().Foreground(styles.Grey0).Render("○")
	p.ArabicFormat = "%d / %d"

	return &AuditLogModel{
		table:     t,
		paginator: p,
		pageSize:  auditPageSize,
		page:      0,
		brokenIDX: -1,
	}
}

func (m *AuditLogModel) MarkIntact(total int) {
	m.brokenIDX = -1
	m.brokenID = ""
}

func (m *AuditLogModel) MarkBroken(brokenAt int, brokenID, reason string) {
	m.brokenID = brokenID
	m.refreshPage()
}

func (m *AuditLogModel) SetEntries(entries []map[string]interface{}) {
	sort.SliceStable(entries, func(i, j int) bool {
		ti, _ := entries[i]["timestamp"].(float64)
		tj, _ := entries[j]["timestamp"].(float64)
		return ti > tj
	})
	m.entries = entries
	m.page = 0
	m.refreshPage()
}

func (m *AuditLogModel) refreshPage() {
	if m.brokenID != "" {
		m.brokenIDX = -1
		for i, e := range m.entries {
			if id, _ := e["id"].(string); id == m.brokenID {
				m.brokenIDX = i
				break
			}
		}
	}

	all := m.buildAllRows(m.entries)
	totalPages := (len(all) + m.pageSize - 1) / m.pageSize
	if totalPages < 1 {
		totalPages = 1
	}
	m.paginator.SetTotalPages(totalPages)
	if m.page >= totalPages {
		m.page = totalPages - 1
	}
	if m.page < 0 {
		m.page = 0
	}
	m.paginator.Page = m.page

	start := m.page * m.pageSize
	end := start + m.pageSize
	if end > len(all) {
		end = len(all)
	}
	if start > len(all) {
		start = len(all)
	}
	m.table.SetRows(all[start:end])
}

func (m *AuditLogModel) Update(msg tea.Msg) (*AuditLogModel, tea.Cmd) {
	var cmd tea.Cmd
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "n", "pgdown", "right":
			if m.page < m.paginator.TotalPages-1 {
				m.page++
				m.refreshPage()
			}
			return m, nil
		case "b", "pgup", "left":
			if m.page > 0 {
				m.page--
				m.refreshPage()
			}
			return m, nil
		case "g":
			if m.brokenIDX >= 0 {
				m.page = m.brokenIDX / m.pageSize
				m.refreshPage()
			}
			return m, nil
		case "home":
			m.page = 0
			m.refreshPage()
			return m, nil
		case "end":
			m.page = m.paginator.TotalPages - 1
			m.refreshPage()
			return m, nil
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *AuditLogModel) View(w, h int) string {
	m.table.SetWidth(w)
	m.table.SetHeight(h - 4)
	if m.table.Focused() {
		m.table.Focus()
	} else {
		m.table.Blur()
	}

	header := styles.HeadingStyle.
		Width(w).
		Align(lipgloss.Left).
		Padding(0, 1).
		Render("🔗 AUDIT LOG — Hash Chain Viewer")

	statusLine := styles.DimTextStyle.Render(
		fmt.Sprintf("Press [V] to verify  ·  [N/B] page  ·  %d entries", len(m.entries)),
	)
	if m.VerifyMsg != "" {
		if strings.HasPrefix(m.VerifyMsg, "✅") {
			statusLine = styles.SuccessTextStyle.Render(m.VerifyMsg)
		} else {
			statusLine = styles.ErrorTextStyle.Render(m.VerifyMsg)
		}
	}
	if m.brokenIDX >= 0 {
		statusLine += styles.ErrorTextStyle.Render("  Press [G] to jump to broken entry")
	}
	if m.errMsg != "" {
		statusLine = styles.ErrorTextStyle.Render(m.errMsg)
	}

	if len(m.entries) == 0 {
		empty := styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No audit entries yet")
		return lipgloss.JoinVertical(lipgloss.Left, header, empty)
	}

	paginatorView := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Render(m.paginator.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		paginatorView,
		m.table.View(),
		statusLine,
	)
}

func (m *AuditLogModel) buildAllRows(entries []map[string]interface{}) []table.Row {
	rows := make([]table.Row, 0, len(entries))
	for i, e := range entries {
		action, _ := e["action"].(string)
		actor, _ := e["actor_id"].(string)
		target, _ := e["target_id"].(string)
		ts, _ := e["timestamp"].(float64)
		hash, _ := e["hash"].(string)

		timeStr := ""
		if ts > 0 {
			t := time.Unix(int64(ts), 0)
			timeStr = t.Format("15:04:05")
		}
		actionStr := formatAction(action)
		if m.brokenIDX == i {
			actionStr = "💥 " + actionStr
			timeStr = timeStr + " ⚠"
		}
		rows = append(rows, table.Row{
			timeStr,
			actionStr,
			shortID(actor),
			shortID(target),
			shortHash(hash),
		})
	}
	return rows
}

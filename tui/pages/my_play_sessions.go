package pages

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type playSessionItem struct {
	session models.GameSession
}

func (p playSessionItem) Title() string { return p.session.GameTitle }
func (p playSessionItem) Description() string {
	remaining := p.session.ExpiresAt - time.Now().Unix()
	if remaining < 0 {
		remaining = 0
	}
	mins := remaining / 60
	secs := remaining % 60
	return fmt.Sprintf("🎮 %dm%02ds remaining", mins, secs)
}
func (p playSessionItem) FilterValue() string { return p.session.GameTitle }

type playSessionDelegate struct{}

func newPlaySessionDelegate() playSessionDelegate { return playSessionDelegate{} }

func (d playSessionDelegate) Height() int                             { return 2 }
func (d playSessionDelegate) Spacing() int                            { return 1 }
func (d playSessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d playSessionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	pi, ok := item.(playSessionItem)
	if !ok {
		return
	}
	s := pi.session

	selected := index == m.Index()

	marker := "  "
	titleStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		titleStyle = lipgloss.NewStyle().Foreground(styles.Orange).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}
	remaining := s.ExpiresAt - time.Now().Unix()
	if remaining < 0 {
		remaining = 0
	}
	mins := remaining / 60
	secs := remaining % 60

	var timeColor lipgloss.Color
	switch {
	case remaining == 0:
		timeColor = styles.Grey1
	case remaining < 60:
		timeColor = styles.Red
	case remaining < 180:
		timeColor = styles.Yellow
	default:
		timeColor = styles.Green
	}
	timeStr := lipgloss.NewStyle().Foreground(timeColor).Bold(true).Render(
		fmt.Sprintf("%dm%02ds", mins, secs),
	)

	var status string
	switch {
	case remaining == 0:
		status = lipgloss.NewStyle().Foreground(styles.Grey1).Render("● EXPIRED")
	case s.Status == "paused":
		status = lipgloss.NewStyle().Foreground(styles.Yellow).Render("⏸ PAUSED")
	default:
		status = lipgloss.NewStyle().Foreground(styles.Orange).Render("● PLAYING")
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker,
		"🕹️  ",
		titleStyle.Render(truncateStr(s.GameTitle, 36)),
		"  ",
		status,
	)
	meta := styles.DimTextStyle.Render("  ⏱ ") + timeStr +
		styles.DimTextStyle.Render(fmt.Sprintf("  ·  session %s  ·  started %s",
			truncateStr(s.ID, 12),
			time.Unix(s.StartedAt, 0).Format("15:04"),
		))

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, meta))
}

type PlaySessionsReloadMsg struct{}

type PlayTickMsg time.Time

type MyPlaySessionsModel struct {
	list     list.Model
	sessions []models.GameSession
}

func NewMyPlaySessionsModel() *MyPlaySessionsModel {
	delegate := newPlaySessionDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🎮 MY PLAY SESSIONS"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	return &MyPlaySessionsModel{list: l}
}

func (m *MyPlaySessionsModel) Init() tea.Cmd { return playTick() }

func (m *MyPlaySessionsModel) SetSessions(sessions []models.GameSession) {
	active := make([]models.GameSession, 0, len(sessions))
	for _, s := range sessions {
		if s.Status != "active" {
			continue
		}
		// Skip sessions that the server still reports as "active" but
		// whose expiry has already passed locally; the next refresh
		// will reconcile this.
		if s.ExpiresAt > 0 && s.ExpiresAt <= time.Now().Unix() {
			continue
		}
		active = append(active, s)
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].StartedAt > active[j].StartedAt
	})
	m.sessions = active
	items := make([]list.Item, len(active))
	for i, s := range active {
		items[i] = playSessionItem{session: s}
	}
	m.list.SetItems(items)
}

func (m *MyPlaySessionsModel) HasExpired() bool {
	now := time.Now().Unix()
	for _, s := range m.sessions {
		if s.ExpiresAt > 0 && s.ExpiresAt <= now {
			return true
		}
	}
	return false
}

func (m *MyPlaySessionsModel) Update(msg tea.Msg) (*MyPlaySessionsModel, tea.Cmd) {
	if _, ok := msg.(PlayTickMsg); ok {
		return m, playTick()
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *MyPlaySessionsModel) View(w, h int) string {
	if len(m.sessions) == 0 {
		empty := styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No active play sessions — visit a game in the catalog to start one")
		return empty
	}
	m.list.SetSize(w, h)
	return m.list.View()
}

func playTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return PlayTickMsg(t) })
}

package pages

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type rentalItem struct {
	rental models.RentalResponse
}

func (r rentalItem) Title() string { return r.rental.MovieTitle }
func (r rentalItem) Description() string {
	return r.detailLine()
}
func (r rentalItem) FilterValue() string { return r.rental.MovieTitle + " " + r.rental.MovieFormat }

func (r rentalItem) detailLine() string {
	const minute = int64(60)
	now := time.Now().Unix()
	secs := (r.rental.DueDate - now) / minute

	var status string
	if r.rental.Status == "returned" {
		status = "✓ returned"
	} else if secs < 0 {
		status = fmt.Sprintf("🔴 overdue %d min", -secs)
	} else if secs <= 1 {
		status = fmt.Sprintf("🟡 due soon · %d min", secs)
	} else {
		status = fmt.Sprintf("🟢 active · %d min left", secs)
	}
	return status
}

type rentalDelegate struct{}

func newRentalDelegate() rentalDelegate { return rentalDelegate{} }

func (d rentalDelegate) Height() int  { return 2 }
func (d rentalDelegate) Spacing() int { return 1 }

func (d rentalDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d rentalDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ri, ok := item.(rentalItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	r := ri.rental

	now := time.Now().Unix()
	const minute = int64(60)
	secs := (r.DueDate - now) / minute

	// ── Line 1: title + format + status
	marker := "  "
	titleStyle := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true)
	if selected {
		titleStyle = lipgloss.NewStyle().Foreground(styles.Green).Bold(true)
		marker = styles.HighlightStyle.Render("▸ ")
	}

	titleStr := titleStyle.Render(truncateStr(r.MovieTitle, 36))
	formatStr := styles.FormatBadge(r.MovieFormat)

	var statusGlyph string
	var statusColor lipgloss.Color
	var statusText string
	switch {
	case r.Status == "returned":
		statusGlyph = "✓"
		statusColor = styles.Grey1
		statusText = "returned"
	case secs < 0:
		statusGlyph = "🔴"
		statusColor = styles.Red
		statusText = fmt.Sprintf("overdue %d min", -secs)
	case secs <= 1:
		statusGlyph = "🟡"
		statusColor = styles.Yellow
		statusText = fmt.Sprintf("due soon · %d min", secs)
	default:
		statusGlyph = "🟢"
		statusColor = styles.Green
		statusText = fmt.Sprintf("%d min left", secs)
	}
	statusStr := lipgloss.NewStyle().Foreground(statusColor).Render(fmt.Sprintf("%s %s", statusGlyph, statusText))

	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		marker, titleStr, "  ", formatStr, "  ", statusStr,
	)

	// ── Line 2: metadata
	meta := []string{}

	if r.IsFreeRental {
		meta = append(meta, lipgloss.NewStyle().Foreground(styles.Yellow).Render("🎟️ FREE RENTAL"))
	}

	fee := r.LateFee + r.RewindFee
	if fee > 0 {
		meta = append(meta, lipgloss.NewStyle().Foreground(styles.Orange).Render(fmt.Sprintf("💵 $%.2f fees", fee)))
	}

	if r.Status == "returned" && r.PointsEarned != 0 {
		var ptsColor lipgloss.Color
		var prefix string
		if r.PointsEarned > 0 {
			ptsColor = styles.Orange
			prefix = "+"
		} else {
			ptsColor = styles.Red
		}
		meta = append(meta, lipgloss.NewStyle().Foreground(ptsColor).Bold(true).Render(
			fmt.Sprintf("%s%d🍿", prefix, r.PointsEarned),
		))
	}

	if r.Status != "returned" && r.RentedAt > 0 {
		meta = append(meta, styles.DimTextStyle.Render(fmt.Sprintf("rented %s", time.Unix(r.RentedAt, 0).Format("15:04"))))
	}

	if len(meta) == 0 {
		meta = []string{styles.DimTextStyle.Render("no extra fees")}
	}

	metaLine := styles.DimTextStyle.Render("  " + meta[0])
	for _, m := range meta[1:] {
		metaLine += styles.DimTextStyle.Render("  ·  ") + m
	}

	io.WriteString(w, lipgloss.JoinVertical(lipgloss.Left, line1, metaLine))
}

type ReturnRequestMsg struct{ RentalID string }
type ExtendRentalMsg struct{ RentalID string }
type RentalsReloadMsg struct{}

type MyRentalsModel struct {
	list       list.Model
	rentals    []models.RentalResponse
	selectedID string
	Status     string
}

func NewMyRentalsModel() *MyRentalsModel {
	delegate := newRentalDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "📼 MY RENTALS"
	l.Styles = gruvboxListStyles()
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	enableListPagination(&l)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	return &MyRentalsModel{list: l}
}

func (m *MyRentalsModel) SetRentals(rs []models.RentalResponse) {
	m.rentals = rs
	items := make([]list.Item, len(rs))
	for i, r := range rs {
		items[i] = rentalItem{rental: r}
	}
	m.list.SetItems(items)
	if len(rs) > 0 {
		m.list.Select(0)
		m.selectedID = rs[0].ID
	} else {
		m.selectedID = ""
	}
}

func (m *MyRentalsModel) SelectedRental() *models.RentalResponse {
	if i, ok := m.list.SelectedItem().(rentalItem); ok {
		return &i.rental
	}
	return nil
}

func (m *MyRentalsModel) Update(msg tea.Msg) (*MyRentalsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if sel, ok := m.list.SelectedItem().(rentalItem); ok {
		m.selectedID = sel.rental.ID
	}
	return m, cmd
}

func (m *MyRentalsModel) View(w, h int) string {
	m.list.SetSize(w, h)
	if len(m.rentals) == 0 {
		empty := styles.DimTextStyle.
			Width(w).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No rental history")
		return empty
	}

	view := m.list.View()

	if m.Status != "" {
		view += "\n" + styles.SuccessTextStyle.Render(m.Status)
	}
	return view
}

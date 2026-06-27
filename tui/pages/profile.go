package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type ProfileModel struct {
	User      *models.UserResponse
	Stats     *RentalStats
	StatusMsg string
}

type RentalStats struct {
	Total   int
	LateFee float64
	Rewind  float64
}

func NewProfileModel(u *models.UserResponse) *ProfileModel { return &ProfileModel{User: u} }
func (m *ProfileModel) SetStats(s *RentalStats)            { m.Stats = s }

func (m *ProfileModel) View(w, h int) string {
	if m.User == nil {
		return ""
	}
	u := m.User

	tierColor := styles.Green
	if c, ok := styles.TierColors[u.TierName]; ok {
		tierColor = c
	}

	cardW := 52
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tierColor).
		Padding(1, 2).
		Width(cardW)

	header := m.headerView(tierColor, cardW-4)
	body := m.bodyView(tierColor, cardW-4)

	cardBody := lipgloss.JoinVertical(lipgloss.Left,
		header,
		strings.Repeat("─", cardW-4),
		body,
	)

	title := styles.HeadingStyle.
		Width(w).
		Align(lipgloss.Left).
		Padding(0, 1).
		Render("👤 MEMBER PROFILE")

	column := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		card.Render(cardBody),
	)

	if m.StatusMsg != "" {
		column += "\n" + styles.SuccessTextStyle.Render("  ✓ "+m.StatusMsg)
	}

	return column
}

func (m *ProfileModel) headerView(tierColor lipgloss.Color, w int) string {
	u := m.User

	badge := styles.TierBadgeStyle(u.TierName).Render(" ★ " + u.TierName + " ★ ")

	name := lipgloss.NewStyle().
		Foreground(styles.FG1).
		Bold(true).
		Render(u.Username)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, name, "  ", badge)

	var row2 string
	if u.CreatedAt > 0 {
		t := time.Unix(u.CreatedAt, 0)
		row2 = styles.DimTextStyle.Render("Member since " + t.Format("Jan 2006"))
	}

	flags := m.statusFlags()
	row3 := flags

	parts := []string{row1}
	if row2 != "" {
		parts = append(parts, row2)
	}
	if row3 != "" {
		parts = append(parts, row3)
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *ProfileModel) statusFlags() string {
	u := m.User
	var flags []string

	if u.Banned {
		flags = append(flags, lipgloss.NewStyle().Foreground(styles.Red).Bold(true).Render("🚫 BANNED"))
	}
	if u.TOTPEnabled {
		flags = append(flags, lipgloss.NewStyle().Foreground(styles.Blue).Bold(true).Render("🔒 2FA"))
	}

	if len(flags) == 0 {
		return ""
	}
	return lipgloss.NewStyle().Foreground(styles.Grey1).Render(strings.Join(flags, "  "))
}

func (m *ProfileModel) bodyView(tierColor lipgloss.Color, w int) string {
	u := m.User

	grid := lipgloss.JoinHorizontal(lipgloss.Top,
		m.statCol("BALANCE", fmt.Sprintf("$%.2f", u.Balance), styles.Yellow, 18),
		m.statCol("POPCORN", fmt.Sprintf("%d 🍿", u.PopcornPoints), styles.Orange, 16),
		m.statCol("FREE", fmt.Sprintf("%d/%d", u.FreeRentals, models.TierByName(u.Subscription).FreeRentals), styles.Aqua, 12),
	)

	rentalBar := m.rentalProgressBar()

	parts := []string{grid, ""}
	if rentalBar != "" {
		parts = append(parts, rentalBar, "")
	}
	if m.Stats != nil && m.Stats.Total > 0 {
		parts = append(parts, m.statsLine())
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *ProfileModel) statCol(label, value string, color lipgloss.Color, w int) string {
	lbl := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render(label)

	val := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Width(w).
		Align(lipgloss.Center).
		Render(value)

	return lipgloss.JoinVertical(lipgloss.Center, lbl, val)
}

func (m *ProfileModel) rentalProgressBar() string {
	u := m.User
	if u.MaxRentals <= 0 {
		return ""
	}

	barWidth := 30
	filled := int(float64(u.RentalCount) / float64(u.MaxRentals) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	fillColor := styles.Green
	switch {
	case filled >= barWidth:
		fillColor = styles.Red
	case filled >= barWidth*4/5:
		fillColor = styles.Yellow
	}

	bar := strings.Repeat("▰", filled) + strings.Repeat("▱", barWidth-filled)
	coloredBar := lipgloss.NewStyle().Foreground(fillColor).Render(bar)

	label := lipgloss.NewStyle().Foreground(styles.Grey1).Render("RENTALS")
	count := lipgloss.NewStyle().Foreground(styles.FG1).Bold(true).Render(
		fmt.Sprintf("  %d/%d", u.RentalCount, u.MaxRentals),
	)

	return label + count + "  " + coloredBar
}

func (m *ProfileModel) statsLine() string {
	s := m.Stats
	parts := []string{
		lipgloss.NewStyle().Foreground(styles.Grey1).Render("LIFETIME "),
		lipgloss.NewStyle().Foreground(styles.FG0).Render(fmt.Sprintf("%d rentals", s.Total)),
	}
	if s.LateFee > 0 {
		parts = append(parts,
			lipgloss.NewStyle().Foreground(styles.Orange).Render(
				fmt.Sprintf("💵 $%.2f late", s.LateFee)))
	}
	if s.Rewind > 0 {
		parts = append(parts,
			lipgloss.NewStyle().Foreground(styles.Orange).Render(
				fmt.Sprintf("🔄 $%.2f rewind", s.Rewind)))
	}
	return strings.Join(parts, "  ")
}

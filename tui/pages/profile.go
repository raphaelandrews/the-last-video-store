package pages

import (
	"fmt"
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
	bd := styles.Green
	if c, ok := styles.TierColors[m.User.TierName]; ok {
		bd = c
	}

	headerBlock := lipgloss.NewStyle().
		Foreground(styles.Green).
		Bold(true).
		Width(46).
		Align(lipgloss.Center).
		Render("─── THE LAST VIDEO STORE ───")

	subHeader := lipgloss.NewStyle().
		Foreground(styles.Grey1).
		Width(46).
		Align(lipgloss.Center).
		Render("MEMBERSHIP CARD")

	badge := styles.TierBadgeStyle(m.User.TierName).Render(" ★ " + m.User.TierName + " ★ ")
	stats := ""
	if m.Stats != nil {
		stats = fmt.Sprintf("\n  📀 Active: %d/%d\n  🍿 Popcorn: %d\n  ⏱ Total: %d\n  💵 Late fees: $%.2f\n  🔄 Rewind fees: $%.2f",
			m.User.RentalCount, m.User.MaxRentals, m.User.PopcornPoints,
			m.Stats.Total, m.Stats.LateFee, m.Stats.Rewind)
	}

	card := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(bd).
		Padding(1, 4).
		Width(46).
		Align(lipgloss.Center)

	cardBody := lipgloss.JoinVertical(lipgloss.Center,
		headerBlock,
		subHeader,
		"",
		styles.TextStyle.Render("  Username: "+m.User.Username),
		styles.TextStyle.Render("  Plan: "+badge),
		styles.DimTextStyle.Render("  Member since: "+fmt.Sprintf("%d", m.User.CreatedAt)),
		stats,
	)

	title := styles.HeadingStyle.
		Width(w).
		Align(lipgloss.Left).
		Padding(0, 1).
		Render("MEMBER PROFILE")

	// Build the column: title at the top-left, card below it. The whole
	// column is left-aligned (placed horizontally at the left edge) but
	// the card content is centered within its 46-char border.
	result := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		card.Render(cardBody),
	)

	topUpInfo := ""
	if m.User.LastTopUpAt > 0 {
		elapsed := time.Now().Unix() - m.User.LastTopUpAt
		if elapsed < 30 {
			cooldown := 30 - int(elapsed)
			topUpInfo = fmt.Sprintf("💰 Top-up cooldown: %ds — press [$]", cooldown)
		} else {
			topUpInfo = "💰 Top-up available — press [$]"
		}
	} else {
		topUpInfo = "💰 Top-up available — press [$]"
	}
	if topUpInfo != "" {
		result += "\n" + styles.DimTextStyle.Render(topUpInfo)
	}

	if m.StatusMsg != "" {
		result += "\n" + styles.SuccessTextStyle.Render(m.StatusMsg)
	}
	return result
}

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
	bd := styles.SkyBlue
	if c, ok := styles.TierColors[m.User.TierName]; ok {
		bd = c
	}
	card := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(bd).
		Background(styles.BgWhite).
		Padding(1, 4).Width(42).Align(lipgloss.Center)

	badge := styles.TierBadgeStyle(m.User.TierName).Render(" ★ " + m.User.TierName + " ★ ")
	stats := ""
	if m.Stats != nil {
		stats = fmt.Sprintf("\n📀 Active: %d/%d\n🍿 Popcorn: %d\n⏱ Total: %d\n💵 Late fees: $%.2f\n🔄 Rewind fees: $%.2f",
			m.User.RentalCount, m.User.MaxRentals, m.User.PopcornPoints,
			m.Stats.Total, m.Stats.LateFee, m.Stats.Rewind)
	}
	inner := card.Render(lipgloss.JoinVertical(lipgloss.Center,
		"THE LAST VIDEO STORE", "",
		"Username: "+m.User.Username,
		"Plan: "+badge,
		"Member since: "+fmt.Sprintf("%d", m.User.CreatedAt),
		stats,
	))

	topUpInfo := ""
	if m.User.TopUpCount > 0 || m.User.LastTopUpAt > 0 {
		if m.User.TopUpCount == 0 {
			cooldown := ""
			elapsed := time.Now().Unix() - m.User.LastTopUpAt
			if elapsed < 30 {
				cooldown = fmt.Sprintf(" (cooldown: %ds)", 30-int(elapsed))
			}
			topUpInfo = fmt.Sprintf("\n💰 Top-up available in%s — press [$]", cooldown)
		} else {
			topUpInfo = fmt.Sprintf("\n💰 Top-up already used")
		}
	}
	if topUpInfo != "" {
		inner += styles.DimTextStyle.Render(topUpInfo)
	}
	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("MEMBER PROFILE")
	result := lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, title, inner))
	if m.StatusMsg != "" {
		result += "\n" + styles.SuccessTextStyle.Render(m.StatusMsg)
	}
	return result
}

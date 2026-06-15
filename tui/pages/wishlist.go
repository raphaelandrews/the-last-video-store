package pages

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/styles"
)

type WishlistRemoveMsg struct{ MovieID string }
type WishlistReloadMsg struct{}

type WishlistItem struct {
	ID        string `json:"id"`
	MovieID   string `json:"movie_id"`
	Title     string `json:"title"`
	Available bool   `json:"available"`
	AddedAt   int64  `json:"added_at"`
}

type WishlistModel struct {
	Items    []WishlistItem
	Selected int
	Loading  bool
	Status   string
}

func NewWishlistModel() *WishlistModel { return &WishlistModel{Selected: -1, Loading: true} }

func (m *WishlistModel) SetItems(items []WishlistItem) {
	m.Items = items
	m.Loading = false
	m.Status = fmt.Sprintf("%d items", len(items))
	if len(items) > 0 && m.Selected < 0 {
		m.Selected = 0
	}
	if len(items) == 0 {
		m.Selected = -1
	}
}

func (m *WishlistModel) MoveUp() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Items) - 1
	}
}

func (m *WishlistModel) MoveDown() {
	if len(m.Items) == 0 {
		return
	}
	m.Selected++
	if m.Selected >= len(m.Items) {
		m.Selected = 0
	}
}

func (m *WishlistModel) SelectedItem() *WishlistItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return &m.Items[m.Selected]
	}
	return nil
}

func (m *WishlistModel) RemoveSelected() {
	if m.Selected < 0 || m.Selected >= len(m.Items) {
		return
	}
	m.Items = append(m.Items[:m.Selected], m.Items[m.Selected+1:]...)
	if m.Selected >= len(m.Items) {
		m.Selected = len(m.Items) - 1
	}
}

func (m *WishlistModel) View(w, h int) string {
	if m.Loading {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("Loading wishlist..."))
	}
	if len(m.Items) == 0 {
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
			styles.TextStyle.Render("Wishlist is empty — browse and press [W] to add movies"))
	}

	title := styles.HeadingStyle.Width(w).Align(lipgloss.Center).Render("📋 MY WISHLIST")
	var rows []string
	for i, item := range m.Items {
		prefix := "  "
		st := styles.TextStyle
		if i == m.Selected {
			prefix = styles.HighlightStyle.Render("▸ ")
			st = styles.HighlightStyle
		}
		avail := styles.SuccessTextStyle.Render("✓ available")
		if !item.Available {
			avail = styles.ErrorTextStyle.Render("✗ out")
		}
		line := fmt.Sprintf("%s%-30s %s", prefix, item.Title, avail)
		rows = append(rows, st.Render(line))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, rows...)...)
	if m.Status != "" {
		content += "\n" + styles.DimTextStyle.Render(m.Status)
	}
	return content
}

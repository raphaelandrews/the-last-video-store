package tui

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadWishlist() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/wishlist")
		if resp == nil {
			return pages.ErrorMsg{Message: "failed to load wishlist"}
		}
		defer resp.Body.Close()
		var items []pages.WishlistItem
		json.NewDecoder(resp.Body).Decode(&items)
		return loadWishlistMsg{items: items}
	}
}

func (m *Model) doRemoveFromWishlist(movieID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.apiDelete("/api/v1/wishlist/" + movieID)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		m.wishlist.RemoveSelected()
		m.wishlist.Status = "Removed from wishlist"
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			m.wishlist.Status = e.Error
		}
		return nil
	}
}

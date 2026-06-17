package tui

import (
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadWishlist() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/wishlist", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var items []pages.WishlistItem
		json.NewDecoder(resp.Body).Decode(&items)
		return loadWishlistMsg{items: items}
	}
}

func (m *Model) doRemoveFromWishlist(movieID string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("DELETE", m.baseURL+"/api/v1/wishlist/"+movieID, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
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

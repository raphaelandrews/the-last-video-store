package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) detailKey(k string) tea.Cmd {
	switch k {
	case "enter":
		rv := m.detail.SelectedRelated()
		if rv != nil {
			m.detail = pages.NewMovieDetailModel(rv)
			m.setDetailContext()
		} else if m.detail != nil && !m.detail.Rented && m.detail.Movie.Available {
			if m.detail.FreeRentals > 0 && m.detail.Balance >= models.MovieCost(m.detail.Movie.RentalPrice, m.detail.Movie.Format) {
				m.detail.Choosing = true
				m.detail.ErrMsg = ""
			} else if m.detail.FreeRentals > 0 {
				m.detail.UseTicket = true
				return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
			} else if m.detail.Balance >= models.MovieCost(m.detail.Movie.RentalPrice, m.detail.Movie.Format) {
				m.detail.UseTicket = false
				return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
			} else {
				m.detail.ErrMsg = "💵 Insufficient funds — upgrade tier or add funds"
			}
		} else if m.detail != nil && !m.detail.Rented && !m.detail.Movie.Available {
			m.detail.ErrMsg = "🔴 No copies available — press [W] to join the waitlist"
		}
	case "t":
		if m.detail != nil && m.detail.Choosing && m.detail.FreeRentals > 0 {
			m.detail.Choosing = false
			m.detail.UseTicket = true
			return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
		}
	case "m":
		if m.detail != nil && m.detail.Choosing && m.detail.Balance >= models.MovieCost(m.detail.Movie.RentalPrice, m.detail.Movie.Format) {
			m.detail.Choosing = false
			m.detail.UseTicket = false
			return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
		}
	case "esc":
		if m.detail != nil && m.detail.Choosing {
			m.detail.Choosing = false
			m.detail.ErrMsg = ""
		}
	case "down", "j":
		if m.detail != nil {
			m.detail.MoveRelatedDown()
		}
	case "up", "k":
		if m.detail != nil {
			m.detail.MoveRelatedUp()
		}
	case "w":
		if m.detail != nil {
			return m.doAddToWishlist(m.detail.Movie.ID, true)
		}
	case "f5":
		return m.loadMovies(m.browse.Page, m.browse.Genre)
	}
	return nil
}

func (m *Model) gameDetailKey(k string) tea.Cmd {
	switch k {
	case "r":
		if m.gameDetail != nil && m.gameDetail.Game != nil && !m.gameDetail.Playing && !m.gameDetail.ChoosingTime && m.gameDetail.Game.Available && m.gameDetail.Game.RentalPrice > 0 {
			return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.gameDetail.Game.ID} }
		}
	case "p":
		if m.gameDetail != nil && m.gameDetail.Game != nil && !m.gameDetail.Playing && !m.gameDetail.ChoosingTime && m.gameDetail.Game.Available && m.gameDetail.Game.PlayPrice > 0 {
			m.gameDetail.ChoosingTime = true
		}
	case "esc":
		if m.gameDetail != nil && m.gameDetail.ChoosingTime {
			m.gameDetail.ChoosingTime = false
		}
	case "e":
		if m.gameDetail != nil && m.gameDetail.Playing && m.gameDetail.Session != nil {
			return m.doGamePlayEnd(m.gameDetail.Session.ID)
		}
	case "1", "2", "3", "4", "5":
		if m.gameDetail != nil && m.gameDetail.ChoosingTime {
			duration := int(k[0] - '0')
			m.gameDetail.ChoosingTime = false
			return m.doGamePlayStart(m.gameDetail.Game.ID, m.gameDetail.Game.Title, duration)
		}
	case "down", "j":
		m.gameDetail.MoveRelatedDown()
	case "up", "k":
		m.gameDetail.MoveRelatedUp()
	}
	return nil
}

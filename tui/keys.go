package tui

import (
	"github.com/thelastvideostore/tui/pages"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) pageKey(msg tea.KeyMsg) tea.Cmd {
	k := msg.String()

	switch m.screen {
	case scrTOTP:
		if k == "enter" && len(m.totpCode) == 6 {
			return m.doSubmitTOTP(m.tempToken, m.totpCode)
		}
		if k == "backspace" && len(m.totpCode) > 0 {
			m.totpCode = m.totpCode[:len(m.totpCode)-1]
		}
		if len(k) == 1 && k[0] >= '0' && k[0] <= '9' && len(m.totpCode) < 6 {
			m.totpCode += k
		}
	case scrBrowse:
		switch k {
		case "down", "j":
			m.browse.MoveDown()
		case "up", "k":
			m.browse.MoveUp()
		case "enter":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				m.detail = pages.NewMovieDetailModel(mv)
				m.screen = scrDetail
			}
		case "d":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				m.detail = pages.NewMovieDetailModel(mv)
				m.screen = scrDetail
			}
		case "r":
			return func() tea.Msg { return pages.NavigateMsg{Page: "rentals"} }
		case "p":
			return func() tea.Msg { return pages.NavigateMsg{Page: "profile"} }
		case "v":
			m.wishlist = pages.NewWishlistModel()
			m.screen = scrWishlist
			return m.loadWishlist()
		case "/":
			m.searching = true
			m.searchBar.Focus()
		case "s":
			if m.browse.Mode != pages.ModeStaffPicks {
				m.browse.Mode = pages.ModeStaffPicks
				m.browse.Selected = -1
				m.browse.Loading = true
				return m.loadStaffPicks()
			}
		case "l":
			if m.browse.Mode != pages.ModeLastChance {
				m.browse.Mode = pages.ModeLastChance
				m.browse.Selected = -1
				m.browse.Loading = true
				return m.loadLastChance()
			}
		case "a":
			if m.browse.Mode != pages.ModeAll {
				m.browse.Mode = pages.ModeAll
				m.browse.Selected = -1
				m.browse.Loading = true
				return m.loadMovies(1)
			}
		case "n":
			if m.browse.HasNextPage() {
				m.browse.Selected = -1
				m.browse.Loading = true
				return m.loadMovies(m.browse.Page + 1)
			}
		case "b":
			if m.browse.HasPrevPage() {
				m.browse.Selected = -1
				m.browse.Loading = true
				return m.loadMovies(m.browse.Page - 1)
			}
		case "ctrl+a":
			if m.userResp != nil && m.userResp.TierName != "Couch Potato" && m.userResp.TierName != "Matinee Fan" {
				m.adminMovies = pages.NewAdminMoviesModel()
				m.screen = scrAdminMovies
				return m.loadAdminMovies()
			}
		case "ctrl+u":
			if m.userResp != nil && m.userResp.TierName != "Couch Potato" && m.userResp.TierName != "Matinee Fan" && m.userResp.TierName != "Gold Member" {
				m.adminUsers = pages.NewAdminUsersModel()
				m.screen = scrAdminUsers
				return m.loadAdminUsers()
			}
		case "ctrl+g":
			if m.userResp != nil && m.userResp.TierName != "Couch Potato" && m.userResp.TierName != "Matinee Fan" {
				m.auditLog = pages.NewAuditLogModel()
				m.screen = scrAuditLog
				return m.loadAuditLog()
			}
		case "f5":
			return m.loadMovies(m.browse.Page)
		}
	case scrDetail:
		switch k {
		case "enter":
			if m.detail != nil && !m.detail.Rented && m.detail.Movie.Available {
				return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
			} else if m.detail != nil && !m.detail.Rented && !m.detail.Movie.Available {
				m.detail.ErrMsg = "🔴 No copies available — press [W] to join the waitlist"
			}
		case "w":
			if m.detail != nil {
				return m.doAddToWishlist(m.detail.Movie.ID, true)
			}
		case "f5":
			return m.loadMovies(m.browse.Page)
		}
	case scrRentals:
		switch k {
		case "down", "j":
			m.rentals.MoveDown()
		case "up", "k":
			m.rentals.MoveUp()
		case "enter":
			r := m.rentals.SelectedRental()
			if r != nil && r.Status != "returned" {
				return func() tea.Msg { return pages.ReturnRequestMsg{RentalID: r.ID} }
			}
		case "e":
			r := m.rentals.SelectedRental()
			if r != nil && r.Status != "returned" && r.Status != "overdue" {
				return func() tea.Msg { return pages.ExtendRentalMsg{RentalID: r.ID} }
			}
		}
	case scrProfile:
		if k == "l" {
			return func() tea.Msg { return pages.NavigateMsg{Page: "login"} }
		}
		if k == "m" {
			pts := 0
			if m.userResp != nil {
				pts = m.userResp.PopcornPoints
			}
			m.merch = pages.NewMerchModel(pts)
			m.screen = scrMerch
			return m.loadMerch()
		}
		if k == "i" {
			m.inventory = pages.NewInventoryModel()
			m.screen = scrInventory
			return m.loadInventory()
		}
	case scrWishlist:
		switch k {
		case "down", "j":
			m.wishlist.MoveDown()
		case "up", "k":
			m.wishlist.MoveUp()
		case "enter":
			item := m.wishlist.SelectedItem()
			if item != nil {
				m.browse.Status = "Use [W] on browse to rent a waitlisted title"
			}
		case "d", "delete":
			item := m.wishlist.SelectedItem()
			if item != nil {
				return m.doRemoveFromWishlist(item.MovieID)
			}
		}
	case scrAdminMovies:
		switch k {
		case "down", "j":
			m.adminMovies.MoveDown()
		case "up", "k":
			m.adminMovies.MoveUp()
		}
	case scrAdminUsers:
		switch k {
		case "down", "j":
			m.adminUsers.MoveDown()
		case "up", "k":
			m.adminUsers.MoveUp()
		}
	case scrAuditLog:
		switch k {
		case "down", "j":
			m.auditLog.MoveDown()
		case "up", "k":
			m.auditLog.MoveUp()
		}
	case scrMerch:
		switch k {
		case "down", "j":
			m.merch.MoveDown()
		case "up", "k":
			m.merch.MoveUp()
		case "enter":
			item := m.merch.SelectedItem()
			if item != nil && item.Stock > 0 && m.userResp.PopcornPoints >= item.PointsCost {
				return func() tea.Msg { return pages.MerchRedeemMsg{ItemID: item.ID} }
			}
		}
	}
	return nil
}

func (m *Model) searchKey(msg tea.KeyMsg) tea.Cmd {
	k := msg.String()
	switch k {
	case "esc":
		m.searching = false
		m.searchBar.Blur()
		return nil
	case "enter":
		mv := m.searchBar.SelectedMovie()
		if mv != nil {
			m.searching = false
			m.searchBar.Blur()
			m.detail = pages.NewMovieDetailModel(mv)
			m.screen = scrDetail
		}
		return nil
	case "up", "k":
		m.searchBar.MoveSelection(-1)
	case "down", "j":
		m.searchBar.MoveSelection(1)
	default:
		m.searchBar.Update(msg)
		return m.doSearch(m.searchBar.Value())
	}
	return nil
}

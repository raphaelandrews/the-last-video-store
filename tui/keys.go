package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
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
		if len(k) == 1 && k >= "0" && k <= "9" && len(m.totpCode) < 6 {
			m.totpCode += k
		}

	case scrBrowse:
		return m.browseKey(k)

	case scrDetail:
		return m.detailKey(k)

	case scrGameDetail:
		return m.gameDetailKey(k)

	case scrRentals:
		switch k {
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
		case "p":
			return func() tea.Msg { return pages.NavigateMsg{Page: "play_sessions"} }
		}

	case scrProfile:
		return m.navKey(k)

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

	case scrMerch:
		switch k {
		case "enter":
			item := m.merch.SelectedItem()
			if item != nil && item.Stock > 0 && m.userResp.PopcornPoints >= item.PointsCost {
				return func() tea.Msg { return pages.MerchRedeemMsg{ItemID: item.ID} }
			}
		}

	case scrTierShop:
		switch k {
		case "enter":
			tier := m.tierShop.SelectedTier()
			if tier != nil && tier.Name != m.tierShop.Current && (tier.Price == 0 || m.userResp.Balance >= tier.Price) {
				return tea.Sequence(m.doPurchaseTier(tier.Name), m.doRefreshMe())
			}
		}

	case scrMovieForm:
		if k == "esc" {
			m.screen = scrAdminMovies
			m.movieForm = nil
		}

	case scrAdminMovies:
		return m.adminKey(k)

	case scrAdminUsers:
		return m.adminUsersKey(k)

	case scrAuditLog:
		return m.auditLogKey(k)

	case scrAccessDenied:
		if k == "q" {
			m.screen = scrBrowse
		}

	case scrSnackBarMenu:
		return m.snackbarMenuKey(k)

	case scrSnackBarManage:
		return m.snackbarManageKey(k)

	case scrSnackBarOrders:
		_ = k

	case scrGameSessions:
		_ = k
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
			m.setDetailContext()
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

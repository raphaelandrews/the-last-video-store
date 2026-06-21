package tui

import (
	"github.com/thelastvideostore/internal/models"
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
				if mv.MediaType == "game" {
					m.gameDetail = pages.NewGameDetailModel(mv)
					bal := 0.0
					if m.userResp != nil {
						bal = m.userResp.Balance
					}
					m.gameDetail.Balance = bal
					m.screen = scrGameDetail
				} else {
					m.detail = pages.NewMovieDetailModel(mv)
					m.setDetailContext()
					m.screen = scrDetail
				}
			}
		case "d":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				if mv.MediaType == "game" {
					m.gameDetail = pages.NewGameDetailModel(mv)
					bal := 0.0
					if m.userResp != nil {
						bal = m.userResp.Balance
					}
					m.gameDetail.Balance = bal
					m.screen = scrGameDetail
				} else {
					m.detail = pages.NewMovieDetailModel(mv)
					m.setDetailContext()
					m.screen = scrDetail
				}
			}
		case "r":
			return func() tea.Msg { return pages.NavigateMsg{Page: "rentals"} }
		case "c":
			bal := 0.0
			if m.userResp != nil {
				bal = m.userResp.Balance
			}
			m.snackBarMenu = pages.NewSnackBarMenuModel(bal)
			m.snackBarMenu.SetItems(nil)
			m.screen = scrSnackBarMenu
			return m.loadSnackBarMenu()
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
				return m.loadStaffPicks()
			}
		case "l":
			if m.browse.Mode != pages.ModeLastChance {
				m.browse.Mode = pages.ModeLastChance
				m.browse.Selected = -1
				return m.loadLastChance()
			}
		case "a":
			if m.browse.Mode != pages.ModeAll {
				m.browse.Mode = pages.ModeAll
				m.browse.Selected = -1
				m.tabs.SetActive(0)
				m.genreTabs.SetActive(0)
				m.browse.Genre = ""
				m.browse.MediaType = "movie"
				return m.loadMovies(1, "")
			}
		case "[":
			if m.browse.Mode == pages.ModeAll {
				m.tabs.Prev()
				return m.applyMediaTypeFilter()
			}
		case "]":
			if m.browse.Mode == pages.ModeAll {
				m.tabs.Next()
				return m.applyMediaTypeFilter()
			}
		case ",":
			if m.browse.Mode == pages.ModeAll && m.tabs.ActiveTab() != "🍿 SnackBar" && m.genreTabs.ActiveTab() != "" {
				m.genreTabs.Prev()
				m.browse.Genre = m.genreTabs.ActiveTab()
				if m.browse.Genre == "ALL" {
					m.browse.Genre = ""
				}
				m.browse.Selected = -1
				return m.loadMovies(1, m.browse.Genre)
			}
		case ".":
			if m.browse.Mode == pages.ModeAll && m.tabs.ActiveTab() != "🍿 SnackBar" && m.genreTabs.ActiveTab() != "" {
				m.genreTabs.Next()
				m.browse.Genre = m.genreTabs.ActiveTab()
				if m.browse.Genre == "ALL" {
					m.browse.Genre = ""
				}
				m.browse.Selected = -1
				return m.loadMovies(1, m.browse.Genre)
			}
		case "n":
			if m.browse.HasNextPage() {
				m.browse.Selected = -1
				return m.loadMovies(m.browse.Page+1, m.browse.Genre)
			}
		case "b":
			if m.browse.HasPrevPage() {
				m.browse.Selected = -1
				return m.loadMovies(m.browse.Page-1, m.browse.Genre)
			}
		case "ctrl+a":
			m.adminMovies = pages.NewAdminMoviesModel()
			m.screen = scrAdminMovies
			return m.loadAdminMovies(1)
		case "ctrl+u":
			m.adminUsers = pages.NewAdminUsersModel()
			m.screen = scrAdminUsers
			return m.loadAdminUsers()
		case "ctrl+g":
			m.auditLog = pages.NewAuditLogModel()
			m.screen = scrAuditLog
			return m.loadAuditLog()
		case "f5":
			return m.loadMovies(m.browse.Page, m.browse.Genre)
		}
	case scrDetail:
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
		if k == "b" {
			bal := 0.0
			if m.userResp != nil {
				bal = m.userResp.Balance
			}
			m.snackBarMenu = pages.NewSnackBarMenuModel(bal)
			m.snackBarMenu.SetItems(nil)
			m.screen = scrSnackBarMenu
			return m.loadSnackBarMenu()
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
		if k == "t" {
			sub := "wood"
			bal := 0.0
			if m.userResp != nil {
				sub = m.userResp.Subscription
				bal = m.userResp.Balance
			}
			m.tierShop = pages.NewTierShopModel(sub, bal)
			m.screen = scrTierShop
		}
		if k == "2" {
			return m.doProfileTOTP()
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
		case "a":
			m.movieForm = pages.NewMovieFormModel()
			m.screen = scrMovieForm
		case "enter":
			mv := m.adminMovies.SelectedMovie()
			if mv != nil {
				m.movieForm = pages.NewMovieEditFormModel(mv)
				m.screen = scrMovieForm
			}
		case "s":
			mv := m.adminMovies.SelectedMovie()
			if mv != nil {
				return m.doToggleStaffPick(mv.ID, mv.IsStaffPick)
			}
		case "d":
			mv := m.adminMovies.SelectedMovie()
			if mv != nil {
				return m.doDeleteMovie(mv.ID)
			}
		case "n":
			if m.adminMovies.HasNextPage() {
				return m.loadAdminMovies(m.adminMovies.Page + 1)
			}
		case "b":
			if m.adminMovies.HasPrevPage() {
				return m.loadAdminMovies(m.adminMovies.Page - 1)
			}
		}
	case scrAdminUsers:
		switch k {
		case "down", "j":
			m.adminUsers.MoveDown()
		case "up", "k":
			m.adminUsers.MoveUp()
		case "p":
			u := m.adminUsers.SelectedUser()
			if u != nil {
				return m.doUpdateUser(u.ID, "promote")
			}
		case "d":
			u := m.adminUsers.SelectedUser()
			if u != nil {
				return m.doUpdateUser(u.ID, "demote")
			}
		case "b":
			u := m.adminUsers.SelectedUser()
			if u != nil {
				return m.doUpdateUser(u.ID, "ban")
			}
		case "t":
			u := m.adminUsers.SelectedUser()
			if u != nil {
				return m.doTOTPToggle(u.ID)
			}
		}
	case scrAuditLog:
		switch k {
		case "down", "j":
			m.auditLog.MoveDown()
		case "up", "k":
			m.auditLog.MoveUp()
		case "v":
			return m.doVerifyAuditChain()
		}
	case scrAccessDenied:
		if k == "q" {
			m.screen = scrBrowse
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
	case scrTierShop:
		switch k {
		case "down", "j":
			m.tierShop.MoveDown()
		case "up", "k":
			m.tierShop.MoveUp()
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
	case scrSnackBarMenu:
		switch k {
		case "down", "j":
			m.snackBarMenu.MoveDown()
		case "up", "k":
			m.snackBarMenu.MoveUp()
		case "enter":
			item := m.snackBarMenu.SelectedItem()
			if item != nil && item.Stock > 0 && m.userResp != nil && m.userResp.Balance >= item.Price {
				return func() tea.Msg { return pages.SnackBarOrderMsg{ItemID: item.ID} }
			}
		case "o":
			m.snackBarOrders = pages.NewSnackBarOrdersModel()
			m.screen = scrSnackBarOrders
			return m.loadSnackBarOrders()
		case "m":
			m.snackBarManage = pages.NewSnackBarManageModel(m.userResp.Tier)
			m.screen = scrSnackBarManage
			return m.loadSnackBarManage()
		}
	case scrSnackBarManage:
		switch k {
		case "down", "j":
			m.snackBarManage.MoveDown()
		case "up", "k":
			m.snackBarManage.MoveUp()
		case "r":
			item := m.snackBarManage.SelectedItem()
			if item != nil {
				return m.doSnackBarRestock(item.ID)
			}
		}
	case scrSnackBarOrders:
		_ = k
	case scrGameDetail:
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

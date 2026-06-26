package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) adminKey(k string) tea.Cmd {
	switch k {
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
	case "down", "j":
		m.adminMovies.MoveDown()
	case "up", "k":
		m.adminMovies.MoveUp()
	}
	return nil
}

func (m *Model) adminUsersKey(k string) tea.Cmd {
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
	return nil
}

func (m *Model) auditLogKey(k string) tea.Cmd {
	switch k {
	case "down", "j":
		m.auditLog.MoveDown()
	case "up", "k":
		m.auditLog.MoveUp()
	case "v":
		return m.doVerifyAuditChain()
	}
	return nil
}

func (m *Model) snackbarMenuKey(k string) tea.Cmd {
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
	return nil
}

func (m *Model) snackbarManageKey(k string) tea.Cmd {
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
	return nil
}

func (m *Model) navKey(k string) tea.Cmd {
	switch k {
	case "l":
		return func() tea.Msg { return pages.NavigateMsg{Page: "login"} }
	case "b":
		bal := 0.0
		if m.userResp != nil {
			bal = m.userResp.Balance
		}
		m.snackBarMenu = pages.NewSnackBarMenuModel(bal)
		m.snackBarMenu.SetItems(nil)
		m.screen = scrSnackBarMenu
		return m.loadSnackBarMenu()
	case "m":
		pts := 0
		if m.userResp != nil {
			pts = m.userResp.PopcornPoints
		}
		m.merch = pages.NewMerchModel(pts)
		m.screen = scrMerch
		return m.loadMerch()
	case "i":
		m.inventory = pages.NewInventoryModel()
		m.screen = scrInventory
		return m.loadInventory()
	case "t":
		sub := "wood"
		bal := 0.0
		if m.userResp != nil {
			sub = m.userResp.Subscription
			bal = m.userResp.Balance
		}
		m.tierShop = pages.NewTierShopModel(sub, bal)
		m.screen = scrTierShop
	case "2":
		return m.doProfileTOTP()
	case "$":
		return m.doTopUp()
	}
	return nil
}

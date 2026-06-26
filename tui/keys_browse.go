package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) browseKey(k string) tea.Cmd {
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
		if !m.searching {
			m.searching = true
			m.searchBar.Focus()
		}
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
			return m.loadMovies(1, m.browse.Genre)
		}
	case "[":
		m.tabs.Prev()
		return m.applyMediaTypeFilter()
	case "]":
		m.tabs.Next()
		return m.applyMediaTypeFilter()
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
		if m.browse.Mode == pages.ModeAll && m.browse.HasNextPage() {
			m.browse.Selected = -1
			return m.loadMovies(m.browse.Page+1, m.browse.Genre)
		}
	case "b":
		if m.browse.Mode == pages.ModeAll && m.browse.HasPrevPage() {
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
	return nil
}

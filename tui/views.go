package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/tui/pages"
	"github.com/thelastvideostore/tui/styles"
)

func (m *Model) View() string {
	if !m.ready {
		return "loading..."
	}
	fh := lipgloss.Height(m.footerView())
	hh := lipgloss.Height(m.headerView())
	ch := m.h - hh - fh
	if ch < 5 {
		ch = 5
	}

	var body string
	switch m.screen {
	case scrSplash:
		return m.splash.View(m.w, m.h)
	case scrLogin:
		body = m.login.View(m.w, ch)
	case scrRegister:
		body = m.register.View(m.w, ch)
	case scrTOTP:
		body = m.totpView(m.w, ch)
	case scrBrowse:
		body = m.browse.View(m.w, ch)
		if m.browse.Mode == pages.ModeAll {
			tabsView := m.tabs.View(m.w)
			if m.genreTabs != nil && len(m.genreTabs.ActiveTab()) > 0 && m.tabs.ActiveTab() != "🍿 SnackBar" {
				tabsView = lipgloss.JoinVertical(lipgloss.Left, tabsView, m.genreTabs.View(m.w))
			}
			body = lipgloss.JoinVertical(lipgloss.Left, tabsView, "", body)
		}
		if m.searching {
			body = lipgloss.JoinVertical(lipgloss.Left, m.searchBar.View(), "", body)
		}
	case scrDetail:
		body = m.detail.View(m.w, ch)
	case scrRentals:
		body = m.rentals.View(m.w, ch)
	case scrProfile:
		body = m.profile.View(m.w, ch)
	case scrWishlist:
		body = m.wishlist.View(m.w, ch)
	case scrMerch:
		body = m.merch.View(m.w, ch)
	case scrInventory:
		body = m.inventory.View(m.w, ch)
	case scrTierShop:
		body = m.tierShop.View(m.w, ch)
	case scrMovieForm:
		body = m.movieForm.View(m.w, ch)
	case scrAccessDenied:
		body = m.accessDenied.View(m.w, ch)
	case scrAdminMovies:
		body = m.adminMovies.View(m.w, ch)
	case scrAdminUsers:
		body = m.adminUsers.View(m.w, ch)
	case scrAuditLog:
		body = m.auditLog.View(m.w, ch)
	case scrSnackBarMenu:
		body = m.snackBarMenu.View(m.w, ch)
	case scrSnackBarOrders:
		body = m.snackBarOrders.View(m.w, ch)
	case scrSnackBarManage:
		body = m.snackBarManage.View(m.w, ch)
	case scrGameDetail:
		body = m.gameDetail.View(m.w, ch)
	case scrGameSessions:
		body = m.gameSessions.View(m.w, ch)
	case scrMyPlaySessions:
		body = m.myPlaySessions.View(m.w, ch)
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(), body, m.footerView())
}

func (m *Model) headerView() string {
	un, tier, sub := "", "", "wood"
	pts, free, bal := 0, 0, 0.0
	loggedIn := m.userResp != nil
	if loggedIn {
		un = m.userResp.Username
		tier = m.userResp.TierName
		pts = m.userResp.PopcornPoints
		free = m.userResp.FreeRentals
		bal = m.userResp.Balance
		sub = m.userResp.Subscription
	}
	return m.header.View(m.w, loggedIn, un, tier, pts, free, bal, sub)
}

func (m *Model) footerView() string {
	if m.screen == scrSplash {
		return ""
	}
	helpView := m.help.View(m.currentScreenKeys())
	// Add a small visual gap above the help line so it doesn't sit
	// flush against the body content.
	return lipgloss.NewStyle().
		Padding(1, 0, 0, 0).
		Render(helpView)
}

func (m *Model) totpView(w, h int) string {
	masked := ""
	for i := 0; i < 6; i++ {
		if i < len(m.totpCode) {
			masked += styles.HighlightStyle.Render(string(m.totpCode[i]))
		} else {
			masked += styles.DimTextStyle.Render("_")
		}
		if i < 5 {
			masked += " "
		}
	}
	content := lipgloss.JoinVertical(lipgloss.Center,
		styles.HeadingStyle.Render("🔒 TWO-FACTOR AUTHENTICATION"),
		"",
		styles.TextStyle.Render("Enter the 6-digit code from your authenticator app:"),
		"",
		masked,
	)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

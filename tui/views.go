package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/ds/bitmask"
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
	var hints string
	switch m.screen {
	case scrSplash:
		hints = "[ENTER] start  [Ctrl+C] quit"
	case scrLogin:
		hints = "[TAB] switch  [ENTER] login  [Ctrl+R] sign up  [Ctrl+C] quit"
	case scrRegister:
		hints = "[TAB] switch  [ENTER] create account  [Ctrl+L] back to login"
	case scrTOTP:
		hints = "Enter 6-digit TOTP code  [ENTER] submit  [Ctrl+C] quit"
	case scrBrowse:
		if m.searching {
			hints = "[↑↓] results  [ENTER] open  [ESC] cancel search  [Ctrl+C] quit"
		} else {
			hints = "[↑↓] navigate  [ENTER] details  [[/]] tab  [,/.] genre  [N/B] pages  [S] staff picks  [L] last chance  [A] all  [R] rentals  [P] profile  [C] snack bar  [V] wishlist  [/] search  [F5] refresh  [Ctrl+C] quit"
		}
	case scrDetail:
		if m.detail != nil && !m.detail.Rented {
			hints = "[ENTER] rent  [↑↓] related  [W] waitlist  [F5] refresh  [Q] back  [Ctrl+C] quit"
		} else {
			hints = "[↑↓] related  [W] waitlist  [F5] refresh  [Q] back  [Ctrl+C] quit"
		}
	case scrRentals:
		hints = "[↑↓] select  [ENTER] return  [E] extend (30🍿)  [Q] back"
	case scrProfile:
		hints = "[L] logout  [T] tiers  [2] TOTP  [$] top-up  [B] snack bar  [M] rewards  [I] inventory  [Q] back"
	case scrSnackBarMenu:
		hints = "[↑↓] select  [ENTER] order  [O] orders"
		if m.userResp != nil && bitmask.CanSnackBarManage(m.userResp.Tier) {
			hints += "  [M] manage"
		}
		hints += "  [Q] back"
	case scrSnackBarOrders:
		hints = "[Q] back to snack bar"
	case scrSnackBarManage:
		hints = "[↑↓] select  [R] restock  [Q] back to snack bar"
	case scrGameDetail:
		hints = "[R] rent  [P] play  [1-5] duration  [E] end play  [↑↓] related  [Q] back"
	case scrGameSessions:
		hints = "[Q] back"
	case scrWishlist:
		hints = "[↑↓] select  [ENTER] info  [D] remove  [Q] back"
	case scrMerch:
		hints = "[↑↓] select  [ENTER] redeem  [Q] back"
	case scrInventory:
		hints = "[Q] back"
	case scrTierShop:
		hints = "[↑↓] select  [ENTER] purchase  [Q] back"
	case scrMovieForm:
		hints = "[TAB] next field  [ENTER] submit  [ESC] back"
	case scrAccessDenied:
		hints = "[Q] back"
	case scrAdminMovies:
		hints = fmt.Sprintf("[A] Add  [ENTER] Edit  [D] Delete  [S] Staff Pick  [N/B] Page %d/%d  [Q] Back", m.adminMovies.Page, m.adminMovies.TotalPages)
	case scrAdminUsers:
		hints = "[P] Promote  [D] Demote  [B] Toggle Ban  [T] Toggle TOTP  [Q] Back"
	case scrAuditLog:
		hints = "[↑↓] Navigate  [V] Verify Chain  [Q] Back"
	default:
		hints = "[Q] back  [Ctrl+C] quit"
	}
	return lipgloss.NewStyle().Background(styles.BgBlue).Foreground(styles.TextMedium).Width(m.w).Padding(0, 1).Render(hints)
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

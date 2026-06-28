package tui

import (
	"github.com/thelastvideostore/internal/ds/bitmask"
)

func (m *Model) currentScreenKeys() screenKeyMap {
	switch m.screen {
	case scrSplash:
		return helpWith(splashKeys{})
	case scrLogin:
		return helpWith(authKeys{})
	case scrRegister:
		return helpWith(authKeys{isRegister: true})
	case scrTOTP:
		return helpWith(totpKeys{})
	case scrBrowse:
		return helpWith(browseKeys{})
	case scrDetail:
		rented := m.detail != nil && m.detail.Rented
		return helpWith(movieDetailKeys{rented: rented})
	case scrGameDetail:
		return helpWith(gameDetailKeys{})
	case scrRentals:
		return helpWith(rentalsKeys{})
	case scrProfile:
		return helpWith(profileKeys{})
	case scrWishlist:
		return helpWith(wishlistKeys{})
	case scrMerch:
		return helpWith(merchKeys{})
	case scrInventory:
		return helpWith(inventoryKeys{})
	case scrTierShop:
		return helpWith(tierShopKeys{})
	case scrSnackBarMenu:
		canManage := m.userResp != nil && bitmask.CanSnackBarManage(m.userResp.Tier)
		return helpWith(snackBarMenuKeys{canManage: canManage})
	case scrSnackBarOrders:
		return helpWith(snackBarOrdersKeys{})
	case scrSnackBarManage:
		return helpWith(snackBarManageKeys{})
	case scrMyPlaySessions:
		return helpWith(myPlaySessionsKeys{})
	case scrAdminMovies:
		return helpWith(adminMoviesKeys{page: m.adminMovies.Page})
	case scrAdminUsers:
		return helpWith(adminUsersKeys{})
	case scrAuditLog:
		return helpWith(auditLogKeys{})
	}
	return helpWith(accessDeniedKeys{})
}

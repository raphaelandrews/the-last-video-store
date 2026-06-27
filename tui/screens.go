package tui

type screen int

const (
	scrSplash screen = iota
	scrLogin
	scrRegister
	scrTOTP
	scrBrowse
	scrDetail
	scrRentals
	scrProfile
	scrWishlist
	scrMerch
	scrInventory
	scrTierShop
	scrAdminMovies
	scrAdminUsers
	scrAuditLog
	scrMovieForm
	scrAccessDenied
	scrSnackBarMenu
	scrSnackBarOrders
	scrSnackBarManage
	scrGameDetail
	scrMyPlaySessions
)

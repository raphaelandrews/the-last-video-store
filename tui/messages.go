package tui

import (
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

type loadMoviesMsg struct {
	movies []models.MovieResponse
	total  int
	page   int
	reqID  int
}
type loadRentalsMsg struct {
	rentals []models.RentalResponse
}
type loadMyPlaySessionsMsg struct {
	sessions []models.GameSession
}
type loadProfileMsg struct {
	stats *pages.RentalStats
}
type loadWishlistMsg struct {
	items []pages.WishlistItem
}
type searchResultsMsg struct {
	results []models.MovieResponse
}
type wishlistResultMsg struct{}
type refreshMeMsg struct {
	user *models.UserResponse
}
type autoRefreshMsg struct{}
type refreshDetailMsg struct {
	movie *models.MovieResponse
}
type loadAdminMoviesMsg struct {
	movies []models.MovieResponse
	total  int
	page   int
}
type loadAdminUsersMsg struct {
	users []models.UserResponse
}
type loadAuditLogMsg struct {
	entries []map[string]interface{}
}
type loadMerchMsg struct {
	items []models.MerchItem
}
type loadInventoryMsg struct {
	items []pages.InventoryItem
}
type loadSnackBarMenuMsg struct {
	items []models.SnackBarItem
}
type loadSnackBarOrdersMsg struct {
	orders []models.SnackBarOrder
}
type loadSnackBarManageMsg struct {
	items []models.SnackBarItem
}
type gameRefreshMsg struct{}
type loadCatalogOptionsMsg struct {
	genres  []string
	formats []string
}

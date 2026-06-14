package models

type WishlistItem struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	MovieID string `json:"movie_id"`
	AddedAt int64  `json:"added_at"`
}

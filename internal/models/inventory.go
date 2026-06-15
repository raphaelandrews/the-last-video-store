package models

type InventoryItem struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	MerchID    string `json:"merch_id"`
	Name       string `json:"name"`
	RedeemedAt int64  `json:"redeemed_at"`
}

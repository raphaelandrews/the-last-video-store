package models

type SnackBarItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	Emoji       string  `json:"emoji"`
}

type SnackBarOrder struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	ItemID    string  `json:"item_id"`
	ItemName  string  `json:"item_name"`
	Emoji     string  `json:"emoji"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	OrderedAt int64   `json:"ordered_at"`
}

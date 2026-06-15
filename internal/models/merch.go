package models

type MerchItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PointsCost  int    `json:"points_cost"`
	Stock       int    `json:"stock"`
}

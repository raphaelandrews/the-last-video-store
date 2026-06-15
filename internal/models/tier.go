package models

type TierInfo struct {
	Name          string  `json:"name"`
	Label         string  `json:"label"`
	Price         float64 `json:"price"`
	FreeRentals   int     `json:"free_rentals"`
	MaxConcurrent int     `json:"max_concurrent"`
	NewReleasesOK bool    `json:"new_releases_ok"`
	NoLateFees    bool    `json:"no_late_fees"`
}

var Tiers = []TierInfo{
	{Name: "wood", Label: "Wood", Price: 0, FreeRentals: 0, MaxConcurrent: 2, NewReleasesOK: false, NoLateFees: false},
	{Name: "bronze", Label: "Bronze", Price: 9.99, FreeRentals: 1, MaxConcurrent: 3, NewReleasesOK: false, NoLateFees: false},
	{Name: "silver", Label: "Silver", Price: 19.99, FreeRentals: 3, MaxConcurrent: 5, NewReleasesOK: true, NoLateFees: false},
	{Name: "gold", Label: "Gold", Price: 29.99, FreeRentals: 5, MaxConcurrent: 10, NewReleasesOK: true, NoLateFees: true},
	{Name: "diamond", Label: "Diamond", Price: 49.99, FreeRentals: 999, MaxConcurrent: 999, NewReleasesOK: true, NoLateFees: true},
}

func TierByName(name string) *TierInfo {
	for i := range Tiers {
		if Tiers[i].Name == name {
			return &Tiers[i]
		}
	}
	return &Tiers[0]
}

func RentalCost(format string) float64 {
	switch format {
	case "VHS":
		return 2.99
	case "DVD":
		return 3.99
	case "Blu-ray":
		return 4.99
	default:
		return 3.99
	}
}

func MovieCost(movieRentalPrice float64, format string) float64 {
	if movieRentalPrice > 0 {
		return movieRentalPrice
	}
	return RentalCost(format)
}

package models

import (
	"github.com/thelastvideostore/internal/ds/bitmask"
)

type User struct {
	ID            string             `json:"id"`
	Username      string             `json:"username"`
	PasswordHash  string             `json:"password_hash"`
	Tier          bitmask.Permission `json:"tier"`
	MaxRentals    int                `json:"max_rentals"`
	RentalCount   int                `json:"rental_count"`
	Banned        bool               `json:"banned"`
	TOTPEnabled   bool               `json:"totp_enabled"`
	TOTPSecret    string             `json:"totp_secret"`
	PopcornPoints int                `json:"popcorn_points"`
	FreeRentals   int                `json:"free_rentals"`
	Balance       float64            `json:"balance"`
	Subscription  string             `json:"subscription"`
	CreatedAt     int64              `json:"created_at"`
	UpdatedAt     int64              `json:"updated_at"`
}

type UserResponse struct {
	ID            string             `json:"id"`
	Username      string             `json:"username"`
	Tier          bitmask.Permission `json:"tier"`
	TierName      string             `json:"tier_name"`
	TierNamePT    string             `json:"tier_name_pt"`
	MaxRentals    int                `json:"max_rentals"`
	RentalCount   int                `json:"rental_count"`
	Banned        bool               `json:"banned"`
	TOTPEnabled   bool               `json:"totp_enabled"`
	PopcornPoints int                `json:"popcorn_points"`
	FreeRentals   int                `json:"free_rentals"`
	Balance       float64            `json:"balance"`
	Subscription  string             `json:"subscription"`
	CreatedAt     int64              `json:"created_at"`
	UpdatedAt     int64              `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		ID:            u.ID,
		Username:      u.Username,
		Tier:          u.Tier,
		TierName:      bitmask.TierName(u.Tier),
		TierNamePT:    bitmask.TierNamesPT[u.Tier],
		MaxRentals:    u.MaxRentals,
		RentalCount:   u.RentalCount,
		Banned:        u.Banned,
		TOTPEnabled:   u.TOTPEnabled,
		PopcornPoints: u.PopcornPoints,
		FreeRentals:   u.FreeRentals,
		Balance:       u.Balance,
		Subscription:  u.Subscription,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
	if resp.TierName == "" {
		resp.TierName = bitmask.OwnerLabel()
	}
	if resp.TierNamePT == "" {
		resp.TierNamePT = bitmask.OwnerLabelPT()
	}
	return resp
}

func (u *User) CanRent() bool {
	return (bitmask.CanRent(u.Tier) || u.FreeRentals > 0) && !u.Banned
}

func (u *User) CanReserve() bool {
	return bitmask.CanReserve(u.Tier)
}

func (u *User) HasStaffAccess() bool {
	return bitmask.IsStaff(u.Tier)
}

func (u *User) CanManageUsers() bool {
	return bitmask.CanManageUsers(u.Tier)
}

func (u *User) CanAdmin() bool {
	return bitmask.CanAdmin(u.Tier)
}

func (u *User) IsOwner() bool {
	return bitmask.IsOwner(u.Tier)
}

func (u *User) AtRentalLimit() bool {
	return u.RentalCount >= u.MaxRentals
}

func (u *User) RemainingRentals() int {
	remaining := u.MaxRentals - u.RentalCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

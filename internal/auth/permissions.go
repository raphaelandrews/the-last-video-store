package auth

import (
	"github.com/thelastvideostore/internal/ds/bitmask"
)

func RequirePermission(userPerms bitmask.Permission, required bitmask.Permission) bool {
	return bitmask.Has(userPerms, required)
}

func TierName(perm bitmask.Permission) string {
	return bitmask.TierName(perm)
}

func MaxRentalsForTier(perm bitmask.Permission) int {
	return bitmask.MaxRentalsForTier(perm)
}

func CanAccessAdmin(perm bitmask.Permission) bool {
	return bitmask.CanAdmin(perm)
}

func IsStaff(perm bitmask.Permission) bool {
	return bitmask.IsStaff(perm)
}

func CanManageUsers(perm bitmask.Permission) bool {
	return bitmask.CanManageUsers(perm)
}

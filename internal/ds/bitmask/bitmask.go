package bitmask

type Permission uint16

const (
	PermBrowse         Permission = 0b0000000001
	PermRent           Permission = 0b0000000010
	PermReserve        Permission = 0b0000000100
	PermManageUsers    Permission = 0b0000001000
	PermStaff          Permission = 0b0000010000
	PermAdmin          Permission = 0b0000100000
	PermSnackBarAccess Permission = 0b0001000000
	PermSnackBarManage Permission = 0b0010000000
	PermGameAccess     Permission = 0b0100000000
	PermGameManage     Permission = 0b1000000000
)

const (
	TierBronze            Permission = PermBrowse | PermSnackBarAccess | PermGameAccess
	TierSilver            Permission = PermBrowse | PermRent | PermSnackBarAccess | PermGameAccess
	TierGold              Permission = PermBrowse | PermRent | PermReserve | PermSnackBarAccess | PermGameAccess
	TierEmployee          Permission = PermBrowse | PermRent | PermReserve | PermStaff | PermSnackBarAccess | PermGameAccess
	TierSupervisor        Permission = PermBrowse | PermRent | PermReserve | PermManageUsers | PermStaff | PermSnackBarAccess | PermGameAccess
	TierManager           Permission = PermBrowse | PermRent | PermReserve | PermManageUsers | PermStaff | PermAdmin | PermSnackBarAccess | PermSnackBarManage | PermGameAccess | PermGameManage
	TierOwner             Permission = TierManager
	TierSnackBarAttendant Permission = PermBrowse | PermStaff | PermSnackBarAccess
	TierSnackBarManager   Permission = PermBrowse | PermStaff | PermSnackBarAccess | PermSnackBarManage
	TierGameAttendant     Permission = PermBrowse | PermStaff | PermGameAccess
	TierGameManager       Permission = PermBrowse | PermStaff | PermGameAccess | PermGameManage
)

var TierPromotionOrder = []Permission{
	TierBronze, TierSilver, TierGold, TierEmployee, TierSupervisor, TierManager, TierOwner,
}

var TierLabels = map[Permission]string{
	TierBronze:            "Bronze",
	TierSilver:            "Silver",
	TierGold:              "Gold",
	TierEmployee:          "Employee",
	TierSupervisor:        "Supervisor",
	TierManager:           "Manager",
	TierSnackBarAttendant: "SnackBar Attendant",
	TierSnackBarManager:   "SnackBar Manager",
	TierGameAttendant:     "Game Attendant",
	TierGameManager:       "Game Manager",
}

var TierNamesPT = map[Permission]string{
	TierBronze:            "Cliente Bronze",
	TierSilver:            "Cliente Prata",
	TierGold:              "Cliente Ouro",
	TierEmployee:          "Atendente",
	TierSupervisor:        "Supervisor",
	TierManager:           "Gerente",
	TierSnackBarAttendant: "Atendente da Lanchonete",
	TierSnackBarManager:   "Gerente da Lanchonete",
	TierGameAttendant:     "Atendente de Games",
	TierGameManager:       "Gerente de Games",
}

func IsOwner(p Permission) bool {
	return p == TierOwner
}

func OwnerLabel() string {
	return "Owner"
}

func OwnerLabelPT() string {
	return "Dono"
}

func Has(p, flag Permission) bool {
	return p&flag != 0
}

func Set(p, flag Permission) Permission {
	return p | flag
}

func Clear(p, flag Permission) Permission {
	return p &^ flag
}

func Toggle(p, flag Permission) Permission {
	return p ^ flag
}

func CanRent(p Permission) bool {
	return Has(p, PermRent)
}

func CanReserve(p Permission) bool {
	return Has(p, PermReserve)
}

func IsStaff(p Permission) bool {
	return Has(p, PermStaff)
}

func CanManageUsers(p Permission) bool {
	return Has(p, PermManageUsers)
}

func CanAdmin(p Permission) bool {
	return Has(p, PermAdmin)
}

func CanSnackBarOrder(p Permission) bool {
	return Has(p, PermSnackBarAccess)
}

func CanSnackBarManage(p Permission) bool {
	return Has(p, PermSnackBarManage)
}

func CanGameAccess(p Permission) bool {
	return Has(p, PermGameAccess)
}

func CanGameManage(p Permission) bool {
	return Has(p, PermGameManage)
}

func TierName(p Permission) string {
	if p == TierOwner {
		return "Manager"
	}
	if name, ok := TierLabels[p]; ok {
		return name
	}

	best := uint(0)
	bestName := "Unknown"
	for tier := range TierLabels {
		count := popcount(uint16(p & tier))
		if count > best {
			best = count
			bestName = TierLabels[tier]
		}
	}
	return bestName
}

func MaxRentalsForTier(p Permission) int {
	switch {
	case Has(p, PermAdmin):
		return 10
	case Has(p, PermManageUsers):
		return 5
	case Has(p, PermStaff):
		return 5
	case Has(p, PermReserve):
		return 5
	case Has(p, PermRent):
		return 2
	default:
		return 1
	}
}

func popcount(x uint16) uint {
	var n uint
	for x != 0 {
		n++
		x &= x - 1
	}
	return n
}

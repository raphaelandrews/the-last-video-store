package bitmask

import (
	"testing"
)

func TestHas(t *testing.T) {
	tests := []struct {
		name    string
		base    Permission
		flag    Permission
		wantHas bool
	}{
		{"bronze_has_browse", TierBronze, PermBrowse, true},
		{"bronze_missing_rent", TierBronze, PermRent, false},
		{"bronze_missing_reserve", TierBronze, PermReserve, false},
		{"bronze_missing_staff", TierBronze, PermStaff, false},
		{"bronze_missing_admin", TierBronze, PermAdmin, false},
		{"silver_has_rent", TierSilver, PermRent, true},
		{"silver_missing_reserve", TierSilver, PermReserve, false},
		{"gold_has_reserve", TierGold, PermReserve, true},
		{"gold_missing_staff", TierGold, PermStaff, false},
		{"gold_not_employee", TierGold, PermStaff, false},
		{"employee_has_staff", TierEmployee, PermStaff, true},
		{"employee_not_gold_same", TierEmployee, PermStaff, true},
		{"supervisor_has_users", TierSupervisor, PermManageUsers, true},
		{"supervisor_missing_admin", TierSupervisor, PermAdmin, false},
		{"manager_has_all", TierManager, PermAdmin, true},
		{"owner_has_all", TierOwner, PermAdmin, true},
		{"owner_has_browse", TierOwner, PermBrowse, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Has(tt.base, tt.flag); got != tt.wantHas {
				t.Errorf("Has(%06b, %06b) = %v, want %v", tt.base, tt.flag, got, tt.wantHas)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name    string
		base    Permission
		flag    Permission
		wantSet Permission
	}{
		{"bronze_gets_rent", TierBronze, PermRent, TierSilver},
		{"silver_gets_reserve", TierSilver, PermReserve, TierGold},
		{"set_already_present", TierSilver, PermBrowse, TierSilver},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Set(tt.base, tt.flag); got != tt.wantSet {
				t.Errorf("Set(%06b, %06b) = %06b, want %06b", tt.base, tt.flag, got, tt.wantSet)
			}
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name      string
		base      Permission
		flag      Permission
		wantClear Permission
	}{
		{"silver_loses_rent", TierSilver, PermRent, TierBronze},
		{"gold_loses_reserve", TierGold, PermReserve, TierSilver},
		{"clear_missing_bit", TierBronze, PermStaff, TierBronze},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clear(tt.base, tt.flag); got != tt.wantClear {
				t.Errorf("Clear(%06b, %06b) = %06b, want %06b", tt.base, tt.flag, got, tt.wantClear)
			}
		})
	}
}

func TestToggle(t *testing.T) {
	tests := []struct {
		name       string
		base       Permission
		flag       Permission
		wantToggle Permission
	}{
		{"toggle_on", TierBronze, PermRent, TierSilver},
		{"toggle_off", TierSilver, PermRent, TierBronze},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Toggle(tt.base, tt.flag); got != tt.wantToggle {
				t.Errorf("Toggle(%06b, %06b) = %06b, want %06b", tt.base, tt.flag, got, tt.wantToggle)
			}
		})
	}
}

func TestCanRent(t *testing.T) {
	if CanRent(TierBronze) {
		t.Error("Bronze should not be able to rent with PermBrowse only")
	}
	if !CanRent(TierSilver) {
		t.Error("Silver should be able to rent")
	}
	if !CanRent(TierGold) {
		t.Error("Gold should be able to rent")
	}
}

func TestCanReserve(t *testing.T) {
	if CanReserve(TierSilver) {
		t.Error("Silver should not be able to reserve new releases")
	}
	if !CanReserve(TierGold) {
		t.Error("Gold should be able to reserve")
	}
}

func TestIsStaff(t *testing.T) {
	if IsStaff(TierGold) {
		t.Error("Gold should NOT be staff (no PermStaff bit)")
	}
	if !IsStaff(TierEmployee) {
		t.Error("Employee should be staff")
	}
	if !IsStaff(TierSupervisor) {
		t.Error("Supervisor should be staff")
	}
	if !IsStaff(TierManager) {
		t.Error("Manager should be staff")
	}
}

func TestGoldNotEmployee(t *testing.T) {
	if TierGold == TierEmployee {
		t.Error("Gold and Employee must have different bitmasks")
	}
	if Has(TierGold, PermStaff) {
		t.Error("Gold must not have PermStaff bit")
	}
	if !Has(TierEmployee, PermStaff) {
		t.Error("Employee must have PermStaff bit")
	}
}

func TestTierName(t *testing.T) {
	tests := []struct {
		perm Permission
		want string
	}{
		{TierBronze, "Bronze"},
		{TierSilver, "Silver"},
		{TierGold, "Gold"},
		{TierEmployee, "Employee"},
		{TierSupervisor, "Supervisor"},
		{TierManager, "Manager"},
		{TierOwner, "Manager"},
		{TierBronze | PermRent, "Silver"},
		{0, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := TierName(tt.perm); got != tt.want {
				t.Errorf("TierName(%06b) = %q, want %q", tt.perm, got, tt.want)
			}
		})
	}
}

func TestMaxRentalsForTier(t *testing.T) {
	tests := []struct {
		perm Permission
		want int
	}{
		{TierBronze, 1},
		{TierSilver, 2},
		{TierGold, 5},
		{TierEmployee, 5},
		{TierSupervisor, 5},
		{TierManager, 10},
		{TierOwner, 10},
	}

	for _, tt := range tests {
		t.Run(TierName(tt.perm), func(t *testing.T) {
			if got := MaxRentalsForTier(tt.perm); got != tt.want {
				t.Errorf("MaxRentalsForTier(%s) = %d, want %d", TierName(tt.perm), got, tt.want)
			}
		})
	}
}

func TestIsOwner(t *testing.T) {
	if IsOwner(TierBronze) {
		t.Error("Bronze is not owner")
	}
	if !IsOwner(TierManager) {
		t.Error("TierManager should match IsOwner (same bitmask)")
	}
	if !IsOwner(TierOwner) {
		t.Error("TierOwner should be owner")
	}
}

func TestManagerOwnerSameBitmask(t *testing.T) {
	if TierManager != TierOwner {
		t.Error("TierManager and TierOwner must share the same bitmask")
	}
	if TierName(TierManager) != "Manager" {
		t.Error("TierName for full mask returns Manager by default")
	}
}

func BenchmarkHas(b *testing.B) {
	for b.Loop() {
		Has(TierManager, PermAdmin)
	}
}

func BenchmarkSet(b *testing.B) {
	for b.Loop() {
		Set(TierBronze, PermRent)
	}
}

func BenchmarkClear(b *testing.B) {
	for b.Loop() {
		Clear(TierSilver, PermRent)
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

func seedUsers(s *store.Store) {
	entries := []struct {
		name, pass, sub string
		tier            bitmask.Permission
		banned          bool
		balance         float64
	}{
		{"bronze", "123", "bronze", bitmask.TierBronze, false, 50},
		{"silver", "123", "silver", bitmask.TierSilver, false, 50},
		{"gold", "123", "gold", bitmask.TierGold, false, 50},
		{"employee", "123", "gold", bitmask.TierEmployee, false, 50},
		{"supervisor", "123", "gold", bitmask.TierSupervisor, false, 50},
		{"manager", "123", "diamond", bitmask.TierManager, false, 100},
		{"owner", "123", "diamond", bitmask.TierOwner, false, 100},
		{"banned", "123", "wood", bitmask.TierBronze, true, 5},
		{"bar_attendant", "123", "wood", bitmask.TierSnackBarAttendant, false, 30},
		{"bar_manager", "123", "wood", bitmask.TierSnackBarManager, false, 50},
		{"game_attendant", "123", "wood", bitmask.TierGameAttendant, false, 30},
		{"game_manager", "123", "wood", bitmask.TierGameManager, false, 50},
	}

	for _, e := range entries {
		hash, err := auth.HashPassword(e.pass)
		if err != nil {
			panic(fmt.Errorf("seed user: hash: %w", err))
		}
		now := time.Now().Unix()
		tier := models.TierByName(e.sub)
		user := &models.User{
			ID:            fmt.Sprintf("seed-%s", e.name),
			Username:      e.name,
			PasswordHash:  hash,
			Tier:          e.tier,
			MaxRentals:    tier.MaxConcurrent,
			Banned:        e.banned,
			PopcornPoints: 250,
			FreeRentals:   tier.FreeRentals,
			Balance:       e.balance,
			Subscription:  e.sub,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.CreateUser(user); err != nil {
			panic(fmt.Errorf("seed user: %s: %w", e.name, err))
		}
	}
}

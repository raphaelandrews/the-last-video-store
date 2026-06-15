package store

import (
	"encoding/json"
	"fmt"

	"github.com/thelastvideostore/internal/models"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) PurchaseTier(userID string, tier *models.TierInfo) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		data := ub.Get([]byte(userID))
		if data == nil {
			return fmt.Errorf("user not found")
		}
		var user models.User
		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}
		if tier.Price > 0 && user.Balance < tier.Price {
			return fmt.Errorf("insufficient balance: need $%.2f, have $%.2f", tier.Price, user.Balance)
		}
		if tier.Price > 0 {
			user.Balance -= tier.Price
		}
		user.Subscription = tier.Name
		user.MaxRentals = tier.MaxConcurrent
		user.FreeRentals = tier.FreeRentals

		updated, _ := json.Marshal(&user)
		return ub.Put([]byte(userID), updated)
	})
}

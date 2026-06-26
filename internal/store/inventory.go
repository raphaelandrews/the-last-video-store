package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) AddToInventory(userID string, item *models.InventoryItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketInventory)
		data, _ := json.Marshal(item)
		return b.Put([]byte(userID+":"+item.ID), data)
	})
}

func (s *Store) ListInventory(userID string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	prefix := []byte(userID + ":")
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketInventory)
		c := b.Cursor()
		for k, v := c.Seek(prefix); k != nil && hasBytePrefix(k, prefix); k, v = c.Next() {
			var item models.InventoryItem
			if err := json.Unmarshal(v, &item); err != nil {
				continue
			}
			items = append(items, item)
		}
		return nil
	})
	return items, err
}

func (s *Store) RedeemMerchItem(itemID, userID string) (*models.MerchItem, error) {
	var redeemed *models.MerchItem
	err := s.db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket(bucketMerch)
		data := mb.Get([]byte(itemID))
		if data == nil {
			return fmt.Errorf("merch item not found")
		}
		var m models.MerchItem
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		if m.Stock <= 0 {
			return fmt.Errorf("out of stock")
		}

		ub := tx.Bucket(bucketUsers)
		userData := ub.Get([]byte(userID))
		if userData == nil {
			return fmt.Errorf("user not found")
		}
		var user models.User
		if err := json.Unmarshal(userData, &user); err != nil {
			return err
		}
		if user.PopcornPoints < m.PointsCost {
			return fmt.Errorf("insufficient popcorn points: need %d, have %d", m.PointsCost, user.PopcornPoints)
		}

		user.PopcornPoints -= m.PointsCost

		switch itemID {
		case "merch-free-rental":
			user.FreeRentals++
		case "merch-screening":
			user.FreeRentals += 5
		case "merch-membership-upgrade":
			if err := upgradeTier(&user); err != nil {
				return err
			}
		default:
			ib := tx.Bucket(bucketInventory)
			inv := models.InventoryItem{
				ID:         fmt.Sprintf("inv-%d", time.Now().UnixNano()),
				UserID:     userID,
				MerchID:    m.ID,
				Name:       m.Name,
				RedeemedAt: time.Now().Unix(),
			}
			invData, _ := json.Marshal(&inv)
			ib.Put([]byte(userID+":"+inv.ID), invData)
		}

		m.Stock--

		updatedUser, _ := json.Marshal(&user)
		ub.Put([]byte(userID), updatedUser)

		updatedItem, _ := json.Marshal(&m)
		mb.Put([]byte(itemID), updatedItem)

		redeemed = &m
		return nil
	})
	return redeemed, err
}

func upgradeTier(u *models.User) error {
	tierNames := []string{"Couch Potato", "Matinee Fan", "Gold Member"}
	currentName := bitmask.TierName(u.Tier)
	maxTier := false

	for _, t := range tierNames {
		if currentName == t {
			next := u.Tier << 1
			if next > 0 {
				u.Tier = next
				u.MaxRentals = bitmask.MaxRentalsForTier(u.Tier)
			}
			maxTier = true
			break
		}
	}
	if !maxTier {
		return fmt.Errorf("tier %s cannot be upgraded further", bitmask.TierName(u.Tier))
	}
	return nil
}

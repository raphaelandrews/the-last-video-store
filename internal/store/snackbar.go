package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/models"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) CreateSnackBarItem(item *models.SnackBarItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return b.Put([]byte(item.ID), data)
	})
}

func (s *Store) ListSnackBarItems() ([]models.SnackBarItem, error) {
	var items []models.SnackBarItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		return b.ForEach(func(k, v []byte) error {
			var item models.SnackBarItem
			if err := json.Unmarshal(v, &item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})
	return items, err
}

func (s *Store) GetSnackBarItem(id string) (*models.SnackBarItem, error) {
	var item models.SnackBarItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("snackbar item not found: %s", id)
		}
		return json.Unmarshal(data, &item)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) UpdateSnackBarItem(item *models.SnackBarItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return b.Put([]byte(item.ID), data)
	})
}

func (s *Store) DeleteSnackBarItem(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		return b.Delete([]byte(id))
	})
}

func (s *Store) RestockSnackBarItem(id string, amount int) (*models.SnackBarItem, error) {
	var item models.SnackBarItem
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnackBarItems)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("snackbar item not found: %s", id)
		}
		if err := json.Unmarshal(data, &item); err != nil {
			return err
		}
		item.Stock += amount
		updated, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return b.Put([]byte(id), updated)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) PlaceSnackBarOrder(userID, itemID string, quantity int) (*models.SnackBarOrder, error) {
	var order models.SnackBarOrder
	err := s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		ib := tx.Bucket(bucketSnackBarItems)
		ob := tx.Bucket(bucketSnackBarOrders)
		ou := tx.Bucket(bucketSnackBarOrdersByUser)

		itemData := ib.Get([]byte(itemID))
		if itemData == nil {
			return fmt.Errorf("snackbar item not found: %s", itemID)
		}
		var item models.SnackBarItem
		if err := json.Unmarshal(itemData, &item); err != nil {
			return err
		}
		if item.Stock < quantity {
			return fmt.Errorf("insufficient stock: %s has %d available", item.Name, item.Stock)
		}

		userData := ub.Get([]byte(userID))
		if userData == nil {
			return fmt.Errorf("user not found: %s", userID)
		}
		var user models.User
		if err := json.Unmarshal(userData, &user); err != nil {
			return err
		}

		total := item.Price * float64(quantity)
		if user.Balance < total {
			return fmt.Errorf("insufficient balance: need $%.2f, have $%.2f", total, user.Balance)
		}

		user.Balance -= total
		updatedUser, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := ub.Put([]byte(userID), updatedUser); err != nil {
			return err
		}

		item.Stock -= quantity
		updatedItem, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if err := ib.Put([]byte(itemID), updatedItem); err != nil {
			return err
		}

		order = models.SnackBarOrder{
			ID:        uuid.NewString(),
			UserID:    userID,
			ItemID:    itemID,
			ItemName:  item.Name,
			Emoji:     item.Emoji,
			Quantity:  quantity,
			UnitPrice: item.Price,
			Total:     total,
			Status:    "completed",
			OrderedAt: time.Now().Unix(),
		}
		orderData, err := json.Marshal(order)
		if err != nil {
			return err
		}
		if err := ob.Put([]byte(order.ID), orderData); err != nil {
			return err
		}

		userOrderKey := fmt.Sprintf("%s:%s:%d", userID, order.ID, order.OrderedAt)
		return ou.Put([]byte(userOrderKey), nil)
	})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Store) ListSnackBarOrders(userID string) ([]models.SnackBarOrder, error) {
	var orders []models.SnackBarOrder
	err := s.db.View(func(tx *bolt.Tx) error {
		ob := tx.Bucket(bucketSnackBarOrders)
		ou := tx.Bucket(bucketSnackBarOrdersByUser)

		prefix := []byte(userID + ":")
		c := ou.Cursor()
		for k, _ := c.Seek(prefix); k != nil && hasBytePrefix(k, prefix); k, _ = c.Next() {
			parts := splitCompositeKey(string(k))
			if len(parts) >= 2 {
				orderData := ob.Get([]byte(parts[1]))
				if orderData != nil {
					var order models.SnackBarOrder
					if err := json.Unmarshal(orderData, &order); err != nil {
						return err
					}
					orders = append(orders, order)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.SnackBarOrder{}
	}
	return orders, nil
}

func splitCompositeKey(key string) []string {
	var parts []string
	start := 0
	count := 0
	for i := 0; i < len(key) && count < 3; i++ {
		if key[i] == ':' {
			parts = append(parts, key[start:i])
			start = i + 1
			count++
		}
	}
	if start < len(key) {
		parts = append(parts, key[start:])
	}
	return parts
}

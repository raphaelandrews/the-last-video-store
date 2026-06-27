package store

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/thelastvideostore/internal/models"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) CreateRental(rental *models.Rental) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		rb := tx.Bucket(bucketRentals)
		rub := tx.Bucket(bucketRentalsByUser)

		data, err := encode(rental)
		if err != nil {
			return err
		}

		if err := rb.Put([]byte(rental.ID), data); err != nil {
			return err
		}

		userKey := fmt.Sprintf("%s:%s:%d", rental.UserID, rental.ID, rental.RentedAt)
		if err := rub.Put([]byte(userKey), nil); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) GetRentalByID(id string) (*models.Rental, error) {
	var rental models.Rental
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketRentals).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("rental not found: %s", id)
		}
		return json.Unmarshal(data, &rental)
	})
	if err != nil {
		return nil, err
	}
	return &rental, nil
}

func (s *Store) UpdateRental(rental *models.Rental) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := encode(rental)
		if err != nil {
			return err
		}
		if err := tx.Bucket(bucketRentals).Put([]byte(rental.ID), data); err != nil {
			return err
		}
		rb := tx.Bucket(bucketRentalsByUser)
		if rb != nil {
			keysToDelete := [][]byte{}
			c := rb.Cursor()
			prefix := rental.UserID + ":"
			for k, _ := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, _ = c.Next() {
				parts := splitKey(string(k), ':')
				if len(parts) >= 2 && parts[1] == rental.ID {
					keysToDelete = append(keysToDelete, k)
				}
			}
			for _, k := range keysToDelete {
				rb.Delete(k)
			}
			newKey := []byte(rental.UserID + ":" + rental.ID + ":" + strconv.FormatInt(rental.RentedAt, 10))
			rb.Put(newKey, nil)
		}
		return nil
	})
}

func (s *Store) GetActiveRentalsByUser(userID string) ([]*models.Rental, error) {
	var rentals []*models.Rental
	now := time.Now().Unix()

	err := s.db.View(func(tx *bolt.Tx) error {
		rb := tx.Bucket(bucketRentals)
		rub := tx.Bucket(bucketRentalsByUser)

		c := rub.Cursor()
		prefix := userID + ":"
		for k, _ := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, _ = c.Next() {
			parts := splitKey(string(k), ':')
			if len(parts) < 2 {
				continue
			}
			rentalID := parts[1]
			data := rb.Get([]byte(rentalID))
			if data == nil {
				continue
			}
			var rental models.Rental
			if json.Unmarshal(data, &rental) != nil {
				continue
			}
			if rental.Status != models.RentalReturned {
				if rental.DueDate < now && rental.Status == models.RentalActive {
					rental.Status = models.RentalOverdue
				}
				rentals = append(rentals, &rental)
			}
		}
		return nil
	})

	return rentals, err
}

func (s *Store) GetRentalHistoryByUser(userID string) ([]*models.Rental, error) {
	var rentals []*models.Rental
	now := time.Now().Unix()

	err := s.db.View(func(tx *bolt.Tx) error {
		rb := tx.Bucket(bucketRentals)
		rub := tx.Bucket(bucketRentalsByUser)

		c := rub.Cursor()
		prefix := userID + ":"
		for k, _ := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, _ = c.Next() {
			parts := splitKey(string(k), ':')
			if len(parts) < 2 {
				continue
			}
			rentalID := parts[1]
			data := rb.Get([]byte(rentalID))
			if data == nil {
				continue
			}
			var rental models.Rental
			if json.Unmarshal(data, &rental) != nil {
				continue
			}
			if rental.DueDate < now && rental.Status == models.RentalActive {
				rental.Status = models.RentalOverdue
			}
			rentals = append(rentals, &rental)
		}
		return nil
	})

	return rentals, err
}

func (s *Store) GetOverdueRentals() ([]*models.Rental, error) {
	var overdue []*models.Rental
	now := time.Now().Unix()

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRentals)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rental models.Rental
			if json.Unmarshal(v, &rental) != nil {
				continue
			}
			if rental.Status != models.RentalReturned && rental.DueDate < now {
				rental.Status = models.RentalOverdue
				overdue = append(overdue, &rental)
			}
		}
		return nil
	})

	return overdue, err
}

func (s *Store) CountActiveRentalsByUser(userID string) (int, error) {
	count := 0
	err := s.db.View(func(tx *bolt.Tx) error {
		rb := tx.Bucket(bucketRentals)
		rub := tx.Bucket(bucketRentalsByUser)

		c := rub.Cursor()
		prefix := userID + ":"
		for k, _ := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, _ = c.Next() {
			parts := splitKey(string(k), ':')
			if len(parts) < 2 {
				continue
			}
			rentalID := parts[1]
			data := rb.Get([]byte(rentalID))
			if data == nil {
				continue
			}
			var rental models.Rental
			if json.Unmarshal(data, &rental) != nil {
				continue
			}
			if rental.Status != models.RentalReturned {
				count++
			}
		}
		return nil
	})
	return count, err
}

func (s *Store) ExtendRental(rentalID, userID string, extensionMinutes int64, cost int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		rb := tx.Bucket(bucketRentals)
		data := rb.Get([]byte(rentalID))
		if data == nil {
			return fmt.Errorf("rental not found")
		}
		var rental models.Rental
		if err := json.Unmarshal(data, &rental); err != nil {
			return err
		}
		if rental.UserID != userID {
			return fmt.Errorf("not your rental")
		}
		if rental.Status == models.RentalReturned {
			return fmt.Errorf("already returned")
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
		if user.PopcornPoints < cost {
			return fmt.Errorf("need %d popcorn points, have %d", cost, user.PopcornPoints)
		}

		user.PopcornPoints -= cost
		rental.DueDate += extensionMinutes * 60
		if rental.Status == models.RentalOverdue {
			rental.Status = models.RentalActive
		}

		updatedUser, _ := json.Marshal(&user)
		ub.Put([]byte(userID), updatedUser)

		updatedRental, _ := json.Marshal(&rental)
		rb.Put([]byte(rentalID), updatedRental)

		return nil
	})
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func splitKey(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

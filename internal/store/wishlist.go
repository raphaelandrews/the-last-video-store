package store

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/models"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) AddToWishlist(userID, movieID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketWishlists)
		item := &models.WishlistItem{
			ID:      uuid.NewString(),
			UserID:  userID,
			MovieID: movieID,
			AddedAt: time.Now().Unix(),
		}

		key := []byte(fmt.Sprintf("%s:%s", userID, item.ID))
		data, err := encode(item)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}

func (s *Store) RemoveFromWishlist(userID, movieID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketWishlists)
		prefix := userID + ":"
		c := b.Cursor()
		for k, v := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, v = c.Next() {
			var item models.WishlistItem
			if decode(v, &item) != nil {
				continue
			}
			if item.MovieID == movieID {
				return b.Delete(k)
			}
		}
		return nil
	})
}

func (s *Store) GetWishlist(userID string) ([]*models.WishlistItem, error) {
	var items []*models.WishlistItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketWishlists)
		prefix := userID + ":"
		c := b.Cursor()
		for k, v := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, v = c.Next() {
			var item models.WishlistItem
			if decode(v, &item) != nil {
				continue
			}
			items = append(items, &item)
		}
		return nil
	})
	return items, err
}

func (s *Store) IsInWishlist(userID, movieID string) (bool, error) {
	found := false
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketWishlists)
		prefix := userID + ":"
		c := b.Cursor()
		for k, v := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, v = c.Next() {
			var item models.WishlistItem
			if decode(v, &item) != nil {
				continue
			}
			if item.MovieID == movieID {
				found = true
				return nil
			}
		}
		return nil
	})
	return found, err
}

func (s *Store) GetWishlistSize(userID string) (int, error) {
	count := 0
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketWishlists)
		prefix := userID + ":"
		c := b.Cursor()
		for k, _ := c.Seek([]byte(prefix)); k != nil && hasPrefix(string(k), prefix); k, _ = c.Next() {
			count++
		}
		return nil
	})
	return count, err
}

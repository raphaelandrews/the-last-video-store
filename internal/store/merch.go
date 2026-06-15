package store

import (
	"encoding/json"
	"fmt"

	"github.com/thelastvideostore/internal/models"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) CreateMerchItem(m *models.MerchItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMerch)
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return b.Put([]byte(m.ID), data)
	})
}

func (s *Store) ListMerchItems() ([]models.MerchItem, error) {
	var items []models.MerchItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMerch)
		return b.ForEach(func(k, v []byte) error {
			var item models.MerchItem
			if err := json.Unmarshal(v, &item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})
	return items, err
}

func (s *Store) GetMerchItem(id string) (*models.MerchItem, error) {
	var item *models.MerchItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMerch)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("merch item not found: %s", id)
		}
		var m models.MerchItem
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		item = &m
		return nil
	})
	return item, err
}

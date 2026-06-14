package store

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) AddStaffPick(movieID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketStaffPicks)
		now := []byte(time.Now().Format(time.RFC3339))
		return b.Put([]byte(movieID), now)
	})
}

func (s *Store) RemoveStaffPick(movieID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketStaffPicks).Delete([]byte(movieID))
	})
}

func (s *Store) GetStaffPicks() ([]string, error) {
	var ids []string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketStaffPicks)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ids = append(ids, string(k))
		}
		return nil
	})
	return ids, err
}

func (s *Store) IsStaffPick(movieID string) bool {
	var found bool
	s.db.View(func(tx *bolt.Tx) error {
		found = tx.Bucket(bucketStaffPicks).Get([]byte(movieID)) != nil
		return nil
	})
	return found
}

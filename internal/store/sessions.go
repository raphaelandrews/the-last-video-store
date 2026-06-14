package store

import (
	"encoding/binary"
	"time"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveRefreshToken(userID, tokenID string, expiresAt int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)
		key := []byte("refresh:" + userID + ":" + tokenID)
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, uint64(expiresAt))
		return b.Put(key, val)
	})
}

func (s *Store) ValidateRefreshToken(userID, tokenID string) (bool, error) {
	var valid bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)
		key := []byte("refresh:" + userID + ":" + tokenID)
		data := b.Get(key)
		if data == nil {
			return nil
		}
		expiresAt := int64(binary.LittleEndian.Uint64(data))
		now := time.Now().Unix()
		valid = now < expiresAt
		return nil
	})
	return valid, err
}

func (s *Store) InvalidateRefreshToken(userID, tokenID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)
		key := []byte("refresh:" + userID + ":" + tokenID)
		return b.Delete(key)
	})
}

func (s *Store) InvalidateAllUserSessions(userID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)
		prefix := []byte("refresh:" + userID + ":")
		c := b.Cursor()
		for k, _ := c.Seek(prefix); k != nil && hasPrefix(string(k), string(prefix)); k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) IncrementFailedAttempts(username string) (int, error) {
	var count int
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLockouts)
		key := []byte("attempts:" + username)
		data := b.Get(key)
		if data != nil {
			count = int(binary.LittleEndian.Uint32(data))
		}
		count++
		val := make([]byte, 4)
		binary.LittleEndian.PutUint32(val, uint32(count))
		return b.Put(key, val)
	})
	return count, err
}

func (s *Store) ResetFailedAttempts(username string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketLockouts).Delete([]byte("attempts:" + username))
	})
}

func (s *Store) LockUserUntil(username string, until int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLockouts)
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, uint64(until))
		return b.Put([]byte("lock:"+username), val)
	})
}

func (s *Store) IsUserLocked(username string) (bool, error) {
	var locked bool
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketLockouts).Get([]byte("lock:" + username))
		if data == nil {
			return nil
		}
		until := int64(binary.LittleEndian.Uint64(data))
		locked = time.Now().Unix() < until
		return nil
	})
	return locked, err
}

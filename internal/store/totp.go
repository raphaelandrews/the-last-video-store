package store

import (
	"encoding/binary"
	"time"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveTOTPSecret(userID string, encryptedSecret []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketTOTPSecrets).Put([]byte(userID), encryptedSecret)
	})
}

func (s *Store) GetTOTPSecret(userID string) ([]byte, error) {
	var secret []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketTOTPSecrets).Get([]byte(userID))
		if data != nil {
			secret = make([]byte, len(data))
			copy(secret, data)
		}
		return nil
	})
	return secret, err
}

func (s *Store) DeleteTOTPSecret(userID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketTOTPSecrets).Delete([]byte(userID))
	})
}

func (s *Store) IncrementTOTPFailures(userID string) (int, error) {
	var count int
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLockouts)
		key := []byte("totp_attempts:" + userID)
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

func (s *Store) ResetTOTPFailures(userID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketLockouts).Delete([]byte("totp_attempts:" + userID))
	})
}

func (s *Store) LockTOTPUserUntil(userID string, until int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketLockouts)
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, uint64(until))
		return b.Put([]byte("totp_lock:"+userID), val)
	})
}

func (s *Store) IsTOTPLocked(userID string) (bool, error) {
	var locked bool
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketLockouts).Get([]byte("totp_lock:" + userID))
		if data == nil {
			return nil
		}
		until := int64(binary.LittleEndian.Uint64(data))
		locked = time.Now().Unix() < until
		return nil
	})
	return locked, err
}

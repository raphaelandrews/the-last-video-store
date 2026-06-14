package store

import (
	"encoding/json"
	"fmt"

	"github.com/thelastvideostore/internal/models"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) CreateUser(user *models.User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		unb := tx.Bucket(bucketUsersByUsername)

		data, err := encode(user)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		if err := ub.Put([]byte(user.ID), data); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		if err := unb.Put([]byte(user.Username), []byte(user.ID)); err != nil {
			return fmt.Errorf("create user index: %w", err)
		}

		return nil
	})
}

func (s *Store) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketUsers).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("user not found: %s", id)
		}
		return json.Unmarshal(data, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.View(func(tx *bolt.Tx) error {
		unb := tx.Bucket(bucketUsersByUsername)
		userID := unb.Get([]byte(username))
		if userID == nil {
			return fmt.Errorf("user not found: %s", username)
		}
		data := tx.Bucket(bucketUsers).Get(userID)
		if data == nil {
			return fmt.Errorf("user data missing: %s", username)
		}
		return json.Unmarshal(data, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) UpdateUser(user *models.User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		unb := tx.Bucket(bucketUsersByUsername)

		data, err := encode(user)
		if err != nil {
			return err
		}

		if err := ub.Put([]byte(user.ID), data); err != nil {
			return err
		}

		if err := unb.Put([]byte(user.Username), []byte(user.ID)); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) DeleteUser(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		unb := tx.Bucket(bucketUsersByUsername)

		data := ub.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("user not found: %s", id)
		}

		var user models.User
		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}

		if err := unb.Delete([]byte(user.Username)); err != nil {
			return err
		}

		return ub.Delete([]byte(id))
	})
}

func (s *Store) ListUsers() ([]*models.User, error) {
	var users []*models.User
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user models.User
			if err := json.Unmarshal(v, &user); err != nil {
				continue
			}
			users = append(users, &user)
		}
		return nil
	})
	return users, err
}

func (s *Store) UserExists(username string) bool {
	var exists bool
	s.db.View(func(tx *bolt.Tx) error {
		exists = tx.Bucket(bucketUsersByUsername).Get([]byte(username)) != nil
		return nil
	})
	return exists
}

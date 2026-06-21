package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/models"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) StartGameSession(userID, gameID, gameTitle string, hourlyRate float64, durationMinutes int) (*models.GameSession, error) {
	var session models.GameSession
	err := s.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket(bucketUsers)
		gsb := tx.Bucket(bucketGameSessions)
		gsu := tx.Bucket(bucketGameSessionsByUser)

		userData := ub.Get([]byte(userID))
		if userData == nil {
			return fmt.Errorf("user not found: %s", userID)
		}
		var user models.User
		if err := json.Unmarshal(userData, &user); err != nil {
			return err
		}
		totalCost := hourlyRate * float64(durationMinutes) / 60.0
		if user.Balance < totalCost {
			return fmt.Errorf("insufficient balance: need $%.2f, have $%.2f", totalCost, user.Balance)
		}
		user.Balance -= totalCost
		updatedUser, _ := json.Marshal(user)
		if err := ub.Put([]byte(userID), updatedUser); err != nil {
			return err
		}

		now := time.Now().Unix()
		session = models.GameSession{
			ID:              uuid.NewString(),
			UserID:          userID,
			GameID:          gameID,
			GameTitle:       gameTitle,
			StartedAt:       now,
			ExpiresAt:       now + int64(durationMinutes*60),
			DurationMinutes: durationMinutes,
			Cost:            totalCost,
			Status:          "active",
		}
		sessionData, err := json.Marshal(session)
		if err != nil {
			return err
		}
		if err := gsb.Put([]byte(session.ID), sessionData); err != nil {
			return err
		}
		key := fmt.Sprintf("%s:%s:%d", userID, session.ID, session.StartedAt)
		return gsu.Put([]byte(key), nil)
	})
	return &session, err
}

func (s *Store) EndGameSession(sessionID string) (*models.GameSession, error) {
	var session models.GameSession
	err := s.db.Update(func(tx *bolt.Tx) error {
		gsb := tx.Bucket(bucketGameSessions)
		data := gsb.Get([]byte(sessionID))
		if data == nil {
			return fmt.Errorf("session not found: %s", sessionID)
		}
		if err := json.Unmarshal(data, &session); err != nil {
			return err
		}
		if session.Status != "active" {
			return fmt.Errorf("session already ended")
		}
		session.EndedAt = time.Now().Unix()
		session.Duration = session.EndedAt - session.StartedAt
		session.Status = "ended"
		session.Cost = 0
		updated, err := json.Marshal(session)
		if err != nil {
			return err
		}
		return gsb.Put([]byte(sessionID), updated)
	})
	return &session, err
}

func (s *Store) ListActiveGameSessions() ([]models.GameSession, error) {
	var sessions []models.GameSession
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketGameSessions)
		return b.ForEach(func(k, v []byte) error {
			var session models.GameSession
			if err := json.Unmarshal(v, &session); err != nil {
				return err
			}
			if session.Status == "active" {
				sessions = append(sessions, session)
			}
			return nil
		})
	})
	if sessions == nil {
		sessions = []models.GameSession{}
	}
	return sessions, err
}

func (s *Store) ListGameSessionsByUser(userID string) ([]models.GameSession, error) {
	var sessions []models.GameSession
	err := s.db.View(func(tx *bolt.Tx) error {
		gsb := tx.Bucket(bucketGameSessions)
		gsu := tx.Bucket(bucketGameSessionsByUser)
		prefix := []byte(userID + ":")
		c := gsu.Cursor()
		for k, _ := c.Seek(prefix); k != nil && hasBytePrefix(k, prefix); k, _ = c.Next() {
			parts := splitCompositeKey(string(k))
			if len(parts) >= 2 {
				data := gsb.Get([]byte(parts[1]))
				if data != nil {
					var session models.GameSession
					if err := json.Unmarshal(data, &session); err != nil {
						return err
					}
					sessions = append(sessions, session)
				}
			}
		}
		return nil
	})
	if sessions == nil {
		sessions = []models.GameSession{}
	}
	return sessions, err
}

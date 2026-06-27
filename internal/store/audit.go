package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thelastvideostore/internal/models"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) AppendAuditEntry(entry *models.AuditEntry) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketAuditLogs)

		if entry.Timestamp == 0 {
			entry.Timestamp = time.Now().Unix()
		}

		data, err := encode(entry)
		if err != nil {
			return fmt.Errorf("append audit: %w", err)
		}

		return b.Put([]byte(entry.ID), data)
	})
}

func (s *Store) GetAllAuditEntries() ([]*models.AuditEntry, error) {
	var entries []*models.AuditEntry
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketAuditLogs)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var entry models.AuditEntry
			if err := json.Unmarshal(v, &entry); err != nil {
				continue
			}
			entries = append(entries, &entry)
		}
		return nil
	})
	return entries, err
}

func (s *Store) GetAuditEntriesByUser(userID string) ([]*models.AuditEntry, error) {
	var entries []*models.AuditEntry
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketAuditLogs)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var entry models.AuditEntry
			if err := json.Unmarshal(v, &entry); err != nil {
				continue
			}
			if entry.ActorID == userID || entry.TargetID == userID {
				entries = append(entries, &entry)
			}
		}
		return nil
	})
	return entries, err
}

// UpdateAuditEntry overwrites the stored JSON for the given entry ID.
// This is normally never called — the audit chain is append-only — but
// it is exposed so the demo tamper tool (cmd/tamper) can simulate an
// attacker modifying a row directly in the DB.
func (s *Store) UpdateAuditEntry(entry *models.AuditEntry) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketAuditLogs)
		data, err := encode(entry)
		if err != nil {
			return fmt.Errorf("update audit: %w", err)
		}
		return b.Put([]byte(entry.ID), data)
	})
}

package auth

import (
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

func AppendAuditEntry(s *store.Store, hc *crypto.HashChain, action, actorID, targetID, data string) error {
	entry := hc.Append(action, actorID, targetID, data)

	auditEntry := &models.AuditEntry{
		ID:        uuid.NewString(),
		Timestamp: entry.Timestamp,
		Action:    action,
		ActorID:   actorID,
		TargetID:  targetID,
		Data:      data,
		Hash:      entry.Hash,
		PrevHash:  entry.PrevHash,
	}

	return s.AppendAuditEntry(auditEntry)
}

func VerifyAuditChain(s *store.Store) (bool, error) {
	entries, err := s.GetAllAuditEntries()
	if err != nil {
		return false, fmt.Errorf("verify audit chain: %w", err)
	}

	if len(entries) == 0 {
		return true, nil
	}

	chainEntries := make([]crypto.HashChainEntry, len(entries))
	for i, e := range entries {
		chainEntries[i] = crypto.HashChainEntry{
			Timestamp: e.Timestamp,
			Action:    e.Action,
			ActorID:   e.ActorID,
			TargetID:  e.TargetID,
			Data:      e.Data,
			Hash:      e.Hash,
			PrevHash:  e.PrevHash,
		}
	}

	valid := crypto.VerifyChain(chainEntries)
	if !valid {
		return false, nil
	}

	prev := []byte("GENESIS")
	for _, e := range entries {
		if hex.EncodeToString(e.PrevHash) != hex.EncodeToString(prev) {
			return false, nil
		}
		prev = e.Hash
	}

	return true, nil
}

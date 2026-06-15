package auth

import (
	"fmt"
	"sort"

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

	sort.Slice(entries, func(i, j int) bool { return entries[i].Timestamp < entries[j].Timestamp })

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

	return crypto.VerifyChain(chainEntries), nil
}

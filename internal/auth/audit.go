package auth

import (
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

type ChainVerifyResult struct {
	Valid    bool
	BrokenAt int
	BrokenID string
	Reason   string
	Entries  []*models.AuditEntry
}

func VerifyAuditChain(s *store.Store) (ChainVerifyResult, error) {
	entries, err := s.GetAllAuditEntries()
	if err != nil {
		return ChainVerifyResult{
			Valid:    false,
			BrokenAt: -1,
			Reason:   "verify audit chain: " + err.Error(),
		}, err
	}

	if len(entries) == 0 {
		return ChainVerifyResult{Valid: true, BrokenAt: -1, Entries: entries}, nil
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

	valid, chainErr := crypto.VerifyChain(chainEntries)
	if valid {
		return ChainVerifyResult{Valid: true, BrokenAt: -1, Entries: entries}, nil
	}

	result := ChainVerifyResult{
		Valid:    false,
		BrokenAt: chainErr.BrokenAt,
		Reason:   chainErr.Reason,
		Entries:  entries,
	}
	if chainErr.BrokenAt >= 0 && chainErr.BrokenAt < len(entries) {
		result.BrokenID = entries[chainErr.BrokenAt].ID
	}
	return result, nil
}

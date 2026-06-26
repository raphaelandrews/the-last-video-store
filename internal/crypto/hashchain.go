package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HashChainEntry struct {
	Timestamp int64
	Action    string
	ActorID   string
	TargetID  string
	Data      string
	Hash      []byte
	PrevHash  []byte
}

type HashChain struct {
	entries  []HashChainEntry
	lastHash []byte
}

func New() *HashChain {
	return &HashChain{
		lastHash: []byte("GENESIS"),
	}
}

func (hc *HashChain) Append(action, actorID, targetID, data string) HashChainEntry {
	now := int64(0)
	entry := HashChainEntry{
		Timestamp: now,
		Action:    action,
		ActorID:   actorID,
		TargetID:  targetID,
		Data:      data,
		PrevHash:  make([]byte, len(hc.lastHash)),
	}
	copy(entry.PrevHash, hc.lastHash)

	entry.Hash = hc.computeHash(entry)
	hc.entries = append(hc.entries, entry)
	hc.lastHash = entry.Hash
	return entry
}

func (hc *HashChain) Verify() bool {
	prev := []byte("GENESIS")
	for i, entry := range hc.entries {
		if hex.EncodeToString(entry.PrevHash) != hex.EncodeToString(prev) {
			return false
		}
		recomputed := hc.computeHash(hc.entries[i])
		if hex.EncodeToString(recomputed) != hex.EncodeToString(entry.Hash) {
			return false
		}
		prev = entry.Hash
	}
	return true
}

func (hc *HashChain) GetAll() []HashChainEntry {
	return hc.entries
}

func (hc *HashChain) Len() int {
	return len(hc.entries)
}

func (hc *HashChain) computeHash(entry HashChainEntry) []byte {
	input := fmt.Sprintf("%x|%d|%s|%s|%s|%s",
		entry.PrevHash, entry.Timestamp, entry.Action,
		entry.ActorID, entry.TargetID, entry.Data)
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

// ChainError describes where the chain broke. A non-nil error means the
// chain is invalid; the int field is the index of the first entry whose
// hash did not match (or whose PrevHash did not match the previous
// entry's Hash). Index 0 is the first entry after GENESIS.
type ChainError struct {
	BrokenAt int // 0-based index of first broken entry; -1 if chain is empty
	Reason   string
}

func (e *ChainError) Error() string {
	if e == nil {
		return ""
	}
	if e.BrokenAt >= 0 {
		return fmt.Sprintf("chain broken at entry %d: %s", e.BrokenAt, e.Reason)
	}
	return "chain broken: " + e.Reason
}

// VerifyChain returns nil on success, or a *ChainError describing the
// first entry whose hash did not match. It returns true as a convenience
// for callers that only want the boolean; callers that need the broken
// index should check for a non-nil error.
func VerifyChain(entries []HashChainEntry) (bool, *ChainError) {
	prev := []byte("GENESIS")
	for i, entry := range entries {
		if hex.EncodeToString(entry.PrevHash) != hex.EncodeToString(prev) {
			return false, &ChainError{
				BrokenAt: i,
				Reason:   "PrevHash does not match the previous entry's Hash",
			}
		}
		input := fmt.Sprintf("%x|%d|%s|%s|%s|%s",
			entry.PrevHash, entry.Timestamp, entry.Action,
			entry.ActorID, entry.TargetID, entry.Data)
		hash := sha256.Sum256([]byte(input))
		if hex.EncodeToString(hash[:]) != hex.EncodeToString(entry.Hash) {
			return false, &ChainError{
				BrokenAt: i,
				Reason:   "Hash does not match recomputed value (data was modified)",
			}
		}
		prev = entry.Hash
	}
	return true, nil
}

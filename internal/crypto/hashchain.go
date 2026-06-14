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

func VerifyChain(entries []HashChainEntry) bool {
	prev := []byte("GENESIS")
	for _, entry := range entries {
		if hex.EncodeToString(entry.PrevHash) != hex.EncodeToString(prev) {
			return false
		}
		input := fmt.Sprintf("%x|%d|%s|%s|%s|%s",
			entry.PrevHash, entry.Timestamp, entry.Action,
			entry.ActorID, entry.TargetID, entry.Data)
		hash := sha256.Sum256([]byte(input))
		if hex.EncodeToString(hash[:]) != hex.EncodeToString(entry.Hash) {
			return false
		}
		prev = entry.Hash
	}
	return true
}

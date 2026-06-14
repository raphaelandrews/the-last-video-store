package crypto

import (
	"encoding/hex"
	"testing"
)

func TestGenesisEntry(t *testing.T) {
	hc := New()
	entry := hc.Append("ActionLogin", "user-1", "", "test login")

	if hex.EncodeToString(entry.PrevHash) != hex.EncodeToString([]byte("GENESIS")) {
		t.Errorf("first entry PrevHash should be GENESIS, got %x", entry.PrevHash)
	}
	if entry.Hash == nil {
		t.Error("Hash should not be nil")
	}
	if hc.Len() != 1 {
		t.Errorf("Len = %d, want 1", hc.Len())
	}
}

func TestChainLinking(t *testing.T) {
	hc := New()
	e1 := hc.Append("ActionLogin", "user-1", "", "")
	e2 := hc.Append("ActionRent", "user-1", "movie-1", "")

	if hex.EncodeToString(e2.PrevHash) != hex.EncodeToString(e1.Hash) {
		t.Error("second entry PrevHash should equal first entry Hash")
	}
}

func TestVerifyIntact(t *testing.T) {
	hc := New()
	hc.Append("ActionLogin", "user-1", "", "1")
	hc.Append("ActionRent", "user-1", "movie-1", "2")
	hc.Append("ActionReturn", "user-1", "movie-1", "3")

	if !hc.Verify() {
		t.Error("intact chain should verify as true")
	}
}

func TestVerifyTampered(t *testing.T) {
	hc := New()
	hc.Append("ActionLogin", "user-1", "", "1")
	hc.Append("ActionRent", "user-1", "movie-1", "2")
	hc.Append("ActionReturn", "user-1", "movie-1", "3")

	entries := hc.GetAll()
	entries[1].Data = "TAMPERED"

	if VerifyChain(entries) {
		t.Error("tampered chain should not verify")
	}
}

func TestVerifyTamperedHash(t *testing.T) {
	hc := New()
	hc.Append("ActionLogin", "user-1", "", "1")
	hc.Append("ActionRent", "user-1", "movie-1", "2")

	entries := hc.GetAll()
	entries[0].Hash[0] ^= 0xFF

	if VerifyChain(entries) {
		t.Error("chain with corrupted hash should not verify")
	}
}

func TestVerifyBrokenLink(t *testing.T) {
	hc := New()
	hc.Append("ActionLogin", "user-1", "", "1")
	hc.Append("ActionRent", "user-1", "movie-1", "2")

	entries := hc.GetAll()
	entries[1].PrevHash[0] ^= 0xFF

	if VerifyChain(entries) {
		t.Error("chain with broken link should not verify")
	}
}

func TestEmptyChain(t *testing.T) {
	hc := New()
	if !hc.Verify() {
		t.Error("empty chain should verify as true")
	}
	if hc.Len() != 0 {
		t.Errorf("empty chain Len = %d, want 0", hc.Len())
	}
}

func TestGetAll(t *testing.T) {
	hc := New()
	hc.Append("Action1", "u1", "t1", "d1")
	hc.Append("Action2", "u2", "t2", "d2")

	entries := hc.GetAll()
	if len(entries) != 2 {
		t.Errorf("GetAll len = %d, want 2", len(entries))
	}
	if entries[0].Action != "Action1" {
		t.Errorf("first entry Action = %s, want Action1", entries[0].Action)
	}
}

func BenchmarkAppend(b *testing.B) {
	hc := New()
	for b.Loop() {
		hc.Append("ActionRent", "user-1", "movie-1", "benchmark data")
	}
}

func BenchmarkVerify(b *testing.B) {
	hc := New()
	for range 100 {
		hc.Append("ActionRent", "user-1", "movie-1", "benchmark data")
	}
	b.ResetTimer()
	for b.Loop() {
		hc.Verify()
	}
}

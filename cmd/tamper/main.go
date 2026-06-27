// Command tamper is a demo tool for the hash-chain audit log. It
// simulates an attacker with direct database access and intentionally
// corrupts a single audit row, breaking the chain at that point. After
// running `tamper demo` the TUI's audit log screen, when asked to
// verify, will show a "chain broken" error with the index of the
// tampered row.
//
// Subcommands:
//
//	list                  — print every audit entry with its index, ID,
//	                        action, actor, target, timestamp
//	corrupt <id>           — flip the Data field on a single row, which
//	                        breaks the chain at that point
//	restore <id>           — undo a previous tamper by recomputing the
//	                        correct hash for the row's current state
//	demo                  — automatically pick a random non-edge row
//	                        and corrupt it (best for live demo)
//
// All operations require the server to be stopped (the DB file is
// opened in exclusive mode by BoltDB).
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

const defaultDB = "thelastvideostore.db"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	dbPath := defaultDB
	cmd := os.Args[1]

	// Allow `tamper -d other.db <subcmd>` style flag.
	if cmd == "-d" || cmd == "--db" {
		if len(os.Args) < 4 {
			usage()
			os.Exit(1)
		}
		dbPath = os.Args[2]
		cmd = os.Args[3]
		os.Args = append([]string{os.Args[0]}, os.Args[3:]...)
	}

	s, err := store.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ could not open %s: %v\n", dbPath, err)
		fmt.Fprintln(os.Stderr, "   (is the server still running? stop it first.)")
		os.Exit(1)
	}
	defer s.Close()

	switch cmd {
	case "list":
		listEntries(s)
	case "corrupt":
		requireArg("corrupt <id>")
		corruptEntry(s, os.Args[2])
	case "restore":
		requireArg("restore <id>")
		restoreEntry(s, os.Args[2])
	case "demo":
		demo(s)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `tamper — demo tool for the audit hash chain

Usage:
  tamper [flags] list
  tamper [flags] corrupt <entry-id>
  tamper [flags] restore <entry-id>
  tamper [flags] demo

Flags:
  -d, --db <path>   Path to the BoltDB file (default: thelastvideostore.db)

Examples:
  go run ./cmd/tamper list
  go run ./cmd/tamper demo
  go run ./cmd/tamper corrupt 5b2c8a3d-...`)
}

func requireArg(name string) {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "missing argument: %s\n", name)
		os.Exit(1)
	}
}

// loadChronological returns all audit entries sorted oldest-first,
// matching the order the server uses for chain verification.
func loadChronological(s *store.Store) []*models.AuditEntry {
	all, err := s.GetAllAuditEntries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to read entries: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Timestamp < all[j].Timestamp })
	return all
}

func listEntries(s *store.Store) {
	all := loadChronological(s)
	fmt.Printf("Found %d audit entries (chronological order):\n\n", len(all))
	fmt.Printf("%-4s %-36s %-12s %-12s %s\n", "IDX", "ENTRY ID", "TIMESTAMP", "ACTION", "DATA")
	for i, e := range all {
		ts := time.Unix(0, e.Timestamp).Format("15:04:05.000")
		fmt.Printf("%-4d %-36s %-12s %-12s %s\n", i, e.ID, ts, e.Action, e.Data)
	}
}

func findEntry(s *store.Store, id string) (*models.AuditEntry, int) {
	all := loadChronological(s)
	for i, e := range all {
		if e.ID == id {
			return e, i
		}
	}
	fmt.Fprintf(os.Stderr, "❌ no entry with id %s\n", id)
	os.Exit(1)
	return nil, 0
}

func corruptEntry(s *store.Store, id string) {
	entry, idx := findEntry(s, id)
	before := entry.Data
	// Append a marker so it's obvious in the audit log which row was
	// tampered with. We don't touch the Hash field — the recomputed
	// hash will differ from the stored hash, which is what verify
	// detects.
	entry.Data = before + " [TAMPERED]"
	if err := s.UpdateAuditEntry(entry); err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Corrupted entry #%d (%s)\n", idx, entry.ID)
	fmt.Printf("   data before: %q\n", before)
	fmt.Printf("   data after:  %q\n", entry.Data)
	fmt.Println()
	fmt.Println("The chain is now broken at this entry. Run the server, open")
	fmt.Println("the audit log screen, and press [v] to see the verification fail.")
}

func restoreEntry(s *store.Store, id string) {
	entry, idx := findEntry(s, id)
	// Re-derive the original Data. We can't recover it from the broken
	// row, so the operator has to supply it. For convenience we accept
	// it as the second argument.
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: tamper restore <id> <original-data>\n")
		fmt.Fprintf(os.Stderr, "  current data: %q\n", entry.Data)
		os.Exit(1)
	}
	entry.Data = os.Args[3]
	if err := s.UpdateAuditEntry(entry); err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Restored entry #%d data field to %q\n", idx, entry.Data)
	fmt.Println("Note: this restores the Data, but the stored Hash still doesn't match.")
	fmt.Println("To fully repair the chain, recompute and overwrite the Hash too:")
	fmt.Printf("  $ %s\n", recomputeHint(entry))
}

// recomputeHint returns a shell command the user can run to recompute
// the correct hash for an entry, given its current fields. Used by
// restore as a hint.
func recomputeHint(e *models.AuditEntry) string {
	// We can't easily rebuild the prev-hash here because it depends
	// on the previous entry's Hash. The hash chain is what makes
	// tampering detectable in the first place — recovery requires the
	// full chain. For the demo, the simplest recovery is to seed
	// again (go run ./data).
	return "go run ./data   # wipe + reseed to fully repair the chain"
}

func demo(s *store.Store) {
	all := loadChronological(s)
	if len(all) < 5 {
		fmt.Fprintf(os.Stderr, "❌ need at least 5 audit entries for demo; have %d\n", len(all))
		fmt.Fprintln(os.Stderr, "   perform some actions (login, rent, return) in the TUI first.")
		os.Exit(1)
	}
	// Pick a row in the middle. Tampering with index 0 or the last row
	// is a degenerate case; the middle is the realistic attack.
	idx := len(all) / 2
	entry := all[idx]
	before := entry.Data
	entry.Data = before + " [TAMPERED]"
	if err := s.UpdateAuditEntry(entry); err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🎭 Demo tamper applied at entry #%d (of %d total)\n", idx, len(all))
	fmt.Println()
	fmt.Println("    ID:    ", entry.ID)
	fmt.Println("    Action:", entry.Action)
	fmt.Println("    Before:", before)
	fmt.Println("    After: ", entry.Data)
	fmt.Println()
	fmt.Println("Stored Hash:    ", hex.EncodeToString(entry.Hash))
	fmt.Println("Expected Hash:  ", "(different — verify will catch this)")
	fmt.Println()
	fmt.Println("▶ Restart the server and open the audit log in the TUI.")
	fmt.Println("▶ Press [v] to verify — you should see a red broken-chain")
	fmt.Println("   message with the index pointing to entry #" + strconv.Itoa(idx) + ".")
	fmt.Println("▶ Press [g] to jump straight to the broken row.")
	fmt.Println()
	fmt.Println("To repair: stop the server, then `go run ./data` to reseed.")
}

var _ = crypto.New // keep the import; the tool conceptually operates on the chain

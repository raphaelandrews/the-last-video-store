package store

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	bucketUsers                = []byte("users")
	bucketUsersByUsername      = []byte("users_by_username")
	bucketMovies               = []byte("movies")
	bucketMoviesByGenre        = []byte("movies_by_genre")
	bucketMoviesByTitle        = []byte("movies_by_title")
	bucketRentals              = []byte("rentals")
	bucketRentalsByUser        = []byte("rentals_by_user")
	bucketAuditLogs            = []byte("audit_logs")
	bucketSessions             = []byte("sessions")
	bucketBanned               = []byte("banned")
	bucketWishlists            = []byte("wishlists")
	bucketStaffPicks           = []byte("staff_picks")
	bucketTOTPSecrets          = []byte("totp_secrets")
	bucketLockouts             = []byte("lockouts")
	bucketMerch                = []byte("merch")
	bucketInventory            = []byte("inventory")
	bucketSnackBarItems        = []byte("snackbar_items")
	bucketSnackBarOrders       = []byte("snackbar_orders")
	bucketSnackBarOrdersByUser = []byte("snackbar_orders_by_user")
	bucketGameSessions         = []byte("game_sessions")
	bucketGameSessionsByUser   = []byte("game_sessions_by_user")
)

type Store struct {
	db *bolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			bucketUsers, bucketUsersByUsername, bucketMovies,
			bucketMoviesByGenre, bucketMoviesByTitle, bucketRentals,
			bucketRentalsByUser, bucketAuditLogs, bucketSessions,
			bucketBanned, bucketWishlists, bucketStaffPicks,
			bucketTOTPSecrets, bucketLockouts, bucketMerch, bucketInventory,
			bucketSnackBarItems, bucketSnackBarOrders, bucketSnackBarOrdersByUser,
			bucketGameSessions, bucketGameSessionsByUser,
		}
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return fmt.Errorf("store: create bucket %s: %w", b, err)
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DB() *bolt.DB {
	return s.db
}

func encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func keyUID(prefix string, id string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefix, id))
}

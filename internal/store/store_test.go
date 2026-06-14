package store

import (
	"os"
	"testing"
	"time"

	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
)

func setupStore(t *testing.T) *Store {
	t.Helper()
	path := t.TempDir() + "/test.db"
	store, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { store.Close(); os.Remove(path) })
	return store
}

func TestUserCRUD(t *testing.T) {
	s := setupStore(t)

	user := &models.User{
		ID:          "user-1",
		Username:    "testuser",
		Tier:        bitmask.TierGold,
		MaxRentals:  5,
		RentalCount: 0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := s.CreateUser(user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	found, err := s.GetUserByID("user-1")
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if found.Username != "testuser" {
		t.Errorf("username = %q, want testuser", found.Username)
	}
	if found.Tier != bitmask.TierGold {
		t.Errorf("tier = %v, want TierGold", found.Tier)
	}

	found, err = s.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if found.ID != "user-1" {
		t.Errorf("id = %q, want user-1", found.ID)
	}

	if !s.UserExists("testuser") {
		t.Error("UserExists should be true")
	}
	if s.UserExists("nonexistent") {
		t.Error("UserExists should be false")
	}

	user.Username = "renamed"
	user.Tier = bitmask.TierManager
	user.MaxRentals = 10
	if err := s.UpdateUser(user); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}

	found, err = s.GetUserByID("user-1")
	if err != nil {
		t.Fatalf("GetUserByID after update: %v", err)
	}
	if found.Username != "renamed" {
		t.Errorf("username = %q, want renamed", found.Username)
	}
	if found.Tier != bitmask.TierManager {
		t.Errorf("tier = %v, want TierManager", found.Tier)
	}

	users, err := s.ListUsers()
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("ListUsers len = %d, want 1", len(users))
	}

	if err := s.DeleteUser("user-1"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	_, err = s.GetUserByID("user-1")
	if err == nil {
		t.Error("GetUserByID should fail after delete")
	}
}

func TestMovieCRUD(t *testing.T) {
	s := setupStore(t)

	movie := &models.Movie{
		ID:              "movie-1",
		Title:           "The Matrix",
		Year:            1999,
		Genre:           "SciFi",
		Format:          "DVD",
		Director:        "Wachowski",
		Cast:            []string{"Keanu Reeves", "Laurence Fishburne"},
		Synopsis:        "A computer hacker learns about the true nature of reality.",
		Rating:          4.7,
		RatingCount:     5000,
		Available:       true,
		CopiesTotal:     5,
		CopiesAvailable: 3,
		IsNewRelease:    false,
		CreatedAt:       time.Now().Unix(),
	}

	if err := s.CreateMovie(movie); err != nil {
		t.Fatalf("CreateMovie: %v", err)
	}

	found, err := s.GetMovieByID("movie-1")
	if err != nil {
		t.Fatalf("GetMovieByID: %v", err)
	}
	if found.Title != "The Matrix" {
		t.Errorf("title = %q, want The Matrix", found.Title)
	}

	movies, total, err := s.ListMovies("SciFi", 0, 10)
	if err != nil {
		t.Fatalf("ListMovies by genre: %v", err)
	}
	if total != 1 || len(movies) != 1 {
		t.Errorf("ListMovies(SciFi) = %d/%d, want 1/1", len(movies), total)
	}

	results, err := s.SearchMoviesByPrefix("the", 10)
	if err != nil {
		t.Fatalf("SearchMoviesByPrefix: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchMoviesByPrefix(the) = %d, want 1", len(results))
	}

	newReleases, err := s.GetNewReleases()
	if err != nil {
		t.Fatalf("GetNewReleases: %v", err)
	}
	if len(newReleases) != 0 {
		t.Errorf("GetNewReleases = %d, want 0", len(newReleases))
	}

	movie.Title = "The Matrix Reloaded"
	if err := s.UpdateMovie(movie); err != nil {
		t.Fatalf("UpdateMovie: %v", err)
	}

	results, err = s.SearchMoviesByPrefix("the matrix reloaded", 10)
	if err != nil {
		t.Fatalf("SearchMoviesByPrefix after update: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchMoviesByPrefix after rename = %d, want 1", len(results))
	}

	if err := s.DeleteMovie("movie-1"); err != nil {
		t.Fatalf("DeleteMovie: %v", err)
	}
	_, err = s.GetMovieByID("movie-1")
	if err == nil {
		t.Error("GetMovieByID should fail after delete")
	}
}

func TestRentalCRUD(t *testing.T) {
	s := setupStore(t)

	rental := &models.Rental{
		ID:          "rental-1",
		UserID:      "user-1",
		MovieID:     "movie-1",
		MovieFormat: "DVD",
		RentedAt:    time.Now().Unix(),
		DueDate:     time.Now().Unix() + int64(5*24*time.Hour),
		Status:      models.RentalActive,
	}

	if err := s.CreateRental(rental); err != nil {
		t.Fatalf("CreateRental: %v", err)
	}

	found, err := s.GetRentalByID("rental-1")
	if err != nil {
		t.Fatalf("GetRentalByID: %v", err)
	}
	if found.Status != models.RentalActive {
		t.Errorf("status = %q, want active", found.Status)
	}

	active, err := s.GetActiveRentalsByUser("user-1")
	if err != nil {
		t.Fatalf("GetActiveRentalsByUser: %v", err)
	}
	if len(active) != 1 {
		t.Errorf("active rentals = %d, want 1", len(active))
	}

	count, err := s.CountActiveRentalsByUser("user-1")
	if err != nil {
		t.Fatalf("CountActiveRentalsByUser: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	rental.Status = models.RentalReturned
	rental.ReturnedAt = time.Now().Unix()
	if err := s.UpdateRental(rental); err != nil {
		t.Fatalf("UpdateRental: %v", err)
	}

	count, err = s.CountActiveRentalsByUser("user-1")
	if err != nil {
		t.Fatalf("CountActiveRentalsByUser after return: %v", err)
	}
	if count != 0 {
		t.Errorf("count after return = %d, want 0", count)
	}

	history, err := s.GetRentalHistoryByUser("user-1")
	if err != nil {
		t.Fatalf("GetRentalHistoryByUser: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("history = %d, want 1", len(history))
	}
}

func TestOverdueRentals(t *testing.T) {
	s := setupStore(t)

	overdue := &models.Rental{
		ID:          "rental-overdue",
		UserID:      "user-1",
		MovieID:     "movie-1",
		MovieFormat: "VHS",
		RentedAt:    time.Now().Unix() - 5*24*3600,
		DueDate:     time.Now().Unix() - 2*24*3600,
		Status:      models.RentalActive,
	}
	s.CreateRental(overdue)

	results, err := s.GetOverdueRentals()
	if err != nil {
		t.Fatalf("GetOverdueRentals: %v", err)
	}
	if len(results) < 1 {
		t.Errorf("overdue rentals = %d, want >= 1", len(results))
	}
}

func TestSessionOperations(t *testing.T) {
	s := setupStore(t)
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	if err := s.SaveRefreshToken("user-1", "token-1", expiresAt); err != nil {
		t.Fatalf("SaveRefreshToken: %v", err)
	}

	valid, err := s.ValidateRefreshToken("user-1", "token-1")
	if err != nil {
		t.Fatalf("ValidateRefreshToken: %v", err)
	}
	if !valid {
		t.Error("token should be valid")
	}

	valid, err = s.ValidateRefreshToken("user-1", "token-nonexistent")
	if err != nil {
		t.Fatalf("ValidateRefreshToken missing: %v", err)
	}
	if valid {
		t.Error("missing token should not be valid")
	}

	if err := s.InvalidateRefreshToken("user-1", "token-1"); err != nil {
		t.Fatalf("InvalidateRefreshToken: %v", err)
	}

	valid, err = s.ValidateRefreshToken("user-1", "token-1")
	if err != nil {
		t.Fatalf("ValidateRefreshToken after invalidation: %v", err)
	}
	if valid {
		t.Error("invalidated token should not be valid")
	}

	s.SaveRefreshToken("user-1", "tok1", expiresAt)
	s.SaveRefreshToken("user-1", "tok2", expiresAt)
	if err := s.InvalidateAllUserSessions("user-1"); err != nil {
		t.Fatalf("InvalidateAllUserSessions: %v", err)
	}
	valid, _ = s.ValidateRefreshToken("user-1", "tok1")
	if valid {
		t.Error("all sessions should be invalidated")
	}
}

func TestLockoutOperations(t *testing.T) {
	s := setupStore(t)

	count, err := s.IncrementFailedAttempts("testuser")
	if err != nil {
		t.Fatalf("first IncrementFailedAttempts: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	count, err = s.IncrementFailedAttempts("testuser")
	if err != nil {
		t.Fatalf("second IncrementFailedAttempts: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	if err := s.ResetFailedAttempts("testuser"); err != nil {
		t.Fatalf("ResetFailedAttempts: %v", err)
	}

	count, err = s.IncrementFailedAttempts("testuser")
	if err != nil {
		t.Fatalf("after reset: %v", err)
	}
	if count != 1 {
		t.Errorf("count after reset = %d, want 1", count)
	}

	locked, err := s.IsUserLocked("testuser")
	if err != nil {
		t.Fatalf("IsUserLocked: %v", err)
	}
	if locked {
		t.Error("user should not be locked yet")
	}

	futureTime := time.Now().Add(30 * time.Minute).Unix()
	if err := s.LockUserUntil("testuser", futureTime); err != nil {
		t.Fatalf("LockUserUntil: %v", err)
	}

	locked, err = s.IsUserLocked("testuser")
	if err != nil {
		t.Fatalf("IsUserLocked after lock: %v", err)
	}
	if !locked {
		t.Error("user should be locked")
	}
}

func TestWishlistCRUD(t *testing.T) {
	s := setupStore(t)

	if err := s.AddToWishlist("user-1", "movie-1"); err != nil {
		t.Fatalf("AddToWishlist: %v", err)
	}
	if err := s.AddToWishlist("user-1", "movie-2"); err != nil {
		t.Fatalf("AddToWishlist 2: %v", err)
	}

	items, err := s.GetWishlist("user-1")
	if err != nil {
		t.Fatalf("GetWishlist: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("wishlist length = %d, want 2", len(items))
	}

	inWishlist, err := s.IsInWishlist("user-1", "movie-1")
	if err != nil {
		t.Fatalf("IsInWishlist: %v", err)
	}
	if !inWishlist {
		t.Error("movie-1 should be in wishlist")
	}

	size, err := s.GetWishlistSize("user-1")
	if err != nil {
		t.Fatalf("GetWishlistSize: %v", err)
	}
	if size != 2 {
		t.Errorf("wishlist size = %d, want 2", size)
	}

	if err := s.RemoveFromWishlist("user-1", "movie-1"); err != nil {
		t.Fatalf("RemoveFromWishlist: %v", err)
	}

	inWishlist, _ = s.IsInWishlist("user-1", "movie-1")
	if inWishlist {
		t.Error("movie-1 should no longer be in wishlist")
	}
}

func TestStaffPicks(t *testing.T) {
	s := setupStore(t)

	if err := s.AddStaffPick("movie-1"); err != nil {
		t.Fatalf("AddStaffPick: %v", err)
	}
	if err := s.AddStaffPick("movie-2"); err != nil {
		t.Fatalf("AddStaffPick 2: %v", err)
	}

	if !s.IsStaffPick("movie-1") {
		t.Error("movie-1 should be a staff pick")
	}

	picks, err := s.GetStaffPicks()
	if err != nil {
		t.Fatalf("GetStaffPicks: %v", err)
	}
	if len(picks) != 2 {
		t.Errorf("staff picks = %d, want 2", len(picks))
	}

	if err := s.RemoveStaffPick("movie-1"); err != nil {
		t.Fatalf("RemoveStaffPick: %v", err)
	}
	if s.IsStaffPick("movie-1") {
		t.Error("movie-1 should no longer be a staff pick")
	}
}

func TestLastChanceMovies(t *testing.T) {
	s := setupStore(t)

	lastChance := &models.Movie{
		ID:              "lc-1",
		Title:           "Last Chance Movie",
		Year:            1994,
		Genre:           "Drama",
		Format:          "VHS",
		CopiesTotal:     1,
		CopiesAvailable: 1,
		IsNewRelease:    false,
	}
	s.CreateMovie(lastChance)

	results, err := s.GetLastChanceMovies()
	if err != nil {
		t.Fatalf("GetLastChanceMovies: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("last chance = %d, want 1", len(results))
	}
}

func TestTOTPStore(t *testing.T) {
	s := setupStore(t)

	if err := s.SaveTOTPSecret("user-1", []byte("test-secret")); err != nil {
		t.Fatalf("SaveTOTPSecret: %v", err)
	}

	secret, err := s.GetTOTPSecret("user-1")
	if err != nil {
		t.Fatalf("GetTOTPSecret: %v", err)
	}
	if string(secret) != "test-secret" {
		t.Errorf("secret = %q, want test-secret", secret)
	}

	count, err := s.IncrementTOTPFailures("user-1")
	if err != nil {
		t.Fatalf("IncrementTOTPFailures: %v", err)
	}
	if count != 1 {
		t.Errorf("totp failures = %d, want 1", count)
	}

	if err := s.ResetTOTPFailures("user-1"); err != nil {
		t.Fatalf("ResetTOTPFailures: %v", err)
	}

	count, err = s.IncrementTOTPFailures("user-1")
	if err != nil {
		t.Fatalf("IncrementTOTPFailures after reset: %v", err)
	}
	if count != 1 {
		t.Errorf("totp failures after reset = %d, want 1", count)
	}

	futureTime := time.Now().Add(10 * time.Minute).Unix()
	if err := s.LockTOTPUserUntil("user-1", futureTime); err != nil {
		t.Fatalf("LockTOTPUserUntil: %v", err)
	}

	locked, err := s.IsTOTPLocked("user-1")
	if err != nil {
		t.Fatalf("IsTOTPLocked: %v", err)
	}
	if !locked {
		t.Error("TOTP should be locked")
	}

	if err := s.DeleteTOTPSecret("user-1"); err != nil {
		t.Fatalf("DeleteTOTPSecret: %v", err)
	}

	secret, err = s.GetTOTPSecret("user-1")
	if err != nil {
		t.Fatalf("GetTOTPSecret after delete: %v", err)
	}
	if len(secret) != 0 {
		t.Error("secret should be empty after delete")
	}
}

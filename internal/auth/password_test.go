package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("securepassword123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" {
		t.Error("hash should not be empty")
	}
	if hash == "securepassword123" {
		t.Error("hash should not be the plaintext password")
	}
}

func TestHashPasswordDifferentEachTime(t *testing.T) {
	h1, _ := HashPassword("samepassword")
	h2, _ := HashPassword("samepassword")
	if h1 == h2 {
		t.Error("two hashes of the same password should differ (salt)")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, _ := HashPassword("mypassword")

	if !CheckPassword(hash, "mypassword") {
		t.Error("correct password should verify")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Error("wrong password should not verify")
	}
}

func TestEmptyPassword(t *testing.T) {
	_, err := HashPassword("")
	if err == nil {
		t.Error("empty password should return error")
	}

	if CheckPassword("", "test") {
		t.Error("empty hash should not verify")
	}
	if CheckPassword("somehash", "") {
		t.Error("empty password should not verify")
	}
}

func TestBcryptCost(t *testing.T) {
	hash, _ := HashPassword("test")
	if len(hash) < 50 {
		t.Error("bcrypt hash should be sufficiently long")
	}
}

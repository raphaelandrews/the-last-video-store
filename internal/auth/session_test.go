package auth

import (
	"testing"
	"time"

	"github.com/thelastvideostore/internal/ds/bitmask"
)

func TestGenerateTokenPair(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	pair, err := GenerateTokenPair("user-1", uint16(bitmask.TierGold), secret)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	if pair.AccessToken == "" {
		t.Error("access token should not be empty")
	}
	if pair.RefreshToken == "" {
		t.Error("refresh token should not be empty")
	}
	if pair.ExpiresAt <= time.Now().Unix() {
		t.Error("expires_at should be in the future")
	}
}

func TestValidateAccessToken(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	pair, _ := GenerateTokenPair("user-1", uint16(bitmask.TierManager), secret)

	claims, err := ValidateAccessToken(pair.AccessToken, secret)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Errorf("subject = %q, want user-1", claims.Subject)
	}
	if claims.Permissions != uint16(bitmask.TierManager) {
		t.Errorf("permissions = %v, want TierManager", claims.Permissions)
	}
}

func TestValidateAccessTokenWrongSecret(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	pair, _ := GenerateTokenPair("user-1", 0, secret)

	_, err := ValidateAccessToken(pair.AccessToken, "wrong-secret")
	if err == nil {
		t.Error("wrong secret should fail validation")
	}
}

func TestValidateMalformedToken(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	_, err := ValidateAccessToken("not-a-real-token", secret)
	if err == nil {
		t.Error("malformed token should fail")
	}
}

func TestValidateRefreshToken(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	pair, _ := GenerateTokenPair("user-1", 0, secret)

	claims, err := ValidateRefreshToken(pair.RefreshToken, secret)
	if err != nil {
		t.Fatalf("ValidateRefreshToken: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Errorf("subject = %q, want user-1", claims.Subject)
	}
}

func TestGenerateAccessToken(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	token, expiresAt, err := GenerateAccessToken("user-1", uint16(bitmask.TierSilver), secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
	if expiresAt <= time.Now().Unix() {
		t.Error("expires_at should be in the future")
	}

	claims, err := ValidateAccessToken(token, secret)
	if err != nil {
		t.Fatalf("validate generated token: %v", err)
	}
	if claims.Permissions != uint16(bitmask.TierSilver) {
		t.Errorf("permissions = %v, want TierSilver", claims.Permissions)
	}
}

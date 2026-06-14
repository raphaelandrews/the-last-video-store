package auth

import (
	"testing"
	"time"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("GenerateTOTPSecret: %v", err)
	}
	if secret == "" {
		t.Error("secret should not be empty")
	}

	secret2, _ := GenerateTOTPSecret()
	if secret == secret2 {
		t.Error("two secrets should differ")
	}
}

func TestGenerateTOTPCode(t *testing.T) {
	secret, _ := GenerateTOTPSecret()

	code, err := GenerateTOTPCode(secret, time.Now())
	if err != nil {
		t.Fatalf("GenerateTOTPCode: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("code length = %d, want 6", len(code))
	}

	code2, _ := GenerateTOTPCode(secret, time.Now())
	if code != code2 {
		t.Error("same secret + time should give same code")
	}
}

func TestValidateTOTPCode(t *testing.T) {
	secret, _ := GenerateTOTPSecret()
	code, _ := GenerateTOTPCode(secret, time.Now())

	if !ValidateTOTPCode(secret, code) {
		t.Error("valid code should verify")
	}
}

func TestValidateTOTPCodeWrong(t *testing.T) {
	secret, _ := GenerateTOTPSecret()

	if ValidateTOTPCode(secret, "000000") {
		t.Error("wrong code should not verify")
	}
}

func TestValidateTOTPCodeEmpty(t *testing.T) {
	if ValidateTOTPCode("", "123456") {
		t.Error("empty secret should not verify")
	}
	if ValidateTOTPCode("SECRET", "") {
		t.Error("empty code should not verify")
	}
}

func TestValidateTOTPCodeWindow(t *testing.T) {
	secret, _ := GenerateTOTPSecret()

	past := time.Now().Add(-30 * time.Second)
	pastCode, _ := GenerateTOTPCode(secret, past)
	if !ValidateTOTPCode(secret, pastCode) {
		t.Error("code from -30s should verify (window tolerance)")
	}

	future := time.Now().Add(30 * time.Second)
	futureCode, _ := GenerateTOTPCode(secret, future)
	if !ValidateTOTPCode(secret, futureCode) {
		t.Error("code from +30s should verify (window tolerance)")
	}

	stale := time.Now().Add(-60 * time.Second)
	staleCode, _ := GenerateTOTPCode(secret, stale)
	if ValidateTOTPCode(secret, staleCode) {
		t.Error("code from -60s should not verify (outside window)")
	}
}

func TestGenerateTOTPURL(t *testing.T) {
	secret, _ := GenerateTOTPSecret()
	url := GenerateTOTPURL("TLVS", "testuser", secret)

	if url == "" {
		t.Error("TOTP URL should not be empty")
	}
	if url[:15] != "otpauth://totp/" {
		t.Errorf("URL should start with otpauth://totp/, got: %s", url)
	}
}

func TestGenerateTOTPTempToken(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	token, expiresAt, err := GenerateTOTPTempToken("user-1", secret)
	if err != nil {
		t.Fatalf("GenerateTOTPTempToken: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty")
	}
	if expiresAt <= time.Now().Unix() {
		t.Error("expires_at should be in the future")
	}

	userID, err := ValidateTOTPTempToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateTOTPTempToken: %v", err)
	}
	if userID != "user-1" {
		t.Errorf("userID = %q, want user-1", userID)
	}
}

func TestValidateTOTPTempTokenWrongSecret(t *testing.T) {
	secret := "test-jwt-secret-32-bytes-long!!"
	token, _, _ := GenerateTOTPTempToken("user-1", secret)

	_, err := ValidateTOTPTempToken(token, "wrong-secret")
	if err == nil {
		t.Error("wrong secret should fail")
	}
}

package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrTOTPCodeInvalid = errors.New("invalid TOTP code")

type TOTPVerifyResult struct {
	Valid   bool
	Message string
}

func GenerateTOTPSecret() (string, error) {
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate totp secret: %w", err)
	}
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)
	return secret, nil
}

func GenerateTOTPCode(secret string, t time.Time) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("decode totp secret: %w", err)
	}

	counter := t.Unix() / 30
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	code := int(binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF)
	code %= 1000000

	return fmt.Sprintf("%06d", code), nil
}

func ValidateTOTPCode(secret, code string) bool {
	if secret == "" || code == "" {
		return false
	}

	now := time.Now()

	for _, offset := range []int64{0, -30, 30} {
		t := now.Add(time.Duration(offset) * time.Second)
		expectedCode, err := GenerateTOTPCode(secret, t)
		if err != nil {
			continue
		}
		if hmac.Equal([]byte(expectedCode), []byte(code)) {
			return true
		}
	}

	return false
}

func GenerateTOTPURL(issuer, accountName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		issuer, accountName, secret, issuer)
}

func GenerateTOTPTempToken(userID string, secret string) (string, int64, error) {
	now := time.Now()
	expires := now.Add(5 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expires),
		ID:        uuid.NewString(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return token, expires.Unix(), nil
}

func ValidateTOTPTempToken(tokenString, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", ErrTokenInvalid
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", ErrTokenExpired
	}

	return claims.Subject, nil
}

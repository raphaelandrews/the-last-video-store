package crypto

import (
	"bytes"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("GenerateAESKey: %v", err)
	}

	plaintext := []byte("sensitive audit log data")
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}

	ciphertext2, _ := Encrypt(plaintext, key)
	if bytes.Equal(ciphertext, ciphertext2) {
		t.Error("two encryptions should differ (random nonce)")
	}
}

func TestAESWrongKey(t *testing.T) {
	key1, _ := GenerateAESKey()
	key2, _ := GenerateAESKey()

	ciphertext, _ := Encrypt([]byte("test"), key1)
	_, err := Decrypt(ciphertext, key2)
	if err != ErrDecryptionFailed {
		t.Errorf("wrong key should fail decryption, got: %v", err)
	}
}

func TestAESTamperedCiphertext(t *testing.T) {
	key, _ := GenerateAESKey()
	ciphertext, _ := Encrypt([]byte("test"), key)

	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err := Decrypt(ciphertext, key)
	if err != ErrDecryptionFailed {
		t.Errorf("tampered ciphertext should fail, got: %v", err)
	}
}

func TestAESEmptyPlaintext(t *testing.T) {
	key, _ := GenerateAESKey()
	ciphertext, err := Encrypt([]byte{}, key)
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("decrypted empty = %q, want empty", decrypted)
	}
}

func TestAESKeyLength(t *testing.T) {
	_, err := Encrypt([]byte("test"), []byte("short-key"))
	if err != ErrKeyLength {
		t.Errorf("short key should return ErrKeyLength, got: %v", err)
	}

	_, err = Decrypt([]byte("test"), []byte("short-key"))
	if err != ErrKeyLength {
		t.Errorf("short key should return ErrKeyLength, got: %v", err)
	}
}

func TestGenerateAESKey(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("GenerateAESKey: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("key length = %d, want 32", len(key))
	}

	key2, _ := GenerateAESKey()
	if bytes.Equal(key, key2) {
		t.Error("two generated keys should differ")
	}
}

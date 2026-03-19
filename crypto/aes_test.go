package crypto

import (
	"bytes"
	"testing"
)

func TestAESGCM_EncryptDecrypt(t *testing.T) {
	key, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	aes, err := NewAESGCM(key)
	if err != nil {
		t.Fatalf("failed to create AESGCM: %v", err)
	}

	plaintext := []byte("hello, world!")
	ciphertext, err := aes.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Fatalf("ciphertext should be different from plaintext")
	}

	decrypted, err := aes.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("decrypted text does not match plaintext: expected %s, got %s", plaintext, decrypted)
	}
}

func TestAESGCM_InvalidKey(t *testing.T) {
	invalidKey := []byte("short")
	_, err := NewAESGCM(invalidKey)
	if err == nil {
		t.Fatalf("expected error for invalid key size, got nil")
	}
}

func TestAESGCM_DecryptionFailure(t *testing.T) {
	key, _ := GenerateKey(32)
	aes, _ := NewAESGCM(key)

	ciphertext := []byte("invalid-ciphertext")
	_, err := aes.Decrypt(ciphertext)
	if err == nil {
		t.Fatalf("expected error for invalid ciphertext, got nil")
	}
}

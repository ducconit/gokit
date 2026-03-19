package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// AESGCM represents an AES-GCM cipher for symmetric encryption and decryption.
type AESGCM struct {
	key []byte
}

// NewAESGCM creates a new AESGCM instance with the provided key.
// The key must be 16, 24, or 32 bytes long for AES-128, AES-192, or AES-256 respectively.
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("crypto/aes: invalid key size: %d, must be 16, 24, or 32", len(key))
	}
	return &AESGCM{key: key}, nil
}

// Encrypt encrypts the plaintext using AES-GCM.
// It returns a byte slice containing the nonce followed by the ciphertext.
func (a *AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to generate nonce: %w", err)
	}

	// Seal appends the authentication tag to the ciphertext
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts the ciphertext using AES-GCM.
// The ciphertext must be in the format: [nonce][encrypted_data][tag].
func (a *AESGCM) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("crypto/aes: ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: decryption failed: %w", err)
	}

	return plaintext, nil
}

// GenerateKey generates a random key of the specified size (16, 24, or 32 bytes).
func GenerateKey(size int) ([]byte, error) {
	if size != 16 && size != 24 && size != 32 {
		return nil, fmt.Errorf("crypto/aes: invalid key size: %d", size)
	}
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("crypto/aes: failed to generate key: %w", err)
	}
	return key, nil
}

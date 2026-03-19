package crypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	basecrypto "github.com/ducconit/gokit/crypto"
)

// APIEncryptor handles encryption and decryption of API payloads.
type APIEncryptor struct {
	aes *basecrypto.AESGCM
}

// NewAPIEncryptor creates a new APIEncryptor with the provided AES key.
func NewAPIEncryptor(key []byte) (*APIEncryptor, error) {
	aes, err := basecrypto.NewAESGCM(key)
	if err != nil {
		return nil, err
	}
	return &APIEncryptor{aes: aes}, nil
}

// EncryptPayload marshals the payload to JSON, encrypts it using AES-GCM,
// and returns the result as a base64 encoded string.
func (e *APIEncryptor) EncryptPayload(payload any) (string, error) {
	plaintext, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("crypto/api: failed to marshal payload: %w", err)
	}

	ciphertext, err := e.aes.Encrypt(plaintext)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPayload decodes the base64 string, decrypts it using AES-GCM,
// and unmarshals the resulting JSON into the destination object.
func (e *APIEncryptor) DecryptPayload(data string, dst any) error {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("crypto/api: failed to decode base64: %w", err)
	}

	plaintext, err := e.aes.Decrypt(ciphertext)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(plaintext, dst); err != nil {
		return fmt.Errorf("crypto/api: failed to unmarshal payload: %w", err)
	}

	return nil
}

// HybridPayload represents the encrypted data sent between client and server.
type HybridPayload struct {
	EncryptedKey string `json:"encrypted_key"` // Base64(RSA_Encrypt(AES_Key))
	Payload      string `json:"payload"`       // Base64(AES_Encrypt(Data))
}

// HybridEncryptor handles hybrid encryption using RSA for key exchange and AES for payload encryption.
type HybridEncryptor struct {
	rsa *basecrypto.RSA
}

// NewHybridEncryptor creates a new HybridEncryptor with the provided RSA instance.
func NewHybridEncryptor(rsa *basecrypto.RSA) *HybridEncryptor {
	return &HybridEncryptor{rsa: rsa}
}

// Encrypt encrypts the data using a randomly generated AES key, which is itself encrypted with RSA.
// This is typically used by the client to send data to the server.
func (h *HybridEncryptor) Encrypt(data any) (*HybridPayload, error) {
	// 1. Generate random AES key (32 bytes for AES-256)
	aesKey, err := basecrypto.GenerateKey(32)
	if err != nil {
		return nil, fmt.Errorf("crypto/hybrid: failed to generate session key: %w", err)
	}

	// 2. Encrypt AES key with RSA Public Key
	encryptedKey, err := h.rsa.Encrypt(aesKey)
	if err != nil {
		return nil, fmt.Errorf("crypto/hybrid: failed to encrypt session key: %w", err)
	}

	// 3. Create AES encryptor for payload
	aesEncryptor, err := NewAPIEncryptor(aesKey)
	if err != nil {
		return nil, err
	}

	// 4. Encrypt payload
	encryptedPayload, err := aesEncryptor.EncryptPayload(data)
	if err != nil {
		return nil, err
	}

	return &HybridPayload{
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKey),
		Payload:      encryptedPayload,
	}, nil
}

// Decrypt decrypts the HybridPayload using the RSA Private Key to recover the AES key.
// This is typically used by the server to receive data from the client.
func (h *HybridEncryptor) Decrypt(hp *HybridPayload, dst any) error {
	// 1. Decode and decrypt the AES session key
	encryptedKey, err := base64.StdEncoding.DecodeString(hp.EncryptedKey)
	if err != nil {
		return fmt.Errorf("crypto/hybrid: failed to decode encrypted key: %w", err)
	}

	aesKey, err := h.rsa.Decrypt(encryptedKey)
	if err != nil {
		return fmt.Errorf("crypto/hybrid: failed to decrypt session key: %w", err)
	}

	// 2. Create AES encryptor with recovered key
	aesEncryptor, err := NewAPIEncryptor(aesKey)
	if err != nil {
		return err
	}

	// 3. Decrypt payload
	return aesEncryptor.DecryptPayload(hp.Payload, dst)
}

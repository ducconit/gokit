package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// RSA represents an RSA cipher for asymmetric encryption and decryption.
type RSA struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
}

// NewRSA creates an RSA instance from a private key.
func NewRSA(priv *rsa.PrivateKey) *RSA {
	return &RSA{
		priv: priv,
		pub:  &priv.PublicKey,
	}
}

// NewRSAWithPublicKey creates an RSA instance from a public key (for encryption only).
func NewRSAWithPublicKey(pub *rsa.PublicKey) *RSA {
	return &RSA{
		pub: pub,
	}
}

// Encrypt encrypts the plaintext using RSA-OAEP with SHA-256.
func (r *RSA) Encrypt(plaintext []byte) ([]byte, error) {
	if r.pub == nil {
		return nil, fmt.Errorf("crypto/rsa: missing public key for encryption")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.pub, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto/rsa: encryption failed: %w", err)
	}

	return ciphertext, nil
}

// Decrypt decrypts the ciphertext using RSA-OAEP with SHA-256.
func (r *RSA) Decrypt(ciphertext []byte) ([]byte, error) {
	if r.priv == nil {
		return nil, fmt.Errorf("crypto/rsa: missing private key for decryption")
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.priv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto/rsa: decryption failed: %w", err)
	}

	return plaintext, nil
}

// GenerateRSAKey generates a new RSA private key of the specified bits.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, fmt.Errorf("crypto/rsa: key size too small, must be at least 2048 bits")
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// EncodePublicKeyToPEM encodes an RSA public key to PEM format.
func EncodePublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubPEM, nil
}

// EncodePrivateKeyToPEM encodes an RSA private key to PEM format.
func EncodePrivateKeyToPEM(priv *rsa.PrivateKey) []byte {
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	return privPEM
}

// DecodePrivateKeyFromPEM decodes an RSA private key from PEM format.
func DecodePrivateKeyFromPEM(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, fmt.Errorf("crypto/rsa: failed to decode PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto/rsa: failed to parse private key: %w", err)
	}

	return priv, nil
}

// DecodePublicKeyFromPEM decodes an RSA public key from PEM format.
func DecodePublicKeyFromPEM(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, fmt.Errorf("crypto/rsa: failed to decode PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto/rsa: failed to parse public key: %w", err)
	}

	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("crypto/rsa: not an RSA public key")
	}

	return pub, nil
}

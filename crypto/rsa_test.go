package crypto

import (
	"bytes"
	"testing"
)

func TestRSA_EncryptDecrypt(t *testing.T) {
	priv, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	rsaCipher := NewRSA(priv)

	plaintext := []byte("hello, world!")
	ciphertext, err := rsaCipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Fatalf("ciphertext should be different from plaintext")
	}

	decrypted, err := rsaCipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("decrypted text does not match plaintext: expected %s, got %s", plaintext, decrypted)
	}
}

func TestRSA_PEMEncoding(t *testing.T) {
	priv, _ := GenerateRSAKey(2048)
	privPEM := EncodePrivateKeyToPEM(priv)
	pubPEM, _ := EncodePublicKeyToPEM(&priv.PublicKey)

	decodedPriv, err := DecodePrivateKeyFromPEM(privPEM)
	if err != nil {
		t.Fatalf("failed to decode private key: %v", err)
	}

	decodedPub, err := DecodePublicKeyFromPEM(pubPEM)
	if err != nil {
		t.Fatalf("failed to decode public key: %v", err)
	}

	if priv.N.Cmp(decodedPriv.N) != 0 {
		t.Fatalf("decoded private key N mismatch")
	}

	if priv.PublicKey.N.Cmp(decodedPub.N) != 0 {
		t.Fatalf("decoded public key N mismatch")
	}
}

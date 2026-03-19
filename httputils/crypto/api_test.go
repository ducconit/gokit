package crypto

import (
	"testing"

	basecrypto "github.com/ducconit/gokit/crypto"
)

func TestAPIEncryptor_EncryptDecryptPayload(t *testing.T) {
	key, _ := basecrypto.GenerateKey(32)
	encryptor, err := NewAPIEncryptor(key)
	if err != nil {
		t.Fatalf("failed to create APIEncryptor: %v", err)
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	user := User{ID: 1, Name: "Alice"}
	encrypted, err := encryptor.EncryptPayload(user)
	if err != nil {
		t.Fatalf("failed to encrypt payload: %v", err)
	}

	if encrypted == "" {
		t.Fatalf("encrypted string should not be empty")
	}

	var decryptedUser User
	err = encryptor.DecryptPayload(encrypted, &decryptedUser)
	if err != nil {
		t.Fatalf("failed to decrypt payload: %v", err)
	}

	if decryptedUser.ID != user.ID || decryptedUser.Name != user.Name {
		t.Fatalf("decrypted user does not match original: expected %+v, got %+v", user, decryptedUser)
	}
}

func TestHybridEncryptor_EncryptDecrypt(t *testing.T) {
	// 1. Setup Server RSA Key
	priv, _ := basecrypto.GenerateRSAKey(2048)
	rsa := basecrypto.NewRSA(priv)
	hybrid := NewHybridEncryptor(rsa)

	type SecretData struct {
		Message string `json:"message"`
	}

	data := SecretData{Message: "this is a hybrid secret"}

	// 2. Client Encrypts (using Server's RSA Public Key context)
	hp, err := hybrid.Encrypt(data)
	if err != nil {
		t.Fatalf("failed to hybrid encrypt: %v", err)
	}

	if hp.EncryptedKey == "" || hp.Payload == "" {
		t.Fatalf("hybrid payload components should not be empty")
	}

	// 3. Server Decrypts (using RSA Private Key)
	var decryptedData SecretData
	err = hybrid.Decrypt(hp, &decryptedData)
	if err != nil {
		t.Fatalf("failed to hybrid decrypt: %v", err)
	}

	if decryptedData.Message != data.Message {
		t.Fatalf("decrypted message mismatch: expected %s, got %s", data.Message, decryptedData.Message)
	}
}

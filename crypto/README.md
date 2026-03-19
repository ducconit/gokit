# crypto

Package cung cấp các giải pháp mã hóa và giải mã dữ liệu an toàn dựa trên thư viện chuẩn của Go.

## Tính năng

- **AES-GCM**: Mã hóa đối xứng với cơ chế xác thực (Authenticated Encryption), an toàn và hiệu suất cao.
- **RSA-OAEP**: Mã hóa bất đối xứng với chuẩn OAEP (SHA-256), dùng cho trao đổi khóa hoặc dữ liệu nhỏ nhạy cảm.
- **PEM Support**: Hỗ trợ encode/decode RSA keys sang định dạng PEM để lưu trữ.

## Cách dùng

### AES-GCM (Mã hóa đối xứng)

```go
// Tạo key 32 bytes (AES-256)
key, _ := crypto.GenerateKey(32)
cipher, _ := crypto.NewAESGCM(key)

// Mã hóa
plaintext := []byte("thông tin bí mật")
ciphertext, _ := cipher.Encrypt(plaintext)

// Giải mã
decrypted, _ := cipher.Decrypt(ciphertext)
```

### RSA (Mã hóa bất đối xứng)

```go
// Tạo cặp khóa mới
priv, _ := crypto.GenerateRSAKey(2048)
rsa := crypto.NewRSA(priv)

// Mã hóa bằng Public Key
ciphertext, _ := rsa.Encrypt([]byte("dữ liệu nhạy cảm"))

// Giải mã bằng Private Key
plaintext, _ := rsa.Decrypt(ciphertext)
```

### Quản lý PEM Key

```go
// Encode sang PEM
pubPEM, _ := crypto.EncodePublicKeyToPEM(&priv.PublicKey)
privPEM := crypto.EncodePrivateKeyToPEM(priv)

// Decode từ PEM
decodedPriv, _ := crypto.DecodePrivateKeyFromPEM(privPEM)
decodedPub, _ := crypto.DecodePublicKeyFromPEM(pubPEM)
```

## Security Notes

- Luôn sử dụng key đủ độ dài (32 bytes cho AES-256).
- RSA nên sử dụng kích thước ít nhất 2048 bits.
- Package sử dụng `crypto/rand` để đảm bảo tính ngẫu nhiên an toàn về mặt mật mã.

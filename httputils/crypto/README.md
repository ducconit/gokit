# httputils/crypto

Tiện ích mã hóa dành riêng cho các giao tiếp HTTP API, giúp bảo vệ payload (request body, response body) một cách tự động.

## Tính năng

- **API Payload Encryption**: Tự động Serialize JSON -> Encrypt -> Base64 Encode.
- **Hybrid Encryption (RSA + AES)**: Giải pháp an toàn nhất cho Client-Server. Dùng RSA để trao đổi AES Session Key và AES-GCM để mã hóa payload.
- **Base64 Encoding**: Kết quả mã hóa luôn được encode sang Base64 để truyền qua giao thức HTTP an toàn.
- **Tích hợp sẵn AES-GCM**: Sử dụng package `crypto` nền tảng để đảm bảo tính toàn vẹn và bảo mật của dữ liệu.

## Cách dùng

### 1. Hybrid Encryption (Khuyên dùng cho Client-Server)

Đây là giải pháp giải quyết vấn đề lộ Secret Key ở Client. Client sẽ tự tạo AES key tạm thời và gửi cho Server thông qua RSA Public Key.

```go
import (
	basecrypto "github.com/ducconit/gokit/crypto"
	apicrypto "github.com/ducconit/gokit/httputils/crypto"
)

// --- PHÍA SERVER ---
// Khởi tạo với RSA Private Key (để giải mã AES key từ client)
privKey, _ := basecrypto.DecodePrivateKeyFromPEM(pemBytes)
rsa := basecrypto.NewRSA(privKey)
hybrid := apicrypto.NewHybridEncryptor(rsa)

// Giải mã dữ liệu từ Client
var input YourStruct
err := hybrid.Decrypt(payloadFromClient, &input)

// --- PHÍA CLIENT ---
// Khởi tạo với RSA Public Key của Server (để mã hóa AES key)
pubKey, _ := basecrypto.DecodePublicKeyFromPEM(serverPubPEM)
rsa := basecrypto.NewRSAWithPublicKey(pubKey)
hybrid := apicrypto.NewHybridEncryptor(rsa)

// Mã hóa dữ liệu gửi lên Server
payload, err := hybrid.Encrypt(yourData)
// payload sẽ có dạng: { "encrypted_key": "...", "payload": "..." }
```

### 2. Simple API Payload Encryption (Dùng khi cả 2 bên đều giữ Secret an toàn)

```go
import (
	basecrypto "github.com/ducconit/gokit/crypto"
	apicrypto "github.com/ducconit/gokit/httputils/crypto"
)

// Khởi tạo (sử dụng 32 bytes AES key)
key, _ := basecrypto.GenerateKey(32)
encryptor, _ := apicrypto.NewAPIEncryptor(key)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// 1. Mã hóa Payload (User -> JSON -> Encrypt -> Base64)
user := User{ID: 1, Name: "Alice"}
encryptedStr, err := encryptor.EncryptPayload(user)

// 2. Giải mã Payload (Base64 -> Decrypt -> JSON -> User)
var decodedUser User
err = encryptor.DecryptPayload(encryptedStr, &decodedUser)
```

## Luồng hoạt động

1. **Mã hóa**: 
   - Marshal object sang JSON bytes.
   - Sử dụng AES-GCM (từ package `crypto`) để mã hóa.
   - Encode kết quả cuối cùng sang chuỗi Base64.
2. **Giải mã**:
   - Decode chuỗi Base64.
   - Giải mã bằng AES-GCM (kiểm tra tính toàn vẹn).
   - Unmarshal JSON về struct đích.

## Ưu điểm

- Đơn giản hóa việc bảo mật API đầu cuối (End-to-End Encryption).
- Tránh rò rỉ dữ liệu nhạy cảm ngay cả khi HTTPS bị phá vỡ hoặc có proxy can thiệp.
- Đảm bảo payload không bị thay đổi trên đường truyền.

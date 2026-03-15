# otp

Quản lý OTP theo purpose (quên mật khẩu, verify account, verify transaction, ...) và gửi qua nhiều channel (sms/email/telegram/discord/...).

Package này tập trung vào:

- generate + store OTP (hash)
- verify OTP với TTL + giới hạn số lần thử
- routing sender theo channel

## Cách dùng

### Tạo manager (memory store + router sender)

```go
store := otp.NewMemoryStore()

sender := otp.Router{
	Senders: map[otp.Channel]otp.Sender{
		"sms":   smsSender,
		"email": emailSender,
	},
}

m, err := otp.New(otp.Config{
	CodeLength:  6,
	TTL:         5 * time.Minute,
	MaxAttempts: 5,
}, store, sender)
if err != nil {
	panic(err)
}
```

### Request OTP

```go
// Gửi kèm metadata (ví dụ tên user, địa chỉ IP, ...)
metadata := map[string]string{
	"user_name": "Nguyen Van A",
	"ip":        "1.1.1.1",
}

msg, err := m.Request(ctx, "forgot_password", "sms", "+8490xxxxxxx", metadata)
if err != nil {
	panic(err)
}
_ = msg.ExpiresAt
```

### Verify OTP

```go
if err := m.Verify(ctx, "forgot_password", "+8490xxxxxxx", "123456"); err != nil {
	// otp.ErrInvalid / otp.ErrExpired / otp.ErrMaxAttempts ...
	panic(err)
}
```

## Stores

- `otp.NewMemoryStore()`: Lưu trữ trong bộ nhớ (phù hợp cho testing/single instance).
- `otp.NewCachingStore(cacheManager)`: Lưu trữ thông qua package `cache` nội bộ (hỗ trợ Redis, Memcache, BigCache, ...).
- `otp.NewSQLStore(db, dialect, table)`: Lưu trữ trong cơ sở dữ liệu (Postgres, MySQL, SQLite).

## Senders

Package hỗ trợ interface `Sender` để bạn tự triển khai logic gửi qua các provider khác nhau.

```go
type Sender interface {
	Send(ctx context.Context, msg Message) error
}
```

### Ví dụ triển khai EmailSender với Metadata

```go
type EmailSender struct {}

func (s *EmailSender) Send(ctx context.Context, msg otp.Message) error {
    userName := msg.Metadata["user_name"]
    if userName == "" {
        userName = "Khách hàng"
    }

    var body string
    switch msg.Purpose {
    case "forgot_password":
        body = fmt.Sprintf("Chào %s, mã khôi phục mật khẩu của bạn là: %s", userName, msg.Code)
    default:
        body = fmt.Sprintf("Mã xác thực của bạn là: %s", msg.Code)
    }

    fmt.Printf("Gửi email đến %s: %s\n", msg.Recipient, body)
    return nil
}
```

### Sử dụng Router

Sử dụng `otp.Router` để tự động định tuyến `Sender` dựa trên `Channel` của yêu cầu:

```go
sender := otp.Router{
	Senders: map[otp.Channel]otp.Sender{
		"sms":   &SMSSender{},
		"email": &EmailSender{},
	},
}
```


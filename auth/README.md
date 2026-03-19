# auth

JWT auth (access/refresh token) + session management theo nhiều driver (memory/redis/memcache/sql).

## Concepts

- Access token: TTL ngắn, dùng để authorize request
- Refresh token: TTL dài hơn, dùng để rotate/renew access token
- Session: record lưu phía server để revoke/banned/rotate refresh token an toàn hơn

## Cách dùng

### Khởi tạo (memory store)

```go
store := auth.NewMemoryStore()
m, err := auth.New(auth.Config{
	Issuer:     "my-service",
	Audience:   []string{"my-service"},
	AccessTTL:  15 * time.Minute,
	RefreshTTL: 30 * 24 * time.Hour,
	HMACSecret: []byte("change-me"),
}, store)
if err != nil {
	panic(err)
}
```

### Khởi tạo (SQL store với logic mặc định)

Nếu bạn truyền `DB` vào `SQLStore`, hệ thống sẽ tự động sử dụng logic thực thi mặc định.

```go
store := &auth.SQLStore{
    DB:      db,
    Dialect: auth.SQLDialectQuestion,
}
m, err := auth.New(auth.Config{
    // ... các config khác
}, store)
```

### Khởi tạo (SQL store với logic tùy chỉnh)

Bạn có thể ghi đè logic thực thi SQL thông qua `Config` (ví dụ: để logging hoặc custom transaction).

### Login → issue token pair

Hệ thống hỗ trợ nhiều loại đối tượng (vd: `user`, `customer`, `admin`) để phân quyền linh hoạt.

```go
// Metadata lưu thông tin thiết bị, user data...
metadata := map[string]string{
    "device": "iPhone 15",
    "fcm_token": "token-123",
}

// Issue token cho 'admin'
pair, err := m.Issue(ctx, "admin-001", "admin", metadata)

// Issue token cho 'customer'
pair, err := m.Issue(ctx, "cust-999", "customer", metadata)
```

### Verify access token → lấy thông tin định danh

Hệ thống sẽ tự động kiểm tra:
1. Tính hợp lệ của JWT (signature, exp, ...)
2. Sự tồn tại của Session trong Store (Session-based auth)
3. Trạng thái Banned của Subject

```go
claims, err := m.VerifyAccess(ctx, pair.AccessToken)
if err != nil {
	// auth.ErrUnauthorized / auth.ErrForbidden
	panic(err)
}

fmt.Println(claims.SubjectID)   // "admin-001"
fmt.Println(claims.SubjectType) // "admin"
```

### Refresh token → sinh access token mới

Bạn có thể chọn có xoay (rotate) Refresh Token hay không thông qua `RefreshOptions`. Mặc định nên để `Rotate: false` trừ khi bạn muốn thu hồi token cũ ngay lập tức.

```go
// Chỉ sinh Access Token mới, giữ nguyên Refresh Token cũ
pair, err := m.Refresh(ctx, refreshToken, auth.RefreshOptions{Rotate: false})

// Sinh cả Access Token và Refresh Token mới (Token Rotation)
pair, err := m.Refresh(ctx, refreshToken, auth.RefreshOptions{Rotate: true})
```

### Logout / logout all

```go
// Đăng xuất thiết bị hiện tại
_ = m.Logout(ctx, pair.SessionID)

// Đăng xuất khỏi tất cả các thiết bị của subject
_ = m.LogoutAll(ctx, "user-123")
```

### Ban subject

```go
_ = m.BanSubject(ctx, "user-123", time.Now().Add(24*time.Hour))
```

## Social Login (Goth) Integration

Package `auth` có thể tích hợp dễ dàng với [goth](https://github.com/markbates/goth) để hỗ trợ đăng nhập qua mạng xã hội (Google, GitHub, Facebook...).

**Quy trình:**
1. Sử dụng `goth` để thực hiện luồng OAuth và lấy thông tin User từ Provider.
2. Sử dụng `auth.Manager.Issue` để cấp phát Session/JWT nội bộ cho ứng dụng.

```go
func CallbackHandler(ctx *gin.Context) {
    // 1. Goth thực hiện lấy thông tin User từ mạng xã hội
    gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
    if err != nil {
        return
    }

    // 2. Map thông tin sang Metadata (Avatar, Name...)
    metadata := map[string]string{
        "name": gothUser.Name,
        "avatar": gothUser.AvatarURL,
        "provider": gothUser.Provider,
    }

    // 3. Issue Token/Session nội bộ (Dùng userID từ social làm định danh subject)
    pair, err := authManager.Issue(ctx, gothUser.UserID, "customer", metadata)
    
    // Trả về token cho Client
    ctx.JSON(200, pair)
}
```

## Stores

- `auth.NewMemoryStore()`: Lưu trữ trong bộ nhớ (phù hợp cho testing/single instance).
- `auth.NewCachingStore(sessionCache, banCache)`: Lưu trữ thông qua package `cache` nội bộ (hỗ trợ Redis, Memcache, ...).
- `auth.SQLStore{Dialect: auth.SQLDialectQuestion}`: Lưu trữ trong database thông qua các hàm thực thi được cấu hình trong `Config`.


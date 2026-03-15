# GoKit (gobase)

Bộ thư viện (SDK) chuẩn hóa cho các dự án Go tại `github.com/ducconit/gokit`. Dự án cung cấp các package thiết yếu để xây dựng ứng dụng backend nhanh chóng, nhất quán và dễ bảo trì.

## Danh sách các Packages

| Package | Mô tả | Tính năng chính |
| :--- | :--- | :--- |
| [auth](./auth) | Quản lý định danh | JWT, Session-based auth, đa loại User (Admin/Customer), Ban user, Metadata. |
| [cache](./cache) | Hệ thống lưu trữ đệm | Hỗ trợ nhiều driver: Redis, Memcache, BigCache, Ristretto, v.v. |
| [config](./config) | Quản lý cấu hình | Load config từ Local file, S3, Redis, SQL, Cloudflare R2. |
| [id](./id) | Định danh duy nhất | Mặc định sử dụng **ULID** (thay thế cho UUID) giúp sắp xếp theo thời gian tốt hơn. |
| [logging](./logging) | Hệ thống ghi log | Dựa trên `zerolog`, hỗ trợ ghi log ra Console/File với định dạng JSON. |
| [otp](./otp) | Mã xác thực một lần | Generate/Verify OTP, hỗ trợ nhiều Store (Memory/Cache/SQL) và Metadata. |
| [pagination](./pagination) | Phân trang dữ liệu | Simple (Page/Limit) và Cursor-based pagination. Hỗ trợ tự động chuẩn hóa qua struct tags. |
| [retry](./retry) | Cơ chế thử lại | Tự động thử lại các tác vụ thất bại với các chiến lược backoff. |
| [validate](./validate) | Kiểm tra dữ liệu | Wrapper cho `go-playground/validator` để validate struct. |

## Yêu cầu hệ thống

- **Go version**: 1.25 trở lên.

## Cài đặt

```bash
go get github.com/ducconit/gokit
```

## Hướng dẫn nhanh

### 1. Quản lý ID (ULID)
Sử dụng ULID để có ID vừa duy nhất vừa có thể sắp xếp theo thời gian:
```go
import "github.com/ducconit/gokit/id"

uid := id.New() // "01AN4Z0MSP..."
```

### 2. Phân trang với Gin
Tận dụng tính năng chuẩn hóa tự động:
```go
type SearchReq struct {
    pagination.Simple // Có sẵn Page (mặc định 1), Limit (mặc định 20, max 100)
    Name string `form:"name"`
}

func Handler(c *gin.Context) {
    var req SearchReq
    c.ShouldBindQuery(&req)
    pagination.Normalize(&req) // Áp dụng default và max limit
    
    // Sử dụng req.Offset() và req.Limit cho SQL
}
```

### 3. Xác thực người dùng (Auth)
Hỗ trợ quản lý phiên đăng nhập (Session) và đa loại người dùng:
```go
// Đăng nhập cho Admin kèm thông tin thiết bị
metadata := map[string]string{"device": "iPhone 15"}
pair, _ := authManager.Issue(ctx, "admin-id", "admin", metadata)

// Kiểm tra token (tự động check Session Store và trạng thái Banned)
claims, _ := authManager.VerifyAccess(ctx, pair.AccessToken)
fmt.Println(claims.UserType) // "admin"
```

### 4. Gửi OTP với Metadata
```go
metadata := map[string]string{"user_name": "Anh Duc"}
msg, _ := otpManager.Request(ctx, "verify_email", "email", "duc@example.com", metadata)
// Metadata sẽ được truyền vào Sender để bạn tùy biến nội dung email
```

## Đóng góp

Vui lòng đọc hướng dẫn trong từng package cụ thể để biết thêm chi tiết về cách cấu hình và mở rộng.

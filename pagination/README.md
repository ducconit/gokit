# pagination

Hỗ trợ phân trang theo 2 dạng:

- Simple pagination (page/limit)
- Cursor pagination (after/limit) cho load more / infinity scroll

Hệ thống hỗ trợ tự động chuẩn hóa (normalization) dựa trên struct tags (`default`, `max`) và tích hợp tốt với các framework như Gin.

## Cách dùng

### 1. Tích hợp với Gin (Recommended)

Sử dụng `ctx.ShouldBind` của Gin và sau đó gọi `pagination.Normalize` để áp dụng giá trị mặc định và giới hạn.

```go
type MyRequest struct {
    pagination.Simple // Nhúng Page và Limit
    Keyword string `form:"name"`
}

func ListItems(ctx *gin.Context) {
    var req MyRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        return
    }

    // Tự động điền Page=1, Limit=20 (nếu trống) và cap Limit tối đa 100
    pagination.Normalize(&req)

    offset := req.Offset()
    limit := req.Limit
    // ... query database
}
```

### 2. Sử dụng với struct tùy chỉnh (Tags)

Bạn không nhất thiết phải dùng struct có sẵn, chỉ cần thêm tags:

```go
type CustomQuery struct {
    PageNum int `form:"p" default:"1"`
    Size    int `form:"s" default:"10" max:"50"`
}

q := CustomQuery{}
pagination.Normalize(&q)
// q.PageNum sẽ là 1, q.Size sẽ là 10
```

### 3. Bind từ url.Values (Standard Library)

```go
q := pagination.Simple{}
_ = pagination.BindQuery(r.URL.Query(), &q) // Tự động gọi Normalize() bên trong

offset := q.Offset()
```

### 4. Cursor pagination (Nâng cao: Đa cột)

Khi sắp xếp theo nhiều cột (vd: `Priority DESC`, `CreatedAt DESC`), con trỏ cần chứa giá trị của tất cả các cột này.

```go
type MyCursor struct {
    Priority  int       `json:"p"`
    CreatedAt time.Time `json:"c"`
    ID        string    `json:"i"`
}

// 1. Giải mã cursor từ request
p := pagination.Cursor{}
_ = pagination.BindQuery(r.URL.Query(), &p)

var cur MyCursor
if p.After != "" {
    _ = pagination.DecodeCursor(p.After, &cur)
}

// 2. Tạo cursor mới cho trang tiếp theo (lấy từ bản ghi cuối cùng của kết quả trả về)
next, _ := pagination.EncodeCursor(MyCursor{
    Priority:  lastItem.Priority,
    CreatedAt: lastItem.CreatedAt,
    ID:        lastItem.ID,
})
```

### Tại sao dùng Cursor?
- **Hiệu năng:** Luôn nhanh (`O(1)`) vì dùng index trực tiếp, không dùng `OFFSET` chậm chạp.
- **Tính nhất quán:** Không bị lặp hoặc sót dữ liệu khi có bản ghi mới chèn vào giữa lúc người dùng đang cuộn.

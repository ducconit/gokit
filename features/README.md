# Features

Gói `features` cung cấp một hệ thống quản lý Feature Flags linh hoạt, cho phép bạn bật/tắt các tính năng trong ứng dụng một cách động mà không cần triển khai lại mã nguồn.

## Các tính năng chính

- **Kiểm tra trạng thái Flag**: Bật/tắt tính năng toàn cục.
- **Quy tắc dựa trên ngữ cảnh (Rule-based)**: Bật tính năng cho một nhóm người dùng cụ thể dựa trên ID hoặc các thuộc tính tùy chỉnh (email, role, v.v.).
- **Toán tử linh hoạt**: Hỗ trợ các toán tử như `in`, `not_in`, `regex`.
- **Hỗ trợ lưu trữ đa dạng**: Mặc định cung cấp `MemoryStore`, dễ dàng mở rộng sang SQL, Redis, v.v. thông qua interface `Store`.

## Cách sử dụng

### Khởi tạo và sử dụng cơ bản

```go
import (
	"context"
	"fmt"
	"github.com/ducconit/gokit/features"
)

func main() {
	// 1. Tạo Store (ở đây dùng MemoryStore)
	store := features.NewMemoryStore()
	ctx := context.Background()

	// 2. Định nghĩa một tính năng
	newFeature := &features.Feature{
		Key:     "new-ui",
		Enabled: true,
		Rules: []features.Rule{
			{
				Attribute: "userID",
				Operator:  features.OpIn,
				Values:    []any{"user-1", "user-2"},
			},
		},
	}
	store.Set(ctx, newFeature)

	// 3. Khởi tạo Manager
	manager := features.NewManager(store)

	// 4. Kiểm tra trạng thái tính năng cho người dùng
	userCtx := &features.Context{
		UserID: "user-1",
	}

	if manager.IsEnabled(ctx, "new-ui", userCtx) {
		fmt.Println("Hiển thị giao diện mới cho user-1")
	}

	// Kiểm tra cho người dùng không nằm trong danh sách
	otherUser := &features.Context{
		UserID: "user-3",
	}
	if !manager.IsEnabled(ctx, "new-ui", otherUser) {
		fmt.Println("Hiển thị giao diện cũ cho user-3")
	}
}
```

### Sử dụng với thuộc tính tùy chỉnh (Custom Attributes)

Bạn có thể sử dụng bất kỳ thuộc tính nào để làm điều kiện bật/tắt tính năng:

```go
// Định nghĩa quy tắc dựa trên email
emailFeature := &features.Feature{
	Key:     "beta-tester",
	Enabled: true,
	Rules: []features.Rule{
		{
			Attribute: "email",
			Operator:  features.OpRegex,
			Values:    []any{`.*@company\.com$`},
		},
	},
}
store.Set(ctx, emailFeature)

// Kiểm tra dựa trên context có chứa email
userCtx := &features.Context{
	Attributes: map[string]any{
		"email": "dev@company.com",
	},
}

if manager.IsEnabled(ctx, "beta-tester", userCtx) {
	fmt.Println("Chào mừng thành viên công ty tham gia thử nghiệm beta!")
}
```

## Các thành phần quan trọng

- `Feature`: Đại diện cho một tính năng với mã định danh (`Key`), trạng thái (`Enabled`) và danh sách các quy tắc (`Rules`).
- `Rule`: Một điều kiện logic bao gồm thuộc tính (`Attribute`), toán tử (`Operator`) và các giá trị so sánh (`Values`).
- `Context`: Chứa thông tin về đối tượng đang được kiểm tra (ví dụ: người dùng hiện tại).
- `Store`: Interface để lưu trữ và truy xuất các Feature Flags.

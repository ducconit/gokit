# config

Wrapper cho `github.com/spf13/viper` để load config từ nhiều nguồn (local file, object storage, database, redis) và tùy chọn auto-reload.

## Cách dùng

### Local file + auto reload

```go
ctx := context.Background()

m, err := config.Load(ctx, &config.LocalFileSource{Path: "./config.yaml"}, config.Options{
	Format:     "yaml",
	AutoReload: true,
})
if err != nil {
	panic(err)
}
defer m.Close()

type AppConfig struct {
	Port int ` + "`mapstructure:\"port\"`" + `
}

var c AppConfig
if err := m.Unmarshal(&c); err != nil {
	panic(err)
}
```

### S3 / R2 / Minio

```go
ctx := context.Background()

src := &config.S3Source{
	Region: "ap-southeast-1",
	Bucket: "my-bucket",
	Key:    "configs/app.yaml",
}

m, err := config.Load(ctx, src, config.Options{
	Format:         "yaml",
	AutoReload:     true,
	ReloadInterval: 5 * time.Second,
})
if err != nil {
	panic(err)
}
defer m.Close()
```

### Database (SQLite/MySQL/PostgreSQL/...)

`SQLSource` dùng `database/sql`, nên ứng dụng sẽ tự import driver tương ứng (sqlite/mysql/postgres).

```go
src := &config.SQLSource{
	DB:    db,
	Query: "select value from app_configs where key = ?",
	Args:  []any{"app"},
}
```

## Auto reload

- Local file: dùng cơ chế watch của viper (fsnotify)
- Nguồn khác: polling theo `ReloadInterval` và reload khi nội dung thay đổi


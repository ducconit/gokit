# cache

Wrapper cho `github.com/eko/gocache/lib/v4` để khởi tạo cache theo `Config` và hỗ trợ nhiều store (driver) có sẵn của gocache.

## Cách dùng

### Redis

```go
ctx := context.Background()

m, err := cache.New[string](ctx, cache.Config{
	Driver:            cache.DriverRedis,
	DefaultExpiration: 30 * time.Second,
	Redis: cache.RedisConfig{
		Addr: "127.0.0.1:6379",
	},
})
if err != nil {
	panic(err)
}

_ = m.Set(ctx, "k1", "v1")
v, err := m.Get(ctx, "k1")
if err != nil {
	panic(err)
}
_ = v
```

### Memcache

```go
ctx := context.Background()

m, err := cache.New[[]byte](ctx, cache.Config{
	Driver:            cache.DriverMemcache,
	DefaultExpiration: 10 * time.Second,
	Memcache: cache.MemcacheConfig{
		Servers: []string{"127.0.0.1:11211"},
	},
})
if err != nil {
	panic(err)
}

_ = m.Set(ctx, "k1", []byte("v1"))
```

## Driver hỗ trợ

- bigcache
- freecache
- go-cache
- ristretto
- memcache
- redis
- rediscluster
- rueidis
- hazelcast
- pegasus


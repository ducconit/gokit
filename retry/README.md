# retry

Helper retry với exponential backoff + jitter và hỗ trợ `context.Context`.

## Cách dùng

```go
err := retry.Do(ctx, retry.Options{
	MaxAttempts: 5,
	BaseDelay:   200 * time.Millisecond,
	MaxDelay:    3 * time.Second,
	Jitter:      0.2,
}, func(ctx context.Context) error {
	return callRemote(ctx)
})
if err != nil {
	panic(err)
}
```


# logging

Wrapper cho `github.com/rs/zerolog` với mode output linh hoạt: console/file/both/disabled.

## Cách dùng

### Sử dụng zerolog trực tiếp (New)

```go
logger, closer, err := logging.New(logging.Config{
	Mode:          logging.ModeConsole,
	Level:         "debug",
	ConsolePretty: true,
})
if err != nil {
	panic(err)
}
defer closer.Close()

logger.Info().Str("service", "api").Msg("started")
```

### Sử dụng thông qua Interface (NewWrapper)

```go
logger, closer, err := logging.NewWrapper(logging.Config{
	Mode:          logging.ModeBoth,
	Level:         "info",
	FilePath:      "./app.log",
	ConsolePretty: false,
})
if err != nil {
	panic(err)
}
defer closer.Close()

// Sử dụng với các fields (key-value)
logger.Info("service started", "version", "1.0.0", "env", "production")

// Sử dụng với context (With)
childLogger := logger.With("request_id", "123")
childLogger.Debug("processing request")

// Sử dụng với lỗi
logger.Error("failed to process", err, "user_id", 456)
```

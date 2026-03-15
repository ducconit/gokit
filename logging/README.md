# logging

Wrapper cho `github.com/rs/zerolog` với mode output linh hoạt: console/file/both/disabled.

## Cách dùng

### Console pretty

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

### Ghi file + console

```go
logger, closer, err := logging.New(logging.Config{
	Mode:          logging.ModeBoth,
	Level:         "info",
	FilePath:      "./app.log",
	ConsolePretty: false,
})
if err != nil {
	panic(err)
}
defer closer.Close()

logger.Warn().Str("k", "v").Msg("something")
```


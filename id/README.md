# id

Helpers cho ID (ULID).

## Cách dùng

```go
id := id.New()

u, err := id.Parse(id)
if err != nil {
	panic(err)
}
_ = u
```


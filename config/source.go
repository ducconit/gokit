package config

import "context"

type Source interface {
	Load(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, b []byte) error
}

package config

import (
	"context"
	"os"
)

type LocalFileSource struct {
	Path string
}

func (s *LocalFileSource) Load(ctx context.Context) ([]byte, error) {
	_ = ctx
	return os.ReadFile(s.Path)
}

func (s *LocalFileSource) Write(ctx context.Context, b []byte) error {
	_ = ctx
	return os.WriteFile(s.Path, b, 0644)
}

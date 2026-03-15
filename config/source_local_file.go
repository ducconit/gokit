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

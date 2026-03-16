package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioSource struct {
	Endpoint string
	UseSSL   bool

	AccessKeyID     string
	SecretAccessKey string

	Bucket string
	Key    string

	once    sync.Once
	client  *minio.Client
	initErr error
}

func (s *MinioSource) Load(ctx context.Context) ([]byte, error) {
	if err := s.init(); err != nil {
		return nil, err
	}
	if s.Bucket == "" || s.Key == "" {
		return nil, fmt.Errorf("config: minio missing bucket/key")
	}

	obj, err := s.client.GetObject(ctx, s.Bucket, s.Key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}

func (s *MinioSource) Write(ctx context.Context, b []byte) error {
	if err := s.init(); err != nil {
		return err
	}
	if s.Bucket == "" || s.Key == "" {
		return fmt.Errorf("config: minio missing bucket/key")
	}

	_, err := s.client.PutObject(ctx, s.Bucket, s.Key, bytes.NewReader(b), int64(len(b)), minio.PutObjectOptions{})
	return err
}

func (s *MinioSource) init() error {
	s.once.Do(func() {
		if s.Endpoint == "" {
			s.initErr = fmt.Errorf("config: minio missing endpoint")
			return
		}

		client, err := minio.New(s.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s.AccessKeyID, s.SecretAccessKey, ""),
			Secure: s.UseSSL,
		})
		if err != nil {
			s.initErr = err
			return
		}
		s.client = client
	})
	return s.initErr
}

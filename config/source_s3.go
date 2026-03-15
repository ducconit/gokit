package config

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Source struct {
	Region string
	Bucket string
	Key    string

	once    sync.Once
	client  *s3.Client
	initErr error
}

func (s *S3Source) Load(ctx context.Context) ([]byte, error) {
	if err := s.init(ctx); err != nil {
		return nil, err
	}
	if s.Bucket == "" || s.Key == "" {
		return nil, fmt.Errorf("config: s3 missing bucket/key")
	}

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &s.Key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (s *S3Source) init(ctx context.Context) error {
	s.once.Do(func() {
		region := s.Region
		if region == "" {
			region = "us-east-1"
		}
		awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			s.initErr = err
			return
		}
		s.client = s3.NewFromConfig(awsCfg)
	})
	return s.initErr
}

package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Source struct {
	AccountID string
	Endpoint  string
	Region    string

	AccessKeyID     string
	SecretAccessKey string

	Bucket       string
	Key          string
	UsePathStyle *bool

	once    sync.Once
	client  *s3.Client
	initErr error
}

func (s *R2Source) Load(ctx context.Context) ([]byte, error) {
	if err := s.init(ctx); err != nil {
		return nil, err
	}
	if s.Bucket == "" || s.Key == "" {
		return nil, fmt.Errorf("config: r2 missing bucket/key")
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

func (s *R2Source) Write(ctx context.Context, b []byte) error {
	if err := s.init(ctx); err != nil {
		return err
	}
	if s.Bucket == "" || s.Key == "" {
		return fmt.Errorf("config: r2 missing bucket/key")
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &s.Key,
		Body:   bytes.NewReader(b),
	})
	return err
}

func (s *R2Source) init(ctx context.Context) error {
	s.once.Do(func() {
		if s.Endpoint == "" && s.AccountID != "" {
			s.Endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", s.AccountID)
		}
		if s.Endpoint == "" {
			s.initErr = fmt.Errorf("config: r2 missing endpoint/accountID")
			return
		}

		region := s.Region
		if region == "" {
			region = "auto"
		}

		creds := aws.CredentialsProvider(credentials.NewStaticCredentialsProvider(s.AccessKeyID, s.SecretAccessKey, ""))

		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:               s.Endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		awsCfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(creds),
			config.WithEndpointResolverWithOptions(resolver),
		)
		if err != nil {
			s.initErr = err
			return
		}

		usePathStyle := true
		if s.UsePathStyle != nil {
			usePathStyle = *s.UsePathStyle
		}
		s.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.UsePathStyle = usePathStyle
		})
	})
	return s.initErr
}

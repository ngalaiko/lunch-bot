package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Store struct {
	s3Client *s3.Client
}

func NewS3(cfg aws.Config) *S3Store {
	return &S3Store{
		s3Client: s3.NewFromConfig(cfg),
	}
}

// Store stores object by key in the bucket.
func (store *S3Store) Store(ctx context.Context, bucket, key string, body []byte) error {
	if _, err := store.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(body),
		StorageClass: types.StorageClassOnezoneIa,
	}); err != nil {
		return fmt.Errorf("failed to PutObject to '%s': %w", bucket, err)
	}
	return nil
}

type listKeysOptions struct {
	prefix *string
}

type ListKeysOption func(*listKeysOptions)

func getListKeysOptions(opts []ListKeysOption) *listKeysOptions {
	options := &listKeysOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func WithPrefix(prefix string) ListKeysOption {
	return func(options *listKeysOptions) {
		options.prefix = &prefix
	}
}

// ListKeys returns up to 1000 keys from the bucket.
func (store *S3Store) ListKeys(ctx context.Context, bucket string, opts ...ListKeysOption) ([]string, error) {
	options := getListKeysOptions(opts)

	response, err := store.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: options.prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ListObjectsV2 from '%s': %w", bucket, err)
	}

	resp := make([]string, 0, len(response.Contents))
	for _, o := range response.Contents {
		resp = append(resp, *o.Key)
	}

	return resp, nil
}

// Get returns object content by key from the bucket.
func (store *S3Store) Get(ctx context.Context, bucket, key string) ([]byte, error) {
	response, err := store.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to GetObject from '%s': %w", bucket, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

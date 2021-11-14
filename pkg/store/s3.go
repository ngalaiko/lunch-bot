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

type S3 struct {
	s3Client *s3.Client
}

func NewS3(cfg aws.Config) *S3 {
	return &S3{
		s3Client: s3.NewFromConfig(cfg),
	}
}

// Store stores object by key in the bucket.
func (store *S3) Store(ctx context.Context, bucket, key string, body []byte) error {
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

// ListKeys returns up to 1000 keys from the bucket.
func (store *S3) ListKeys(ctx context.Context, bucket string) ([]string, error) {
	response, err := store.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
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
func (store *S3) Get(ctx context.Context, bucket, key string) ([]byte, error) {
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

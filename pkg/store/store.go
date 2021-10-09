package store

import "context"

type Storage interface {
	Store(ctx context.Context, bucket, key string, data []byte) error
	ListKeys(ctx context.Context, bucket string, opts ...ListKeysOption) ([]string, error)
	Get(ctx context.Context, bucket, key string) ([]byte, error)
}

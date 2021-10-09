package store

import (
	"context"
	"fmt"
	"strings"
)

// verify interface
var _ Storage = &MemoryStore{}

type MemoryStore struct {
	store map[string]map[string][]byte
}

func NewInMemory() *MemoryStore {
	return &MemoryStore{
		store: map[string]map[string][]byte{},
	}
}

// Store stores object by key in the bucket.
func (s *MemoryStore) Store(ctx context.Context, bucket, key string, body []byte) error {
	b, bucketExists := s.store[bucket]
	if !bucketExists {
		b = map[string][]byte{}
		s.store[bucket] = b
	}

	b[key] = body

	return nil
}

// ListKeys returns up to 1000 keys from the bucket.
func (s *MemoryStore) ListKeys(ctx context.Context, bucket string, opts ...ListKeysOption) ([]string, error) {
	b, bucketExists := s.store[bucket]
	if !bucketExists {
		b = map[string][]byte{}
		s.store[bucket] = b
	}

	result := []string{}
	options := getListKeysOptions(opts)
	for key := range b {
		if options.prefix != nil && !strings.HasPrefix(key, *options.prefix) {
			continue
		}
		result = append(result, key)
	}

	return result, nil
}

// Get returns object content by key from the bucket.
func (s *MemoryStore) Get(ctx context.Context, bucket, key string) ([]byte, error) {
	b, bucketExists := s.store[bucket]
	if !bucketExists {
		b = map[string][]byte{}
		s.store[bucket] = b
	}

	value, found := b[key]
	if !found {
		return nil, fmt.Errorf("not found")
	}

	return value, nil
}

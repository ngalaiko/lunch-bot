package storage

import (
	"context"
	"fmt"

	"lunch/pkg/jwt/keys"
	"lunch/pkg/store"
)

type bolt struct {
	db     *store.Bolt
	bucket string
}

func NewBolt(db *store.Bolt) *bolt {
	return &bolt{
		db:     db,
		bucket: "jwt_keys",
	}
}

func (b *bolt) Create(ctx context.Context, key *keys.Key) error {
	if err := b.db.Put(ctx, b.bucket, key.ID, key); err != nil {
		return fmt.Errorf("failed to put: %w", err)
	}
	return nil
}

func (b *bolt) Get(ctx context.Context, id string) (*keys.Key, error) {
	key := &keys.Key{}
	if err := b.db.Get(ctx, b.bucket, id, key); err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return key, nil
}

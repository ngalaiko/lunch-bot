package boosts

import (
	"context"
	"fmt"

	"lunch/pkg/store"
)

type Store struct {
	storage    store.Storage
	bucketName string
}

func NewStore(storage store.Storage) *Store {
	return &Store{
		storage:    storage,
		bucketName: "lunch-boosts",
	}
}

func (rs *Store) Store(ctx context.Context, boost *Boost) error {
	if err := rs.storage.Store(ctx, rs.bucketName, boost.key(), []byte(boost.value())); err != nil {
		return fmt.Errorf("failed to store boost in storage: %w", err)
	}
	return nil
}

func (rs *Store) ListBoosts(ctx context.Context) ([]*Boost, error) {
	keys, err := rs.storage.ListKeys(ctx, rs.bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}

	boosts := make([]*Boost, 0, len(keys))
	for _, key := range keys {
		boost, err := boostFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key: %w", err)
		}
		boosts = append(boosts, boost)
	}

	return boosts, nil

}

package places

import (
	"context"
	"encoding/json"
	"fmt"

	"lunch/pkg/store"
)

type Store struct {
	storage    *store.S3Store
	bucketName string
}

func NewStore(storage *store.S3Store) *Store {
	return &Store{
		storage:    storage,
		bucketName: "lunch-places",
	}
}

func (ps *Store) Add(ctx context.Context, place *Place) error {
	jsonPlace, err := json.Marshal(place)
	if err != nil {
		return fmt.Errorf("failed to marshal place: %w", err)
	}

	if err := ps.storage.Store(ctx, ps.bucketName, place.Name, jsonPlace); err != nil {
		return fmt.Errorf("failed to store place in a storage: %w", err)
	}

	return nil
}

func (ps *Store) ListNames(ctx context.Context) ([]string, error) {
	names, err := ps.storage.ListKeys(ctx, ps.bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}
	return names, nil
}

func (ps *Store) Get(ctx context.Context, id string) (*Place, error) {
	rawPlace, err := ps.storage.Get(ctx, ps.bucketName, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get place from storage: %w", err)
	}

	place := &Place{}
	if err := json.Unmarshal(rawPlace, place); err != nil {
		return nil, fmt.Errorf("failed to unmarshal place: %w", err)
	}

	return place, nil
}

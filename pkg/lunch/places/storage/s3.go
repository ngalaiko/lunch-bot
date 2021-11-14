package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

type Store struct {
	storage    *store.S3
	bucketName string
}

func NewS3(storage *store.S3) *Store {
	return &Store{
		storage:    storage,
		bucketName: "lunch-places",
	}
}
func (ps *Store) Store(ctx context.Context, place *places.Place) error {
	jsonPlace, err := json.Marshal(place)
	if err != nil {
		return fmt.Errorf("failed to marshal place: %w", err)
	}
	if err := ps.storage.Store(ctx, ps.bucketName, string(place.Name), jsonPlace); err != nil {
		return fmt.Errorf("failed to store place in a storage: %w", err)
	}
	return nil
}
func (ps *Store) ListNames(ctx context.Context) ([]places.Name, error) {
	keys, err := ps.storage.ListKeys(ctx, ps.bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}
	names := make([]places.Name, len(keys))
	for i, key := range keys {
		names[i] = places.Name(key)
	}
	return names, nil
}
func (ps *Store) GetByName(ctx context.Context, name places.Name) (*places.Place, error) {
	rawPlace, err := ps.storage.Get(ctx, ps.bucketName, string(name))
	if err != nil {
		return nil, fmt.Errorf("failed to get place from storage: %w", err)
	}
	place := &places.Place{}
	if err := json.Unmarshal(rawPlace, place); err != nil {
		return nil, fmt.Errorf("failed to unmarshal place: %w", err)
	}
	return place, nil
}

package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

var _ Storage = &DynamoDBStorage{}

type DynamoDBStorage struct {
	storage *store.DynamoDB
}

func NewDynamoDB(storage *store.DynamoDB) *DynamoDBStorage {
	return &DynamoDBStorage{
		storage: storage,
	}
}

func (s *DynamoDBStorage) Store(ctx context.Context, place *places.Place) error {
	if err := s.storage.Execute(ctx, `
		INSERT INTO Places 
			value {
				'name': ?,
				'added_at': ?,
				'added_by': ?
			}
	`, place.Name, place.AddedAt.Unix(), place.AddedBy); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s *DynamoDBStorage) ListNames(ctx context.Context) ([]places.Name, error) {
	pp := []*places.Place{}
	if err := s.storage.Query(ctx, &pp, `SELECT * FROM Places`); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	names := make([]places.Name, 0, len(pp))
	for _, place := range pp {
		names = append(names, place.Name)
	}
	return names, nil
}

func (s *DynamoDBStorage) GetByName(ctx context.Context, name places.Name) (*places.Place, error) {
	places := []*places.Place{}
	if err := s.storage.Query(ctx, &places, `SELECT * FROM Places WHERE name = ?`, name); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return places[0], nil
}

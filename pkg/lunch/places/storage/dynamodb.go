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
				'id': ?,
				'name': ?,
				'added_at': ?,
				'added_by': ?
			}
	`, place.ID, place.Name, place.AddedAt.Unix(), place.AddedBy); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s *DynamoDBStorage) GetByID(ctx context.Context, id places.ID) (*places.Place, error) {
	places := []*places.Place{}
	if err := s.storage.Query(ctx, &places, `SELECT * FROM Places WHERE id = ?`, id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return places[0], nil
}

func (s *DynamoDBStorage) ListAll(ctx context.Context) ([]*places.Place, error) {
	pp := []*places.Place{}
	if err := s.storage.Query(ctx, &pp, `SELECT * FROM Places`); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return pp, nil
}

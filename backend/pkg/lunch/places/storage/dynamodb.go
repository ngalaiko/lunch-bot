package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

var _ Storage = &DynamoDBStorage{}

type DynamoDBStorage struct {
	storage   *store.DynamoDB
	tableName string
}

func NewDynamoDB(storage *store.DynamoDB, tableName string) *DynamoDBStorage {
	return &DynamoDBStorage{
		storage:   storage,
		tableName: tableName,
	}
}

func (s *DynamoDBStorage) Store(ctx context.Context, place *places.Place) error {
	if err := s.storage.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s"
			value {
				'id': ?,
				'name': ?,
				'added_at': ?,
				'added_by': ?
			}
	`, s.tableName), place.ID, place.Name, place.AddedAt.Unix(), place.AddedBy); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s *DynamoDBStorage) GetByID(ctx context.Context, id places.ID) (*places.Place, error) {
	places := []*places.Place{}
	if err := s.storage.Query(ctx, &places, fmt.Sprintf(`SELECT * FROM "%s" WHERE id = ?`, s.tableName), id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return places[0], nil
}

func (s *DynamoDBStorage) ListAll(ctx context.Context) (map[places.ID]*places.Place, error) {
	pp := []*places.Place{}
	if err := s.storage.Query(ctx, &pp, fmt.Sprintf(`SELECT * FROM "%s"`, s.tableName)); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	m := make(map[places.ID]*places.Place)
	for _, p := range pp {
		m[p.ID] = p
	}
	return m, nil
}

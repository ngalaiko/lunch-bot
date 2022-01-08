package storage

import (
	"context"
	"fmt"

	"lunch/pkg/jwt/keys"
	"lunch/pkg/store"
)

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

func (dynamodb *DynamoDBStorage) Create(ctx context.Context, key *keys.Key) error {
	if err := dynamodb.storage.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s"
			value {
				'id': ?,
				'public_der': ?
			}
	`, dynamodb.tableName), key.ID, key.PublicDER); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s *DynamoDBStorage) Get(ctx context.Context, id string) (*keys.Key, error) {
	keys := []*keys.Key{}
	if err := s.storage.Query(ctx, &keys, fmt.Sprintf(`SELECT * FROM "%s" WHERE id = ?`, s.tableName), id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	if len(keys) == 0 {
		return nil, ErrNotFound
	}
	return keys[0], nil
}

package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/boosts"
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

func (dynamodb *DynamoDBStorage) Store(ctx context.Context, boost *boosts.Boost) error {
	if err := dynamodb.storage.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s"
			value {
				'id': ?,
				'user_id': ?,
				'place_id': ?,
				'time': ?
			}
	`, dynamodb.tableName), boost.ID, boost.UserID, boost.PlaceID, boost.Time.Unix()); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (dynamo *DynamoDBStorage) ListBoosts(ctx context.Context) ([]*boosts.Boost, error) {
	bb := []*boosts.Boost{}
	if err := dynamo.storage.Query(ctx, &bb, fmt.Sprintf(`SELECT * FROM "%s"`, dynamo.tableName)); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return bb, nil
}

package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/boosts"
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

func (dynamodb *DynamoDBStorage) Store(ctx context.Context, boost *boosts.Boost) error {
	if err := dynamodb.storage.Execute(ctx, `
		INSERT INTO Boosts
			value {
				'id': ?,
				'user_id': ?,
				'place_name': ?,
				'time': ?
			}
	`, boost.ID, boost.UserID, boost.PlaceName, boost.Time.Unix()); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (dynamo *DynamoDBStorage) ListBoosts(ctx context.Context) ([]*boosts.Boost, error) {
	bb := []*boosts.Boost{}
	if err := dynamo.storage.Query(ctx, &bb, `SELECT * FROM Boosts`); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return bb, nil
}

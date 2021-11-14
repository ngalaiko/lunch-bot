package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/rolls"
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

func (dynamodb *DynamoDBStorage) Store(ctx context.Context, roll *rolls.Roll) error {
	if err := dynamodb.storage.Execute(ctx, `
		INSERT INTO Rolls
			value {
				'id': ?,
				'user_id': ?,
				'place_name': ?,
				'time': ?
			}
	`, roll.ID, roll.UserID, roll.PlaceName, roll.Time.Unix()); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (dynamo *DynamoDBStorage) ListRolls(ctx context.Context) ([]*rolls.Roll, error) {
	rr := []*rolls.Roll{}
	if err := dynamo.storage.Query(ctx, &rr, `SELECT * FROM Rolls`); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return rr, nil
}

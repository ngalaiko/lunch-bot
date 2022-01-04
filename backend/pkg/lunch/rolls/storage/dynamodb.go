package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/rolls"
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

func (dynamodb *DynamoDBStorage) Store(ctx context.Context, roll *rolls.Roll) error {
	if err := dynamodb.storage.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s"
			value {
				'id': ?,
				'user_id': ?,
				'place_id': ?,
				'time': ?
			}
	`, dynamodb.tableName), roll.ID, roll.UserID, roll.PlaceID, roll.Time.Unix()); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (dynamo *DynamoDBStorage) ListRolls(ctx context.Context) ([]*rolls.Roll, error) {
	rr := []*rolls.Roll{}
	if err := dynamo.storage.Query(ctx, &rr, fmt.Sprintf(`SELECT * FROM "%s"`, dynamo.tableName)); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return rr, nil
}

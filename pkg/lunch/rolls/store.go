package rolls

import (
	"context"
	"fmt"
	"time"

	"lunch/pkg/store"
)

type Store struct {
	storage    *store.S3Store
	bucketName string
}

func NewStore(storage *store.S3Store) *Store {
	return &Store{
		storage:    storage,
		bucketName: "lunch-rolls",
	}
}

func (rs *Store) Store(ctx context.Context, roll *Roll) error {
	if err := rs.storage.Store(ctx, rs.bucketName, roll.key(), []byte(roll.value())); err != nil {
		return fmt.Errorf("failed to store roll in storage: %w", err)
	}
	return nil
}

func (rs *Store) ListThisWeekRolls(ctx context.Context) ([]*Roll, error) {
	year, week := time.Now().ISOWeek()
	prefix := fmt.Sprintf("%d/%d/", year, week)

	keys, err := rs.storage.ListKeys(ctx, rs.bucketName, store.WithPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}

	rolls := make([]*Roll, 0, len(keys))
	for _, key := range keys {
		roll, err := rollFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key: %w", err)
		}
		rolls = append(rolls, roll)
	}

	return rolls, nil
}

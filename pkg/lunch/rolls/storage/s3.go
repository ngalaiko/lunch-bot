package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/store"
)

type Store struct {
	storage    *store.S3
	bucketName string
}

func NewS3(storage *store.S3) *Store {
	return &Store{
		storage:    storage,
		bucketName: "lunch-rolls",
	}
}

func (rs *Store) Store(ctx context.Context, roll *rolls.Roll) error {
	if err := rs.storage.Store(ctx, rs.bucketName, key(roll), []byte(value(roll))); err != nil {
		return fmt.Errorf("failed to store roll in storage: %w", err)
	}
	return nil
}

func (rs *Store) ListRolls(ctx context.Context) ([]*rolls.Roll, error) {
	keys, err := rs.storage.ListKeys(ctx, rs.bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}

	rolls := make([]*rolls.Roll, 0, len(keys))
	for _, key := range keys {
		roll, err := rollFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key: %w", err)
		}
		rolls = append(rolls, roll)
	}

	return rolls, nil
}

func rollFromKey(key string) (*rolls.Roll, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("unexpected key format")
	}
	t, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse key time: %w", err)
	}
	return &rolls.Roll{
		UserID:    parts[3],
		PlaceName: places.Name(parts[4]),
		Time:      t,
	}, nil
}

func key(r *rolls.Roll) string {
	year, week := r.Time.ISOWeek()
	return fmt.Sprintf("%d/%d/%s/%s/%s", year, week, r.Time.Format(time.RFC3339), r.UserID, r.PlaceName)
}

func value(r *rolls.Roll) string {
	return ""
}

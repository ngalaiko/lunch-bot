package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

type S3Storage struct {
	storage    *store.S3
	bucketName string
}

func NewS3(storage *store.S3) *S3Storage {
	return &S3Storage{
		storage:    storage,
		bucketName: "lunch-boosts",
	}
}

func boostFromKey(key string) (*boosts.Boost, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("unexpected key format")
	}
	t, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse key time: %w", err)
	}
	return &boosts.Boost{
		UserID:    parts[3],
		PlaceName: places.Name(parts[4]),
		Time:      t,
	}, nil
}

func key(r *boosts.Boost) string {
	year, week := r.Time.ISOWeek()
	return fmt.Sprintf("%d/%d/%s/%s/%s", year, week, r.Time.Format(time.RFC3339), r.UserID, r.PlaceName)
}

func value(r *boosts.Boost) string {
	return ""
}

func (rs *S3Storage) Store(ctx context.Context, boost *boosts.Boost) error {
	if err := rs.storage.Store(ctx, rs.bucketName, key(boost), []byte(value(boost))); err != nil {
		return fmt.Errorf("failed to store boost in storage: %w", err)
	}
	return nil
}

func (rs *S3Storage) ListBoosts(ctx context.Context) ([]*boosts.Boost, error) {
	keys, err := rs.storage.ListKeys(ctx, rs.bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys from storage: %w", err)
	}

	boosts := make([]*boosts.Boost, 0, len(keys))
	for _, key := range keys {
		boost, err := boostFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key: %w", err)
		}
		boosts = append(boosts, boost)
	}

	return boosts, nil

}

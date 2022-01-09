package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/rolls"
	"lunch/pkg/store"
)

type BoltDBStorage struct {
	db         *store.Bolt
	bucketName string
}

func NewBolt(db *store.Bolt) *BoltDBStorage {
	return &BoltDBStorage{
		db:         db,
		bucketName: "rolls",
	}
}

func (s *BoltDBStorage) Store(ctx context.Context, boost *rolls.Roll) error {
	return s.db.Put(ctx, s.bucketName, string(boost.ID), boost)
}

func (s *BoltDBStorage) ListRolls(ctx context.Context) ([]*rolls.Roll, error) {
	var rolls = []*rolls.Roll{}
	if _, err := s.db.List(ctx, s.bucketName, &rolls, 100, nil); err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}
	return rolls, nil
}
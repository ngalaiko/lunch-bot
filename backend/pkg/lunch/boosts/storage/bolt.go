package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/store"
)

type BoltDBStorage struct {
	db         *store.Bolt
	bucketName string
}

func NewBolt(db *store.Bolt) *BoltDBStorage {
	return &BoltDBStorage{
		db:         db,
		bucketName: "boosts",
	}
}

func (s *BoltDBStorage) Store(ctx context.Context, boost *boosts.Boost) error {
	return s.db.Put(ctx, s.bucketName, string(boost.ID), boost)
}

func (s *BoltDBStorage) ListBoosts(ctx context.Context) (map[boosts.ID]*boosts.Boost, error) {
	bb := []*boosts.Boost{}
	if err := s.db.List(ctx, s.bucketName, &bb); err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}
	m := make(map[boosts.ID]*boosts.Boost)
	for _, boost := range bb {
		m[boost.ID] = boost
	}
	return m, nil
}

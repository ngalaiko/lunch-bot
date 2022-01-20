package storage

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/places"
	"lunch/pkg/store"
)

type BoltDBStorage struct {
	db         *store.Bolt
	bucketName string
}

func NewBolt(db *store.Bolt) *BoltDBStorage {
	return &BoltDBStorage{
		db:         db,
		bucketName: "places",
	}
}

func (s *BoltDBStorage) Store(ctx context.Context, boost *places.Place) error {
	return s.db.Put(ctx, s.bucketName, string(boost.ID), boost)
}

func (s *BoltDBStorage) GetByID(ctx context.Context, id places.ID) (*places.Place, error) {
	var place *places.Place
	if err := s.db.Get(ctx, s.bucketName, string(id), &place); err != nil {
		return nil, fmt.Errorf("failed to get place by id: %w", err)
	}
	return place, nil
}

func (s *BoltDBStorage) ListAll(ctx context.Context) (map[places.ID]*places.Place, error) {
	pp := []*places.Place{}
	if _, err := s.db.List(ctx, s.bucketName, &pp, 100, nil); err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}
	m := make(map[places.ID]*places.Place)
	for _, p := range pp {
		m[p.ID] = p
	}
	return m, nil
}

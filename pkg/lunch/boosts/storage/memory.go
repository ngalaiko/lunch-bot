package storage

import (
	"context"
	"lunch/pkg/lunch/boosts"
)

var _ Storage = &MemoryStorage{}

type MemoryStorage struct {
	boosts []*boosts.Boost
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{}
}

func (memory *MemoryStorage) Store(_ context.Context, boost *boosts.Boost) error {
	memory.boosts = append(memory.boosts, boost)
	return nil
}

func (memory *MemoryStorage) ListBoosts(_ context.Context) ([]*boosts.Boost, error) {
	return memory.boosts, nil
}

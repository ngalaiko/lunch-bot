package storage

import (
	"context"

	"lunch/pkg/lunch/rolls"
)

var _ Storage = &MemoryStorage{}

type MemoryStorage struct {
	rolls []*rolls.Roll
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{}
}

func (memory *MemoryStorage) Store(_ context.Context, roll *rolls.Roll) error {
	memory.rolls = append(memory.rolls, roll)
	return nil
}

func (memory *MemoryStorage) ListRolls(_ context.Context) ([]*rolls.Roll, error) {
	return memory.rolls, nil
}

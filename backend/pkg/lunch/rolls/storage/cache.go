package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"lunch/pkg/lunch/rolls"
)

var _ Storage = &cache{}

type cache struct {
	storage Storage

	rolls            map[rolls.ID]*rolls.Roll
	rollsGuard       *sync.RWMutex
	rollsInitialized *int64
}

func NewCache(s Storage) *cache {
	var zero int64 = 0
	return &cache{
		storage:          s,
		rolls:            make(map[rolls.ID]*rolls.Roll),
		rollsGuard:       &sync.RWMutex{},
		rollsInitialized: &zero,
	}
}

func (c *cache) Store(ctx context.Context, roll *rolls.Roll) error {
	if err := c.storage.Store(ctx, roll); err != nil {
		return err
	}
	c.rollsGuard.Lock()
	c.rolls[roll.ID] = roll
	c.rollsGuard.Unlock()
	return nil
}

func (c *cache) ListRolls(ctx context.Context) (map[rolls.ID]*rolls.Roll, error) {
	if atomic.LoadInt64(c.rollsInitialized) == 1 {
		c.rollsGuard.RLock()
		defer c.rollsGuard.RUnlock()
		return c.rolls, nil
	}

	rolls, err := c.storage.ListRolls(ctx)
	if err != nil {
		return nil, err
	}

	c.rollsGuard.Lock()
	c.rolls = rolls
	c.rollsGuard.Unlock()
	atomic.StoreInt64(c.rollsInitialized, 1)
	return rolls, nil
}

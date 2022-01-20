package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"lunch/pkg/lunch/boosts"
)

var _ Storage = &cache{}

type cache struct {
	storage Storage

	boosts            map[boosts.ID]*boosts.Boost
	boostsGuard       *sync.RWMutex
	boostsInitialized *int64
}

func NewCache(s Storage) *cache {
	var zero int64 = 0
	return &cache{
		storage:           s,
		boosts:            make(map[boosts.ID]*boosts.Boost),
		boostsGuard:       &sync.RWMutex{},
		boostsInitialized: &zero,
	}
}

func (c *cache) Store(ctx context.Context, boost *boosts.Boost) error {
	if err := c.storage.Store(ctx, boost); err != nil {
		return err
	}

	c.boostsGuard.Lock()
	c.boosts[boost.ID] = boost
	c.boostsGuard.Unlock()
	return nil
}

func (c *cache) ListBoosts(ctx context.Context) (map[boosts.ID]*boosts.Boost, error) {
	if atomic.LoadInt64(c.boostsInitialized) == 1 {
		c.boostsGuard.RLock()
		defer c.boostsGuard.RUnlock()
		return c.boosts, nil
	}

	boosts, err := c.storage.ListBoosts(ctx)
	if err != nil {
		return nil, err
	}

	c.boostsGuard.Lock()
	c.boosts = boosts
	c.boostsGuard.Unlock()
	atomic.StoreInt64(c.boostsInitialized, 1)
	return boosts, nil
}

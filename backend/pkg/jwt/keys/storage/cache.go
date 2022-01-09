package storage

import (
	"context"
	"sync"

	"lunch/pkg/jwt/keys"
)

var _ Storage = &cache{}

type cached struct {
	key *keys.Key
	err error
}

type cache struct {
	storage Storage

	byID      map[string]*cached
	byIDGuard *sync.RWMutex
}

func NewCache(s Storage) *cache {
	return &cache{
		storage:   s,
		byID:      make(map[string]*cached),
		byIDGuard: &sync.RWMutex{},
	}
}

func (m *cache) Create(ctx context.Context, key *keys.Key) error {
	if err := m.storage.Create(ctx, key); err != nil {
		return err
	}

	m.byIDGuard.Lock()
	m.byID[key.ID] = &cached{key: key}
	m.byIDGuard.Unlock()
	return nil
}

func (m *cache) Get(ctx context.Context, id string) (*keys.Key, error) {
	m.byIDGuard.RLock()
	c, ok := m.byID[id]
	m.byIDGuard.RUnlock()
	if ok {
		return c.key, c.err
	}

	key, err := m.storage.Get(ctx, id)
	m.byIDGuard.Lock()
	m.byID[id] = &cached{key: key, err: err}
	m.byIDGuard.Unlock()
	return key, err
}

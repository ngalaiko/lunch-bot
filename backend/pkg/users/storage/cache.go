package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"lunch/pkg/users"
)

type cached struct {
	user *users.User
	err  error
}

var _ Storage = &cache{}

type cache struct {
	storage Storage

	byID      map[string]*cached
	byIDGuard *sync.RWMutex

	list            map[string]*users.User
	listGuard       *sync.RWMutex
	listInitialized *int64
}

func NewCache(s Storage) *cache {
	var zero int64 = 0
	return &cache{
		storage:   s,
		byID:      make(map[string]*cached),
		byIDGuard: &sync.RWMutex{},

		list:            make(map[string]*users.User),
		listGuard:       &sync.RWMutex{},
		listInitialized: &zero,
	}
}

func (c *cache) Create(ctx context.Context, user *users.User) error {
	if err := c.storage.Create(ctx, user); err != nil {
		return err
	}

	c.byIDGuard.Lock()
	c.byID[user.ID] = &cached{
		user: user,
	}
	c.byIDGuard.Unlock()

	c.listGuard.Lock()
	c.list[user.ID] = user
	c.listGuard.Unlock()

	return nil
}

func (c *cache) Get(ctx context.Context, id string) (*users.User, error) {
	c.byIDGuard.RLock()
	if cached, ok := c.byID[id]; ok {
		c.byIDGuard.RUnlock()
		return cached.user, cached.err
	}
	c.byIDGuard.RUnlock()

	user, err := c.storage.Get(ctx, id)

	c.byIDGuard.Lock()
	c.byID[id] = &cached{
		user: user,
		err:  err,
	}
	c.byIDGuard.Unlock()

	c.listGuard.Lock()
	c.list[user.ID] = user
	c.listGuard.Unlock()

	return user, err
}

func (c *cache) List(ctx context.Context) (map[string]*users.User, error) {
	if atomic.LoadInt64(c.listInitialized) == 1 {
		c.listGuard.RLock()
		defer c.listGuard.RUnlock()
		return c.list, nil
	}

	users, err := c.storage.List(ctx)
	if err != nil {
		return nil, err
	}

	c.listGuard.Lock()
	c.list = users
	c.listGuard.Unlock()
	atomic.StoreInt64(c.listInitialized, 1)
	return users, nil
}

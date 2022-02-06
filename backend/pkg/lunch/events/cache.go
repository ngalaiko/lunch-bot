package events

import (
	"context"
	"sync"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

var _ Storage = &cache{}

type cache struct {
	storage Storage

	byRoomID            map[rooms.ID][]*Event
	byRoomIDGuard       *sync.RWMutex
	byRoomIDInitialized bool

	byUserID            map[users.ID][]*Event
	byUserIDGuard       *sync.RWMutex
	byUserIDInitialized bool
}

func NewCache(s Storage) *cache {
	return &cache{
		storage:       s,
		byRoomID:      make(map[rooms.ID][]*Event),
		byRoomIDGuard: &sync.RWMutex{},
		byUserID:      make(map[users.ID][]*Event),
		byUserIDGuard: &sync.RWMutex{},
	}
}

func (c *cache) Create(ctx context.Context, event *Event) error {
	if err := c.storage.Create(ctx, event); err != nil {
		return err
	}

	c.byRoomIDGuard.Lock()
	c.byRoomID[event.RoomID] = append(c.byRoomID[event.RoomID], event)
	c.byRoomIDGuard.Unlock()

	c.byUserIDGuard.Lock()
	c.byUserID[event.UserID] = append(c.byUserID[event.UserID], event)
	c.byUserIDGuard.Unlock()
	return nil
}

func (c *cache) ByUserID(ctx context.Context, userID users.ID, types ...Type) ([]*Event, error) {
	c.byUserIDGuard.RLock()
	isInitialized := c.byUserIDInitialized
	events := c.byUserID[userID]
	c.byUserIDGuard.RUnlock()

	if !isInitialized {
		ee, err := c.storage.ByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		c.byUserIDGuard.Lock()
		events = ee
		c.byUserID[userID] = events
		c.byUserIDInitialized = true
		c.byUserIDGuard.Unlock()
	}

	if len(types) == 0 {
		return events, nil
	}

	var filtered []*Event
	for _, e := range events {
		for _, t := range types {
			if e.Type == t {
				filtered = append(filtered, e)
				break
			}
		}
	}
	return filtered, nil
}

func (c *cache) ByRoomID(ctx context.Context, roomID rooms.ID, types ...Type) ([]*Event, error) {
	c.byRoomIDGuard.RLock()
	isInitialized := c.byRoomIDInitialized
	events := c.byRoomID[roomID]
	c.byRoomIDGuard.RUnlock()

	if !isInitialized {
		ee, err := c.storage.ByRoomID(ctx, roomID)
		if err != nil {
			return nil, err
		}
		c.byRoomIDGuard.Lock()
		events = ee
		c.byRoomID[roomID] = events
		c.byRoomIDInitialized = true
		c.byRoomIDGuard.Unlock()
	}

	if len(types) == 0 {
		return events, nil
	}

	var filtered []*Event
	for _, e := range events {
		for _, t := range types {
			if e.Type == t {
				filtered = append(filtered, e)
				break
			}
		}
	}
	return filtered, nil
}

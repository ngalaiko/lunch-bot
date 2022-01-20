package storage

import (
	"context"
	"sync"
	"sync/atomic"

	"lunch/pkg/lunch/places"
)

var _ Storage = &cache{}

type cached struct {
	place *places.Place
	err   error
}

type cache struct {
	storage Storage

	placesByIDGuard *sync.RWMutex
	placesByID      map[places.ID]*cached

	placesListGuard       *sync.RWMutex
	placesList            map[places.ID]*places.Place
	placesListInitialized *int64
}

func NewCache(s Storage) *cache {
	var zero int64 = 0
	return &cache{
		storage: s,

		placesByIDGuard: &sync.RWMutex{},
		placesByID:      map[places.ID]*cached{},

		placesListGuard:       &sync.RWMutex{},
		placesList:            map[places.ID]*places.Place{},
		placesListInitialized: &zero,
	}
}

func (c *cache) Store(ctx context.Context, place *places.Place) error {
	if err := c.storage.Store(ctx, place); err != nil {
		return err
	}

	c.placesByIDGuard.Lock()
	c.placesByID[place.ID] = &cached{place: place}
	c.placesByIDGuard.Unlock()

	c.placesListGuard.Lock()
	c.placesList[place.ID] = place
	c.placesListGuard.Unlock()
	return nil
}

func (c *cache) GetByID(ctx context.Context, id places.ID) (*places.Place, error) {
	cachedPlace, found := c.placesByID[id]
	if found {
		return cachedPlace.place, cachedPlace.err
	}
	place, err := c.storage.GetByID(ctx, id)

	c.placesByIDGuard.Lock()
	c.placesByID[id] = &cached{place: place, err: err}
	c.placesByIDGuard.Unlock()

	if err != nil {
		return nil, err
	}

	c.placesListGuard.Lock()
	c.placesList[id] = place
	c.placesListGuard.Unlock()

	return place, nil
}

func (c *cache) ListAll(ctx context.Context) (map[places.ID]*places.Place, error) {
	if atomic.LoadInt64(c.placesListInitialized) == 1 {
		c.placesListGuard.RLock()
		defer c.placesListGuard.RUnlock()
		return c.placesList, nil
	}

	places, err := c.storage.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	c.placesListGuard.Lock()
	c.placesList = places
	c.placesListGuard.Unlock()
	atomic.StoreInt64(c.placesListInitialized, 1)

	c.placesByIDGuard.Lock()
	for _, place := range places {
		c.placesByID[place.ID] = &cached{place: place}
	}
	c.placesByIDGuard.Unlock()
	return places, nil
}

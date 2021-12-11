package storage

import (
	"context"
	"fmt"
	"lunch/pkg/lunch/places"
)

var _ Storage = &MemoryStorage{}

type MemoryStorage struct {
	places map[places.ID]*places.Place
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		places: map[places.ID]*places.Place{},
	}
}

func (memory *MemoryStorage) Store(_ context.Context, place *places.Place) error {
	memory.places[place.ID] = place
	return nil
}

func (memory *MemoryStorage) GetByID(_ context.Context, id places.ID) (*places.Place, error) {
	place, found := memory.places[id]
	if !found {
		return nil, fmt.Errorf("not found")
	}
	return place, nil
}

func (memory *MemoryStorage) ListAll(_ context.Context) ([]*places.Place, error) {
	places := []*places.Place{}
	for _, place := range memory.places {
		places = append(places, place)
	}
	return places, nil
}

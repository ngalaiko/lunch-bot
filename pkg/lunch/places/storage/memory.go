package storage

import (
	"context"
	"fmt"
	"lunch/pkg/lunch/places"
)

var _ Storage = &MemoryStorage{}

type MemoryStorage struct {
	places map[places.Name]*places.Place
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		places: map[places.Name]*places.Place{},
	}
}

func (memory *MemoryStorage) Store(_ context.Context, place *places.Place) error {
	memory.places[place.Name] = place
	return nil
}

func (memory *MemoryStorage) ListNames(_ context.Context) ([]places.Name, error) {
	names := []places.Name{}
	for name := range memory.places {
		names = append(names, name)
	}
	return names, nil
}

func (memory *MemoryStorage) GetByName(_ context.Context, name places.Name) (*places.Place, error) {
	place, found := memory.places[name]
	if !found {
		return nil, fmt.Errorf("not found")
	}
	return place, nil
}

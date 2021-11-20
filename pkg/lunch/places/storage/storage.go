package storage

import (
	"context"

	"lunch/pkg/lunch/places"
)

type Storage interface {
	Store(context.Context, *places.Place) error
	ListNames(context.Context) ([]places.Name, error)
	GetByName(context.Context, places.Name) (*places.Place, error)
	ListAll(context.Context) ([]*places.Place, error)
}

package storage

import (
	"context"

	"lunch/pkg/lunch/places"
)

type Storage interface {
	Store(context.Context, *places.Place) error
	GetByID(context.Context, places.ID) (*places.Place, error)
	ListAll(context.Context) (map[places.ID]*places.Place, error)
}

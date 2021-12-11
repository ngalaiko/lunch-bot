package storage

import (
	"context"

	"lunch/pkg/lunch/boosts"
)

type Storage interface {
	Store(context.Context, *boosts.Boost) error
	ListBoosts(context.Context) ([]*boosts.Boost, error)
}

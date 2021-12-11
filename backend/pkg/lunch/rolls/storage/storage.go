package storage

import (
	"context"

	"lunch/pkg/lunch/rolls"
)

type Storage interface {
	Store(context.Context, *rolls.Roll) error
	ListRolls(context.Context) ([]*rolls.Roll, error)
}

package storage

import (
	"context"

	"lunch/pkg/users"
)

type Storage interface {
	Create(context.Context, *users.User) error
	Get(context.Context, string) (*users.User, error)
	List(context.Context) ([]*users.User, error)
}

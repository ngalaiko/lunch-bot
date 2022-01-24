package storage

import (
	"context"
	"fmt"

	"lunch/pkg/users"
)

var ErrNotFound = fmt.Errorf("not found")

type Storage interface {
	Create(context.Context, *users.User) error
	Get(context.Context, users.ID) (*users.User, error)
	List(context.Context) (map[users.ID]*users.User, error)
}

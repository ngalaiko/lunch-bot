package storage

import (
	"context"
	"fmt"

	"lunch/pkg/users"
)

var ErrNotFound = fmt.Errorf("not found")

type Storage interface {
	Create(context.Context, *users.User) error
	Get(context.Context, string) (*users.User, error)
	List(context.Context) ([]*users.User, error)
}

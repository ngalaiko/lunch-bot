package storage

import (
	"context"
	"fmt"

	"lunch/pkg/jwt/keys"
)

var (
	ErrNotFound = fmt.Errorf("key not found")
)

type Storage interface {
	Create(context.Context, *keys.Key) error
	Get(context.Context, string) (*keys.Key, error)
}

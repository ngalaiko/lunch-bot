package storage

import (
	"context"
	"errors"
	"fmt"

	"lunch/pkg/store"
	"lunch/pkg/users"
)

var _ Storage = &bolt{}

type bolt struct {
	db         *store.Bolt
	bucketName string
}

func NewBolt(db *store.Bolt) *bolt {
	return &bolt{
		db:         db,
		bucketName: "users",
	}
}

func (s *bolt) Create(ctx context.Context, user *users.User) error {
	return s.db.Put(ctx, s.bucketName, string(user.ID), user)
}

func (s *bolt) Get(ctx context.Context, id users.ID) (*users.User, error) {
	var user *users.User
	if err := s.db.Get(ctx, s.bucketName, string(id), &user); errors.Is(err, store.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *bolt) List(ctx context.Context) (map[users.ID]*users.User, error) {
	uu := []*users.User{}
	if err := s.db.List(ctx, s.bucketName, &uu); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	m := make(map[users.ID]*users.User, len(uu))
	for _, u := range uu {
		m[u.ID] = u
	}
	return m, nil
}

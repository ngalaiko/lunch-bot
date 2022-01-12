package storage

import (
	"context"
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

func (s *bolt) Get(ctx context.Context, id string) (*users.User, error) {
	var user *users.User
	if err := s.db.Get(ctx, s.bucketName, id, &user); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *bolt) List(ctx context.Context) ([]*users.User, error) {
	var users []*users.User
	if _, err := s.db.List(ctx, s.bucketName, &users, 100, nil); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	"lunch/pkg/users"
	"lunch/pkg/users/storage"
)

type Service struct {
	store storage.Storage
}

func New(store storage.Storage) *Service {
	return &Service{store: store}
}

func (s *Service) Get(ctx context.Context, userID string) (*users.User, error) {
	return s.store.Get(ctx, userID)
}

func (s *Service) List(ctx context.Context) ([]*users.User, error) {
	return s.store.List(ctx)
}

func (s *Service) Create(ctx context.Context, user *users.User) error {
	if _, err := s.store.Get(ctx, user.ID); errors.Is(err, storage.ErrNotFound) {
		return s.store.Create(ctx, user)
	} else if err == nil {
		return nil
	} else {
		return fmt.Errorf("failed to get user: %w", err)
	}
}

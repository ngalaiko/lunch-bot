package storage

import (
	"context"

	"lunch/pkg/jwt/keys"
)

type memory struct {
	byID map[string]*keys.Key
}

func NewMemory() *memory {
	return &memory{
		byID: make(map[string]*keys.Key),
	}
}

func (m *memory) Create(_ context.Context, key *keys.Key) error {
	m.byID[key.ID] = key
	return nil
}

func (m *memory) Get(_ context.Context, id string) (*keys.Key, error) {
	key, ok := m.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return key, nil
}

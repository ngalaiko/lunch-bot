package events

import (
	"context"
	"fmt"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/store"
	"lunch/pkg/users"
)

var _ Storage = &boltStorage{}

type boltStorage struct {
	db         *store.Bolt
	bucketName string
}

func NewBoltStorage(db *store.Bolt) *boltStorage {
	return &boltStorage{
		db:         db,
		bucketName: "events",
	}
}

func (b *boltStorage) Create(ctx context.Context, event *Event) error {
	return b.db.Put(ctx, b.bucketName, fmt.Sprint(event.Timestamp), event)
}

func (b *boltStorage) ByUserID(ctx context.Context, userID users.ID, types ...Type) ([]*Event, error) {
	events := []*Event{}
	if err := b.db.List(ctx, b.bucketName, &events); err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	result := []*Event{}
	tmap := map[Type]bool{}
	for _, t := range types {
		tmap[t] = true
	}
	for _, event := range events {
		if len(types) > 0 && !tmap[event.Type] {
			continue
		}
		if event.UserID == userID {
			result = append(result, event)
		}
	}
	return result, nil
}

func (b *boltStorage) ByRoomID(ctx context.Context, roomID rooms.ID, types ...Type) ([]*Event, error) {
	events := []*Event{}
	if err := b.db.List(ctx, b.bucketName, &events); err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	result := []*Event{}
	tmap := map[Type]bool{}
	for _, t := range types {
		tmap[t] = true
	}
	for _, event := range events {
		if len(types) > 0 && !tmap[event.Type] {
			continue
		}
		if event.RoomID == roomID {
			result = append(result, event)
		}
	}
	return result, nil
}

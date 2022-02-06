package storage

import (
	"context"
	"fmt"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

const (
	placeCreated events.Type = "places/created"
)

type Storage struct {
	storage events.Storage
}

func New(storage events.Storage) *Storage {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Create(ctx context.Context, place *places.Place) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    place.UserID,
		RoomID:    place.RoomID,
		Timestamp: time.Now(),
		Type:      placeCreated,
		PlaceID:   places.ID(uuid.NewString()),
		Name:      place.Name,
	})
}

func (s *Storage) Place(ctx context.Context, roomID rooms.ID, placeID places.ID) (*places.Place, error) {
	events, err := s.storage.ByRoomID(ctx, roomID, placeCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	for _, event := range events {
		if event.PlaceID != placeID {
			continue
		}
		return &places.Place{
			ID:     event.PlaceID,
			Name:   event.Name,
			UserID: event.UserID,
			Time:   event.Timestamp,
			RoomID: event.RoomID,
		}, nil
	}
	return nil, ErrNotFound
}

func (s *Storage) Places(ctx context.Context, roomID rooms.ID) (map[places.ID]*places.Place, error) {
	events, err := s.storage.ByRoomID(ctx, roomID, placeCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	result := make(map[places.ID]*places.Place)
	for _, event := range events {
		result[event.PlaceID] = &places.Place{
			ID:     event.PlaceID,
			Name:   event.Name,
			UserID: event.UserID,
			Time:   event.Timestamp,
			RoomID: event.RoomID,
		}
	}
	return result, nil

}

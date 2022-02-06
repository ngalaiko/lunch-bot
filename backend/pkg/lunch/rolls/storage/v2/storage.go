package storage

import (
	"context"
	"fmt"
	"time"

	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/lunch/rooms"
)

var (
	rollCreated events.Type = "rolls/created"
)

type Storage struct {
	eventsStorage events.Storage
}

func New(eventsStorage events.Storage) *Storage {
	return &Storage{
		eventsStorage: eventsStorage,
	}
}

func (s *Storage) Create(ctx context.Context, roll *rolls.Roll) error {
	return s.eventsStorage.Create(ctx, &events.Event{
		UserID:    roll.UserID,
		PlaceID:   roll.PlaceID,
		RoomID:    roll.RoomID,
		Type:      rollCreated,
		Timestamp: time.Now(),
	})
}

func (s *Storage) Rolls(ctx context.Context, roomID rooms.ID) (map[places.ID]*rolls.Roll, error) {
	events, err := s.eventsStorage.ByRoomID(ctx, roomID, rollCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	result := make(map[places.ID]*rolls.Roll)
	for _, event := range events {
		result[event.PlaceID] = &rolls.Roll{
			UserID:  event.UserID,
			PlaceID: event.PlaceID,
			Time:    event.Timestamp,
		}
	}
	return result, nil
}

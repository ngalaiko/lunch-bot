package storage

import (
	"context"
	"fmt"
	"time"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
)

var (
	boostCreated events.Type = "boosts/created"
)

type Storage struct {
	eventsStorage events.Storage
}

func New(eventsStorage events.Storage) *Storage {
	return &Storage{
		eventsStorage: eventsStorage,
	}
}

func (s *Storage) Create(ctx context.Context, boost *boosts.Boost) error {
	return s.eventsStorage.Create(ctx, &events.Event{
		UserID:    boost.UserID,
		PlaceID:   boost.PlaceID,
		RoomID:    boost.RoomID,
		Type:      boostCreated,
		Timestamp: time.Now(),
	})
}

func (s *Storage) Boosts(ctx context.Context, roomID rooms.ID) (map[places.ID]*boosts.Boost, error) {
	events, err := s.eventsStorage.ByRoomID(ctx, roomID, boostCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	result := make(map[places.ID]*boosts.Boost)
	for _, event := range events {
		result[event.PlaceID] = &boosts.Boost{
			UserID:  event.UserID,
			PlaceID: event.PlaceID,
			Time:    event.Timestamp,
		}
	}
	return result, nil
}

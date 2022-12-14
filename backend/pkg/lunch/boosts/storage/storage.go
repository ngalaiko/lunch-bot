package storage

import (
	"context"
	"fmt"
	"time"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/events"
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
		Timestamp: events.UnixNanoTime(boost.Time),
	})
}

func (s *Storage) Boosts(ctx context.Context, roomID rooms.ID) ([]*boosts.Boost, error) {
	events, err := s.eventsStorage.ByRoomID(ctx, roomID, boostCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	result := make([]*boosts.Boost, 0, len(events))
	for _, event := range events {
		result = append(result, &boosts.Boost{
			UserID:  event.UserID,
			PlaceID: event.PlaceID,
			Time:    time.Time(event.Timestamp),
		})
	}
	return result, nil
}

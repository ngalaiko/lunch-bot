package storage

import (
	"context"
	"fmt"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
	"sort"
	"time"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

const (
	placeCreated  events.Type = "places/created"
	placeDeleted  events.Type = "places/deleted"
	placeRestored events.Type = "places/restored"
)

type Storage struct {
	storage events.Storage
}

func New(storage events.Storage) *Storage {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Restore(ctx context.Context, userID users.ID, place *places.Place) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    userID,
		RoomID:    place.RoomID,
		Timestamp: events.UnixNanoTime(time.Now()),
		Type:      placeRestored,
		PlaceID:   place.ID,
	})
}

func (s *Storage) Delete(ctx context.Context, userID users.ID, place *places.Place) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    userID,
		RoomID:    place.RoomID,
		Timestamp: events.UnixNanoTime(time.Now()),
		Type:      placeDeleted,
		PlaceID:   place.ID,
	})
}

func (s *Storage) Create(ctx context.Context, place *places.Place) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    place.UserID,
		RoomID:    place.RoomID,
		Timestamp: events.UnixNanoTime(place.Time),
		Type:      placeCreated,
		PlaceID:   place.ID,
		Name:      place.Name,
	})
}

func (s *Storage) Place(ctx context.Context, roomID rooms.ID, placeID places.ID) (*places.Place, error) {
	places, err := s.Places(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get places: %w", err)
	}
	place, ok := places[placeID]
	if !ok {
		return nil, ErrNotFound
	}
	return place, nil
}

func (s *Storage) Places(ctx context.Context, roomID rooms.ID) (map[places.ID]*places.Place, error) {
	events, err := s.storage.ByRoomID(ctx, roomID, placeCreated, placeDeleted, placeRestored)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	sort.Slice(events, func(i, j int) bool {
		return time.Time(events[i].Timestamp).Before(time.Time(events[j].Timestamp))
	})
	result := make(map[places.ID]*places.Place)
	for _, event := range events {
		switch event.Type {
		case placeCreated:
			result[event.PlaceID] = &places.Place{
				ID:     event.PlaceID,
				Name:   event.Name,
				UserID: event.UserID,
				Time:   time.Time(event.Timestamp),
				RoomID: event.RoomID,
			}
		case placeDeleted:
			result[event.PlaceID].IsDeleted = true
		case placeRestored:
			result[event.PlaceID].IsDeleted = false
		}
	}
	return result, nil
}

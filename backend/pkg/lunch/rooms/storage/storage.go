package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

const (
	roomCreated events.Type = "rooms/created"
	roomJoined  events.Type = "rooms/joined"
	roomLeft    events.Type = "rooms/left"
)

type Storage struct {
	storage events.Storage
}

func New(storage events.Storage) *Storage {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) Create(ctx context.Context, room *rooms.Room) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    room.UserID,
		Timestamp: time.Now(),
		Type:      roomCreated,
		RoomID:    rooms.ID(uuid.NewString()),
		Name:      room.Name,
	})
}

func (s *Storage) Join(ctx context.Context, user *users.User, roomID rooms.ID) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    user.ID,
		Timestamp: time.Now(),
		Type:      roomJoined,
		RoomID:    roomID,
	})
}

func (s *Storage) Leave(ctx context.Context, user *users.User, roomID rooms.ID) error {
	return s.storage.Create(ctx, &events.Event{
		UserID:    user.ID,
		Timestamp: time.Now(),
		Type:      roomLeft,
		RoomID:    roomID,
	})
}

func (s *Storage) Room(ctx context.Context, roomID rooms.ID) (*rooms.Room, error) {
	events, err := s.storage.ByRoomID(ctx, roomID, roomCreated, roomJoined, roomLeft)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	if len(events) == 0 {
		return nil, ErrNotFound
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	memberIDs := make(map[users.ID]bool)
	for _, event := range events {
		switch event.Type {
		case roomCreated:
			memberIDs[event.UserID] = true
		case roomJoined:
			memberIDs[event.UserID] = true
		case roomLeft:
			delete(memberIDs, event.UserID)
		}
	}

	return &rooms.Room{
		ID:        roomID,
		Name:      events[0].Name,
		UserID:    events[0].UserID,
		Time:      events[0].Timestamp,
		MemberIDs: memberIDs,
	}, nil
}

func (s *Storage) Rooms(ctx context.Context, userID users.ID) (map[rooms.ID]bool, error) {
	events, err := s.storage.ByUserID(ctx, userID, roomCreated, roomJoined, roomLeft)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})
	result := make(map[rooms.ID]bool)
	for _, event := range events {
		switch event.Type {
		case roomCreated:
			result[event.RoomID] = true
		case roomJoined:
			result[event.RoomID] = true
		case roomLeft:
			delete(result, event.RoomID)
		}
	}
	return result, nil
}

package events

import (
	"context"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

type Storage interface {
	// Create stores an new event.
	Create(context.Context, *Event) error
	// ByUserID returns all events for a given user id.
	// If no types are specified, all events are returned, otherwise only events of the given types are returned.
	ByUserID(context.Context, users.ID, ...Type) ([]*Event, error)
	// ByRoomID returns all events for a given room id.
	// If no types are specified, all events are returned, otherwise only events of the given types are returned.
	ByRoomID(context.Context, rooms.ID, ...Type) ([]*Event, error)
}

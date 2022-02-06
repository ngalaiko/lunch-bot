package events

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

type Type string

type Event struct {
	UserID    users.ID  `dynamodbav:"user_id"`
	RoomID    rooms.ID  `dynamodbav:"room_id"`
	Type      Type      `dynamodbav:"type"`
	Timestamp time.Time `dynamodbav:"timestamp,unixtime"`
	PlaceID   places.ID `dynamodbav:"place_id"`
	Name      string    `dynamodbav:"name"`
}

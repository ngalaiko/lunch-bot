package rolls

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Roll struct {
	ID      ID        `dynamodbav:"id" json:"id"`
	UserID  users.ID  `dynamodbav:"user_id" json:"userId"`
	PlaceID places.ID `dynamodbav:"place_id" json:"placeId"`
	Time    time.Time `dynamodbav:"time,unixtime" json:"time"`
	RoomID  rooms.ID
}

func NewRoll(userID users.ID /* roomID rooms.ID,*/, placeID places.ID, now time.Time) *Roll {
	return &Roll{
		ID:     ID(uuid.NewString()),
		UserID: userID,
		// RoomID:  roomID,
		PlaceID: placeID,
		Time:    now,
	}
}

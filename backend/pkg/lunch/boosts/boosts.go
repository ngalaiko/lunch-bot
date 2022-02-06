package boosts

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

type ID string

type Boost struct {
	ID      ID        `dynamodbav:"id" json:"id"`
	UserID  users.ID  `dynamodbav:"user_id" json:"userId"`
	PlaceID places.ID `dynamodbav:"place_id" json:"placeId"`
	Time    time.Time `dynamodbav:"time,unixtime" json:"time"`
	RoomID  rooms.ID  `json:"roomId"`
}

func NewBoost(userID users.ID, roomID rooms.ID, placeID places.ID, now time.Time) *Boost {
	return &Boost{
		UserID:  userID,
		PlaceID: placeID,
		RoomID:  roomID,
		Time:    now,
	}
}

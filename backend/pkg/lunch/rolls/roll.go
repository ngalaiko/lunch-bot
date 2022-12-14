package rolls

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

type Roll struct {
	UserID  users.ID  `dynamodbav:"user_id" json:"userId"`
	PlaceID places.ID `dynamodbav:"place_id" json:"placeId"`
	Time    time.Time `dynamodbav:"time,unixtime" json:"time"`
	RoomID  rooms.ID  `json:"roomId"`
}

func NewRoll(userID users.ID, roomID rooms.ID, placeID places.ID, now time.Time) *Roll {
	return &Roll{
		UserID:  userID,
		RoomID:  roomID,
		PlaceID: placeID,
		Time:    now,
	}
}

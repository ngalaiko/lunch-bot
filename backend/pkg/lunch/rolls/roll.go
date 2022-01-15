package rolls

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Roll struct {
	ID      ID        `dynamodbav:"id" json:"id"`
	UserID  string    `dynamodbav:"user_id" json:"userId"`
	PlaceID places.ID `dynamodbav:"place_id" json:"placeId"`
	Time    time.Time `dynamodbav:"time,unixtime" json:"time"`
}

func NewRoll(user *users.User, placeID places.ID, now time.Time) *Roll {
	return &Roll{
		ID:      ID(uuid.NewString()),
		UserID:  user.ID,
		PlaceID: placeID,
		Time:    now,
	}
}

package boosts

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Boost struct {
	ID      ID        `dynamodbav:"id" json:"id"`
	UserID  string    `dynamodbav:"user_id" json:"userId"`
	PlaceID places.ID `dynamodbav:"place_id" json:"placeId"`
	Time    time.Time `dynamodbav:"time,unixtime" json:"time"`
}

func NewBoost(user *users.User, placeID places.ID, now time.Time) *Boost {
	return &Boost{
		ID:      ID(uuid.NewString()),
		UserID:  user.ID,
		PlaceID: placeID,
		Time:    now,
	}
}

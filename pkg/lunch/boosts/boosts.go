package boosts

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type Boost struct {
	ID        string      `dynamodbav:"id"`
	UserID    string      `dynamodbav:"user_id"`
	PlaceName places.Name `dynamodbav:"place_name"`
	Time      time.Time   `dynamodbav:"time,unixtime"`
}

func NewBoost(user *users.User, placeName places.Name, now time.Time) *Boost {
	return &Boost{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		PlaceName: placeName,
		Time:      now,
	}
}
